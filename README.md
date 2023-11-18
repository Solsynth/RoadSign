# 🚦 RoadSign

A blazing fast reverse proxy with a lot of shining features.

## Features

1. Reverse proxy
2. Static file hosting
3. ~~Analytics and Metrics~~
4. Integrate with CI/CD
5. ~~Webhook integration~~
6. ~~Web management panel~~
7. **Blazing fast ⚡**

> Deleted item means under construction, check out our roadmap!

### How fast is it?

Static Files Hosting

```shell
go-wrk -t=8 -c=100 -n=10000 "http://localhost"
```

```text
==========================BENCHMARK==========================
URL:                            http://localhost

Used Connections:               100
Used Threads:                   8
Total number of calls:          10000

===========================TIMINGS===========================
Total time passed:              11.36s
Avg time per request:           112.69ms
Requests per second:            880.32
Median time per request:        111.14ms
99th percentile time:           160.88ms
Slowest time for request:       217.00ms

=============================DATA=============================
Total response body sizes:              190130000
Avg response body per request:          19013.00 Byte
Transfer rate per second:               16737517.73 Byte/s (16.74 MByte/s)
==========================RESPONSES==========================
20X Responses:          10000   (100.00%)
30X Responses:          0       (0.00%)
40X Responses:          0       (0.00%)
50X Responses:          0       (0.00%)
Errors:                 0       (0.00%)
```