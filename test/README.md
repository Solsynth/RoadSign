# Benchmark

This result is design for test the performance of the roadsign.
Welcome to contribute more tests of others reverse proxy software!

## Platform

All tests are running on my workstation:

```text
                     ..'          littlesheep@LittleSheepdeMacBook-Pro
                 ,xNMM.           ------------------------------------
               .OMMMMo            OS: macOS Sonoma 14.1 23B2073 arm64
               lMM"               Host: MacBook Pro (14-inch, Nov 2023, Three Thunderbolt 4 ports)
     .;loddo:.  .olloddol;.       Kernel: 23.1.0
   cKMMMMMMMMMMNWMMMMMMMMMM0:     Uptime: 2 days, 1 hour, 16 mins
 .KMMMMMMMMMMMMMMMMMMMMMMMWd.     Packages: 63 (brew), 4 (brew-cask)
 XMMMMMMMMMMMMMMMMMMMMMMMX.       Shell: zsh 5.9
;MMMMMMMMMMMMMMMMMMMMMMMM:        Display (Color LCD): 3024x1964 @ 120Hz (as 1512x982) [Built-in]
:MMMMMMMMMMMMMMMMMMMMMMMM:        DE: Aqua
.MMMMMMMMMMMMMMMMMMMMMMMMX.       WM: Quartz Compositor
 kMMMMMMMMMMMMMMMMMMMMMMMMWd.     WM Theme: Multicolor (Dark)
 'XMMMMMMMMMMMMMMMMMMMMMMMMMMk    Font: .AppleSystemUIFont [System], Helvetica [User]
  'XMMMMMMMMMMMMMMMMMMMMMMMMK.    Cursor: Fill - Black, Outline - White (32px)
    kMMMMMMMMMMMMMMMMMMMMMMd      Terminal: iTerm 3.4.22
     ;KMMMMMMMWXXWMMMMMMMk.       Terminal Font: MesloLGMNFM-Regular (12pt)
       "cooc*"    "*coo'"         CPU: Apple M3 Max (14) @ 4.06 GHz
                                  GPU: Apple M3 Max (30) [Integrated]
                                  Memory: 18.45 GiB / 36.00 GiB (51%)
                                  Swap: Disabled
                                  Disk (/): 72.52 GiB / 926.35 GiB (8%) - apfs [Read-only]
                                  Local IP (en0): 192.168.50.0/24 *
                                  Battery: 100% [AC connected]
                                  Power Adapter: 96W USB-C Power Adapter
                                  Locale: zh_CN.UTF-8
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
# => 8905230 requests in 1m56.215762709s, 3.52GB read
# => Requests/sec:		76626.70
# => Transfer/sec:		30.98MB
# => Avg Req Time:		783.016µs
# => Fastest Request:	28.542µs
# => Slowest Request:	46.773083ms
# => Number of Errors:	0
```

### RoadSign w/ Prefork

```shell
go-wrk -c 60 -d 120 http://localhost:8000
# => Running 120s test @ http://localhost:8000
# =>  60 goroutine(s) running concurrently
# => 4784308 requests in 1m59.100307178s, 1.89GB read
# => Requests/sec:		40170.41
# => Transfer/sec:		16.24MB
# => Avg Req Time:		1.493636ms
# => Fastest Request:	34.291µs
# => Slowest Request:	8.727666ms
# => Number of Errors:	0
```
