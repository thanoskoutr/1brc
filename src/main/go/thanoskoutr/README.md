# Prerequisites
The following are the prerequisites already listed in the main README of the repository.

## Install Java 21
For Ubuntu and Debian like distributions run:
```bash
sudo apt install openjdk-21-jre openjdk-21-jdk
```

## Generate measurements
Verify environment and generate 1 billion measurements by running the following, from the project root directory:
```bash
./mvnw clean verify
./create_measurements.sh 1000000000
```

Move 1 billion line measurements generated file to avoid being removed by test scripts:
```bash
mv measurements.txt measurements_1B.txt
```

Calculate the baseline numbers and save results to expected out file:
```bash
./calculate_average_baseline.sh > measurements_1B.out
```

# Build

## Local
To build the binary locally, run from this directory:
```bash
make all
```

## Docker
To build the docker image, run from the project root directory:
```bash
./prepare_thanoskoutr.sh
```

# Run

## Local
To run the binary locally with the measurements, run from this directory:
```bash
./1bc ../../../../measurements.txt
```

## Docker
To run the docker image with the measurements, run from the project root directory:
```bash
docker run -ti --rm 1brc:thanoskoutr /1brc ./measurements.txt
```

# Validate
To validate the results of the solution on the default sample measurements, run from the project root directory:
```bash
# builds and runs the binary with the prepare and calculate scripts
./test.sh thanoskoutr
```

To validate the results of the solution on the 1 billion measures file, run from the project root directory:
```bash
# builds and runs the binary with the prepare and calculate scripts
./test.sh thanoskoutr measurements_1B.txt
```

# Evaluate
To quickly check the execution time of the solution, run from the project root directory:
```bash
./prepare_thanoskoutr.sh
time ./calculate_average_thanoskoutr.sh measurements_1B.txt
```

To benchmark properly the solution and print a performance summary, run from the project root directory:
```bash
./evaluate.sh thanoskoutr
./evaluate_10K.sh thanoskoutr
```

# Profiling

## Enable profiling
Enable CPU and Memory profiling, by editing the calculate script that runs the binary:
```bash
go run main.go -cpuprofile -memprofile
target/thanoskoutr/1brc "$INPUT" -cpuprofile -memprofile
```

Enable live profiling with HTTP server on port `6060`, by editing the calculate script that runs the binary:
```bash
target/thanoskoutr/1brc "$INPUT" -httpprofile
```

To check the live profiling results, while the program is still running navigate to the following path: `http://localhost:6060/debug/pprof/`

## Analyze results
To check the produced CPU and Memory profiling results, use the `pprof` tool to open them in a web page on port `6060`:
```bash
go tool pprof -http=:6060 -no_browser cpuprofile.prof
go tool pprof -http=:6060 -no_browser cpuprofile.prof
```
