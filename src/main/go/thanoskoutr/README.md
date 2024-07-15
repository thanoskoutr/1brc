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

# Performance

## System info
For performance measurements, my laptop was used, with the following hardware characteristics:
- Model: HP ZBook Fury 15.6 inch G8 Mobile Workstation PC
- CPU: 11th Gen Intel(R) Core(TM) i7-11850H @ 2.50GHz
- Memory: 32 GB
- OS: Ubuntu 22.04.4 LTS

CPU details:
```bash
$ cat /proc/cpuinfo

processor       : 0
vendor_id       : GenuineIntel
cpu family      : 6
model           : 141
model name      : 11th Gen Intel(R) Core(TM) i7-11850H @ 2.50GHz
stepping        : 1
microcode       : 0x50
cpu MHz         : 1299.983
cache size      : 24576 KB
physical id     : 0
siblings        : 16
core id         : 0
cpu cores       : 8
apicid          : 0
initial apicid  : 0
fpu             : yes
fpu_exception   : yes
cpuid level     : 27
wp              : yes
flags           : fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc art arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc cpuid aperfmperf tsc_known_freq pni pclmulqdq dtes64 monitor ds_cpl vmx smx est tm2 ssse3 sdbg fma cx16 xtpr pdcm pcid sse4_1 sse4_2 x2apic movbe popcnt tsc_deadline_timer aes xsave avx f16c rdrand lahf_lm abm 3dnowprefetch cpuid_fault epb cat_l2 invpcid_single cdp_l2 ssbd ibrs ibpb stibp ibrs_enhanced tpr_shadow flexpriority ept vpid ept_ad fsgsbase tsc_adjust bmi1 avx2 smep bmi2 erms invpcid rdt_a avx512f avx512dq rdseed adx smap avx512ifma clflushopt clwb intel_pt avx512cd sha_ni avx512bw avx512vl xsaveopt xsavec xgetbv1 xsaves split_lock_detect dtherm ida arat pln pts hwp hwp_notify hwp_act_window hwp_epp hwp_pkg_req vnmi avx512vbmi umip pku ospke avx512_vbmi2 gfni vaes vpclmulqdq avx512_vnni avx512_bitalg tme avx512_vpopcntdq rdpid movdiri movdir64b fsrm avx512_vp2intersect md_clear ibt flush_l1d arch_capabilities
vmx flags       : vnmi preemption_timer posted_intr invvpid ept_x_only ept_ad ept_1gb flexpriority apicv tsc_offset vtpr mtf vapic ept vpid unrestricted_guest vapic_reg vid ple shadow_vmcs pml ept_mode_based_exec tsc_scaling
bugs            : apic_c1e spectre_v1 spectre_v2 spec_store_bypass swapgs eibrs_pbrsb gds bhi
bogomips        : 4992.00
clflush size    : 64
cache_alignment : 64
address sizes   : 39 bits physical, 48 bits virtual
power management:
```

## Commands
The command used to get the execution time of mine and other implementations is the following:
```bash
time ./test.sh <FORK> measurements_1B.txt
```
> The first time is always omitted as it takes some time to first compile the implementation.

## Comparisons
To compare my solution with others submitted, the below table is created which contains the execution times of mine implementation and others.

The selected implementation are the following:
- The baseline version
- The 3 best Java implementations
- The 2 other existing Go implementations

To take the final execution time results, 3 runs were executed for each implementation, and the average time is used (check `average.py`).

| **Local Result (m:s.ms)** | **Official Result (m:s.ms)** | **Language** | **Implementation** |
| :-----------------------: | :--------------------------: | :----------: | :----------------: |
|          3:35.92          |          04:49.679           |     Java     |      baseline      |
|          1:36.74          |              -               |      Go      |    thanoskoutr     |
|          0:07.26          |              -               |      Go      |        elh         |
|          0:04.66          |              -               |      Go      | AlexanderYastrebov |
|          0:02.09          |          00:01.535           |     Java     |     thomaswue      |
|          0:02.03          |          00:01.587           |     Java     |   artsiomkorzun    |
|          0:02.13          |          00:01.608           |     Java     |      jerrinot      |

