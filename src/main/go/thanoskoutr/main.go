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
	"time"
)

// TODO: Delete logs
// TODO: Delete run time calculations

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

func processFile(file *os.File) *map[string]*Stats {
	// Read file line-by-line
	var lineCount = 0
	var countInterval = 1000000
	measurements := make(map[string]*Stats)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Count lines
		lineCount++
		if lineCount%countInterval == 0 {
			log.Printf("Read %d lines so far...\n", lineCount)
		}

		// Parse line
		line := scanner.Text()
		parts := strings.Split(line, ";")
		if len(parts) != 2 {
			fmt.Printf("Skipping invalid line at %d: %s\n", lineCount, line)
			continue
		}
		location := parts[0]
		temperature, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			fmt.Printf("Skipping invalid temperature at line %d: %s\n", lineCount, parts[1])
			continue
		}

		// Calculate min, max, mean
		if _, exists := measurements[location]; !exists {
			measurements[location] = &Stats{
				Min:   temperature,
				Max:   temperature,
				Mean:  temperature,
				Sum:   temperature,
				Count: 1,
			}
		} else {
			stats := measurements[location]
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
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading measurements file: %v\n", err)
	}
	log.Printf("Count of lines in file: %v\n", lineCount)

	return &measurements
}

// rounding floats to 1 decimal place with 0.05 rounding up to 0.1
func round(x float64) float64 {
	return math.Floor((x+0.05)*10) / 10
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
	measurements := processFile(file)
	finished = time.Now()
	total = finished.Sub(started)
	log.Printf("Total reading and calculation time: %v\n", total)

	// Print results
	started = time.Now()
	results := formatMeasurements(*measurements)
	finished = time.Now()
	total = finished.Sub(started)
	log.Printf("Total printing results time: %v\n", total)
	log.Printf("Calculated measurements for each location:\n")
	fmt.Println(results)
}
