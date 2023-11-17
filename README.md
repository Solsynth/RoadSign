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
Total time passed:              8.96s
Avg time per request:           88.63ms
Requests per second:            1115.70
Median time per request:        85.76ms
99th percentile time:           144.05ms
Slowest time for request:       237.00ms

=============================DATA=============================
Total response body sizes:              5431892
Avg response body per request:          543.19 Byte
Transfer rate per second:               606034.29 Byte/s (0.61 MByte/s)
==========================RESPONSES==========================
20X Responses:          277     (2.77%)
30X Responses:          0       (0.00%)
40X Responses:          9723    (97.23%)
50X Responses:          0       (0.00%)
Errors:                 0       (0.00%)
```