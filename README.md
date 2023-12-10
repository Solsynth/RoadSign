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
|:---------------------:|:--------------:|:--------------------:|:--------------------:|:-----------:|:------------:|:------------:|:------------:|
|        _Nginx_        |     515749     |        4299.58       |        2.05MB        | 13.954846ms |      0s (Cached)      |  410.6972ms  |       0      |
|       _RoadSign_      |     8905230    |       76626.70       | 30.98MB       |  783.016µs  |      28.542µs      |   46.773083ms  |       0      |
| _RoadSign w/ Prefork_ |     4784308    |       40170.41       |        16.24MB        | 1.493636ms |      34.291µs      |  8.727666ms  |       0      |

As result, roadsign undoubtedly is the fastest one.

It can be found that the prefork feature makes RoadSign more stable in concurrency. We can see this from the **Slowest Time**. At the same time, the **Fastest Time** is affected because reusing ports requires some extra steps to handle load balancing. Enable this feature at your own discretion depending on your use case.

More details can be found at benchmark's [README.md](./test/README.md)