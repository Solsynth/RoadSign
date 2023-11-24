# Benchmark

This result is design for test the performance of the roadsign.
Welcome to contribute more tests of others reverse proxy software!

## Platform

All tests are running on my workstation:

```text
                                ..,   LittleSheep@PEROPERO-WORKSTATION
                    ....,,:;+ccllll   --------------------------------
      ...,,+:;  cllllllllllllllllll   OS: Windows 10 Pro x86_64
,cclllllllllll  lllllllllllllllllll   Host: LENOVO 82TF
llllllllllllll  lllllllllllllllllll   Kernel: 10.0.19045
llllllllllllll  lllllllllllllllllll   Uptime: 1 hour, 22 mins
llllllllllllll  lllllllllllllllllll   Shell: pwsh 7.4.0
llllllllllllll  lllllllllllllllllll   Resolution: 2560x1600
llllllllllllll  lllllllllllllllllll   DE: Aero
                                      WM: Explorer
llllllllllllll  lllllllllllllllllll   WM Theme: Custom
llllllllllllll  lllllllllllllllllll   Terminal: Windows Terminal
llllllllllllll  lllllllllllllllllll   CPU: 12th Gen Intel i7-12700H (20) @ 2.690GHz
llllllllllllll  lllllllllllllllllll   GPU: Caption
llllllllllllll  lllllllllllllllllll   GPU: NVIDIA GeForce RTX 3070 Laptop GPU
`'ccllllllllll  lllllllllllllllllll   GPU
       `' \*::  :ccllllllllllllllll   Memory: 7318MiB / 16192MiB
                       ````''*::cll
                                 ``
```

## Results

The tests are run in the order `nginx -> roadsign without prefork -> roadsign with prefork`. There is no reason why nginx performance should be affected by hardware temperature.

### Nginx

```shell
go-wrk -c 60 -d 120 http://localhost:8001
# => Running 120s test @ http://localhost:8001
# =>   60 goroutine(s) running concurrently
# => 515749 requests in 1m59.953302003s, 245.92MB read
# => Requests/sec:           4299.58
# => Transfer/sec:           2.05MB
# => Avg Req Time:           13.954846ms
# => Fastest Request:        0s
# => Slowest Request:        410.6972ms
# => Number of Errors:       0
```

### RoadSign

```shell
go-wrk -c 60 -d 120 http://localhost:8000
# => Running 120s test @ http://localhost:8000
# =>   60 goroutine(s) running concurrently
# => 3256820 requests in 1m59.446620043s, 1.43GB read
# => Requests/sec:           27265.90
# => Transfer/sec:           12.27MB
# => Avg Req Time:           2.20055ms
# => Fastest Request:        0s
# => Slowest Request:        56.8726ms
# => Number of Errors:       0
```

### RoadSign w/ Prefork

```shell
go-wrk -c 60 -d 120 http://localhost:8000
# => Running 120s test @ http://localhost:8000
# =>   60 goroutine(s) running concurrently
# => 2188594 requests in 1m59.933175915s, 985.16MB read
# => Requests/sec:           18248.45
# => Transfer/sec:           8.21MB
# => Avg Req Time:           3.287951ms
# => Fastest Request:        0s
# => Slowest Request:        121.5189ms
# => Number of Errors:       0
```