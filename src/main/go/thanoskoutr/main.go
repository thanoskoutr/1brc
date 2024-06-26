package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime/pprof"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// TODO: Delete logs
// TODO: Delete run time calculations

const (
	WorkersCount    = 16
	LinesBufferSize = 1000000 // Max: 1000000000 (not any difference), Min: 0 (significantly slower)
)

type Measurement struct {
	Location    string
	Temperature float64
}

type Stats struct {
	Min   float64
	Max   float64
	Mean  float64
	Count int
	Sum   float64
}

func formatMeasurements(measurements map[string]*Stats) string {
	// Collect and sort the station names
	var stations []string
	for station := range measurements {
		stations = append(stations, station)
	}
	sort.Strings(stations)
	stationsLen := len(stations)

	// Print the results
	var sb strings.Builder
	sb.WriteString("{")
	for i, station := range stations {
		stats := measurements[station]
		// round the sum instead of the final result, to remove float precision errors
		mean := round(round(stats.Sum) / float64(stats.Count))
		sb.WriteString(fmt.Sprintf("%s=%.1f/%.1f/%.1f", station, stats.Min, mean, stats.Max))
		if i != stationsLen-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString("}")
	return sb.String()
}

func processFile(file *os.File) []map[string]*Stats {
	// Read file lines into an array
	var lines []string
	var lineLogCount = 0
	var countInterval = 1000000
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Count lines
		lineLogCount++
		if lineLogCount%countInterval == 0 {
			log.Printf("Read %d lines so far...\n", lineLogCount)
		}
		lines = append(lines, scanner.Text())
		// TODO: Very memory consuming. Max lines read: 300_000_000
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading measurements file: %v", err)
	}

	// Calculate the chunk size
	lineCount := len(lines)
	log.Printf("Count of lines in file: %v\n", lineCount)
	// Ensure that workers do not exceed the line count
	workerCount := min(lineCount, WorkersCount)
	log.Printf("workerCount: %v\n", workerCount)
	// Make sure that all lines are evenly distributed among the workers
	chunkSize := (lineCount + workerCount - 1) / workerCount
	log.Printf("chunkSize: %v\n", chunkSize)

	// Each worker has a partial measurement
	var wg sync.WaitGroup
	partialMeasurements := make([]map[string]*Stats, workerCount)

	// Initialize partial results maps
	for i := 0; i < workerCount; i++ {
		partialMeasurements[i] = make(map[string]*Stats)
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			start := workerID * chunkSize
			// log.Printf("ID=%v: start: %v\n", workerID, start)
			if start >= lineCount {
				return
			}
			end := start + chunkSize
			// log.Printf("ID=%v: end: %v\n", workerID, end)
			if end > lineCount {
				end = lineCount
			}
			processChunk(lines[start:end], partialMeasurements[workerID])
		}(i)
	}

	wg.Wait()

	return partialMeasurements
}

func processChunk(chunk []string, measurements map[string]*Stats) {
	for _, line := range chunk {
		parts := strings.Split(line, ";")
		if len(parts) != 2 {
			fmt.Printf("Skipping invalid line: %s\n", line)
			return
		}
		location := parts[0]
		temperature, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			fmt.Printf("Skipping invalid temperature: %s\n", parts[1])
			return
		}

		// Calculate min, max, mean
		if stats, exists := measurements[location]; !exists {
			measurements[location] = &Stats{
				Min:   temperature,
				Max:   temperature,
				Mean:  temperature,
				Sum:   temperature,
				Count: 1,
			}
		} else {
			if temperature < stats.Min {
				stats.Min = temperature
			}
			if temperature > stats.Max {
				stats.Max = temperature
			}
			stats.Count++
			stats.Sum += temperature
			// Do not have to calculate the Mean here, do it on print
		}
	}
}

func combineMeasurements(partialMeasurements []map[string]*Stats) map[string]*Stats {
	combined := make(map[string]*Stats)
	for _, partial := range partialMeasurements {
		for location, stats := range partial {
			if finalStats, exists := combined[location]; !exists {
				combined[location] = &Stats{
					Min:   stats.Min,
					Max:   stats.Max,
					Mean:  stats.Mean,
					Sum:   stats.Sum,
					Count: stats.Count,
				}
			} else {
				if stats.Min < finalStats.Min {
					finalStats.Min = stats.Min
				}
				if stats.Max > finalStats.Max {
					finalStats.Max = stats.Max
				}
				finalStats.Count += stats.Count
				finalStats.Sum += stats.Sum
				// Do not have to calculate the Mean here, do it on print
			}
		}
	}
	return combined
}

// round rounds floats to 1 decimal place with 0.05 rounding up to 0.1
func round(x float64) float64 {
	return math.Floor((x+0.05)*10) / 10
}

func main() {
	var (
		cpuprofile  = flag.Bool("cpuprofile", false, "write cpu profile to `file`")
		memprofile  = flag.Bool("memprofile", false, "write memory profile to `file`")
		httpprofile = flag.Bool("httpprofile", false, "run HTTP server for runtime profiling")
	)

	// Start CPU profiling
	if *cpuprofile {
		f, err := os.Create("cpuprofile.prof")
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	// Start memory profiling
	if *memprofile {
		fMem, err := os.Create("memprofile.prof")
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer func() {
			if err := pprof.WriteHeapProfile(fMem); err != nil {
				log.Fatal("could not write memory profile: ", err)
			}
			fMem.Close()
		}()
	}

	// Start live server for profiling on run time
	if *httpprofile {
		go func() {
			log.Println(http.ListenAndServe("localhost:6060", nil))
		}()
	}

	// Open file
	if len(os.Args) < 2 {
		fmt.Printf("No measurements file input given\n")
		os.Exit(1)
	}
	filename := os.Args[1]
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening measurements file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Keep execution times
	var started, finished time.Time
	var total time.Duration

	// Process file
	started = time.Now()
	partialMeasurements := processFile(file)
	finished = time.Now()
	total = finished.Sub(started)
	log.Printf("Total reading and calculation time: %v\n", total)

	// Combine results
	started = time.Now()
	measurements := combineMeasurements(partialMeasurements)
	// Print results
	results := formatMeasurements(measurements)
	finished = time.Now()
	total = finished.Sub(started)
	log.Printf("Total printing results time: %v\n", total)
	log.Printf("Calculated measurements for each location:\n")
	fmt.Println(results)
}
