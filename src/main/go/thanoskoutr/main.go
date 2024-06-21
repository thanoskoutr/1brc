package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// TODO: Delete logs
// TODO: Delete run time calculations

const (
	WorkersCount    = 2048
	LinesBufferSize = 10000
	ShardCount      = 2048
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

type ShardedMeasurements struct {
	Shards [ShardCount]map[string]*Stats
	Mu     [ShardCount]sync.Mutex
}

func NewShardedMeasurements() *ShardedMeasurements {
	sm := &ShardedMeasurements{}
	for i := 0; i < ShardCount; i++ {
		sm.Shards[i] = make(map[string]*Stats)
	}
	return sm
}

func (sm *ShardedMeasurements) getShard(key string) (map[string]*Stats, *sync.Mutex) {
	hash := fnv32(key)
	index := hash % ShardCount
	return sm.Shards[index], &sm.Mu[index]
}

func (sm *ShardedMeasurements) update(location string, temperature float64) {
	shard, mu := sm.getShard(location)
	mu.Lock()
	defer mu.Unlock()

	// Calculate min, max, mean
	if stats, exists := shard[location]; !exists {
		shard[location] = &Stats{
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

func (sm *ShardedMeasurements) combineMeasurements() map[string]*Stats {
	combined := make(map[string]*Stats)
	for i := 0; i < ShardCount; i++ {
		sm.Mu[i].Lock()
		for location, stats := range sm.Shards[i] {
			if existingStats, exists := combined[location]; !exists {
				combined[location] = &Stats{
					Min:   stats.Min,
					Max:   stats.Max,
					Mean:  stats.Mean,
					Sum:   stats.Sum,
					Count: stats.Count,
				}
			} else {
				if stats.Min < existingStats.Min {
					existingStats.Min = stats.Min
				}
				if stats.Max > existingStats.Max {
					existingStats.Max = stats.Max
				}
				existingStats.Count += stats.Count
				existingStats.Sum += stats.Sum
				// Do not have to calculate the Mean here, do it on print
			}
		}
		sm.Mu[i].Unlock()
	}
	return combined
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

func processFile(file *os.File) *ShardedMeasurements {
	// Read file line-by-line
	var lineCount = 0
	var countInterval = 1000000
	shardedMeasurements := NewShardedMeasurements()

	// Channel for reading line in parallel
	var wg sync.WaitGroup
	lines := make(chan string, LinesBufferSize)

	go func() {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			// Count lines
			lineCount++
			if lineCount%countInterval == 0 {
				log.Printf("Read %d lines so far...\n", lineCount)
			}
			// Read line
			lines <- scanner.Text()
		}
		close(lines)
		if err := scanner.Err(); err != nil {
			fmt.Printf("Error reading measurements file: %v\n", err)
		}
		log.Printf("Count of lines in file: %v\n", lineCount)
	}()

	for i := 0; i < WorkersCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for line := range lines {
				processLine(line, shardedMeasurements)
			}
		}()
	}

	wg.Wait()

	return shardedMeasurements
}

func processLine(line string, shardedMeasurements *ShardedMeasurements) {
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

	shardedMeasurements.update(location, temperature)
}

// round rounds floats to 1 decimal place with 0.05 rounding up to 0.1
func round(x float64) float64 {
	return math.Floor((x+0.05)*10) / 10
}

// fnv32 implements the Fowler-Noll-Vo (FNV) algorithm for 32-bit numbers, which
// a popular non-cryptographic hash algorithm, with good distribution and speed.
func fnv32(key string) uint32 {
	const prime32 = 16777619
	hash := uint32(2166136261)
	for i := 0; i < len(key); i++ {
		hash *= prime32
		hash ^= uint32(key[i])
	}
	return hash
}

func main() {
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
	shardedMeasurements := processFile(file)
	finished = time.Now()
	total = finished.Sub(started)
	log.Printf("Total reading and calculation time: %v\n", total)

	// Print results
	started = time.Now()
	measurements := shardedMeasurements.combineMeasurements()
	results := formatMeasurements(measurements)
	finished = time.Now()
	total = finished.Sub(started)
	log.Printf("Total printing results time: %v\n", total)
	log.Printf("Calculated measurements for each location:\n")
	fmt.Println(results)
}