## Personal Results
The progression of my implementation in execution times, based on each improvement is the following:
| **Execution Time** |                              **Implementation**                              |                            **Comment**                             |
| :----------------: | :--------------------------------------------------------------------------: | :----------------------------------------------------------------: |
|       ~5min        |                           Simple Go implementation                           |                                                                    |
|       ~10min       |           Worker implementation with goroutines, channel and locks           |           Slower than original due to mutex lock/unlock            |
|       ~13min       |         Break measurements map into shard with locks for each worker         |                      Still slow, due to locks                      |
|       ~10min       |   Calculate partial result measurements per worker and combine in the end    |               Locking is avoided, Parallel execution               |
|      ~1.5min       | Read multiple input lines (chunks) and pass to each worker through a channel | Every worker processes a chunk of input and not one line at a time |

## Future improvements
There is a lot of room for improvements, to get close to the other Go implementations. Some of the already thoughs for improvement are:
- Parse input manually and avoid using the `strings.Split` and `strconv.ParseFloat` functions to achieve better performance
- Play with worker, chunk size and buffer size values to achieve better performance

## Results
All the execution results in details are below:
```bash
./test.sh baseline measurements_1B.txt  213,90s user 5,06s system 101% cpu 3:36,03 total
./test.sh baseline measurements_1B.txt  211,32s user 5,36s system 101% cpu 3:34,03 total
./test.sh baseline measurements_1B.txt  215,32s user 5,13s system 101% cpu 3:37,69 total

./test.sh thomaswue measurements_1B.txt  0,10s user 0,06s system 7% cpu 2,043 total
./test.sh thomaswue measurements_1B.txt  0,12s user 0,04s system 7% cpu 2,199 total
./test.sh thomaswue measurements_1B.txt  0,09s user 0,05s system 7% cpu 2,017 total

./test.sh artsiomkorzun measurements_1B.txt  0,12s user 0,05s system 8% cpu 1,990 total
./test.sh artsiomkorzun measurements_1B.txt  0,10s user 0,03s system 6% cpu 2,000 total
./test.sh artsiomkorzun measurements_1B.txt  0,12s user 0,02s system 6% cpu 2,105 total

./test.sh jerrinot measurements_1B.txt  0,11s user 0,05s system 7% cpu 2,035 total
./test.sh jerrinot measurements_1B.txt  0,13s user 0,03s system 7% cpu 2,171 total
./test.sh jerrinot measurements_1B.txt  0,11s user 0,05s system 7% cpu 2,189 total

./test.sh elh measurements_1B.txt  91,99s user 5,30s system 1328% cpu 7,323 total
./test.sh elh measurements_1B.txt  91,74s user 4,71s system 1333% cpu 7,230 total
./test.sh elh measurements_1B.txt  91,29s user 5,00s system 1330% cpu 7,239 total

./test.sh AlexanderYastrebov measurements_1B.txt  55,85s user 1,51s system 1220% cpu 4,701 total
./test.sh AlexanderYastrebov measurements_1B.txt  55,91s user 1,53s system 1236% cpu 4,644 total
./test.sh AlexanderYastrebov measurements_1B.txt  55,66s user 1,58s system 1236% cpu 4,630 total

./test.sh thanoskoutr measurements_1B.txt  392,33s user 17,93s system 420% cpu 1:37,49 total
./test.sh thanoskoutr measurements_1B.txt  392,04s user 17,69s system 423% cpu 1:36,66 total
./test.sh thanoskoutr measurements_1B.txt  388,85s user 17,89s system 423% cpu 1:36,07 total
```

# References
Some blogs and repos with other Go implementation for the 1brc challenge:
- [One Billion Rows Challenge in Golang](https://www.bytesizego.com/blog/one-billion-row-challenge-go)
- [1brc - AlexanderYastrebov](https://github.com/gunnarmorling/1brc/tree/main/src/main/go/AlexanderYastrebov)
- [1brc - elh](https://github.com/gunnarmorling/1brc/tree/main/src/main/go/elh)
