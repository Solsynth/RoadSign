# 🚦 RoadSign

A blazing fast reverse proxy with a lot of shining features.

## Features

1. Reverse proxy
2. Static file hosting
3. ~~Analytics and Metrics~~
4. Integrate with CI/CD
5. Webhook integration
6. ~~Web management panel~~
7. **Blazing fast ⚡**

> Deleted item means under construction, check out our roadmap!

### How fast is it?

We use roadsign and nginx to host a same static file, and test them with [go-wrk](https://github.com/tsliwowicz/go-wrk). 
Here's the result:

|      **Software**     | Total Requests | Requests per Seconds | Transfer per Seconds |   Avg Time  | Fastest Time | Slowest Time | Errors Count |
|:---------------------:|----------------|:--------------------:|:--------------------:|:-----------:|:------------:|:------------:|:------------:|
|        _Nginx_        |     515749     |        4299.58       |        2.05MB        | 13.954846ms |      0s      |  410.6972ms  |       0      |
|       _RoadSign_      |     3256820    |       27265.90       |        12.27MB       |  2.20055ms  |      0s      |   56.8726ms  |       0      |
| _RoadSign w/ Prefork_ |     2188594    |       18248.45       |        8.21MB        |  3.287951ms |      0s      |  121.5189ms  |       0      |

As result, roadsign undoubtedly is the fastest one.

More details can be found at benchmark's [README.md](./test/benchmark/README.md)