# 🚦 RoadSign

A blazing fast reverse proxy with a lot of shining features.

## Features

1. Reverse proxy
2. Static file hosting
3. ~~Analytics and Metrics~~
4. Integrate with CI/CD
5. Webhook integration
6. ~~Web management panel~~
7. One-liner CLI
8. **Blazing fast ⚡**

> Deleted item means under construction, check out our roadmap!

### How fast is it?

We use roadsign and nginx to host a same static file, and test them with [go-wrk](https://github.com/tsliwowicz/go-wrk).
Here's the result:

|     **Software**      | Total Requests | Requests per Seconds | Transfer per Seconds |  Avg Time   | Fastest Time | Slowest Time | Errors Count |
|:---------------------:|:--------------:|:--------------------:|:--------------------:|:-----------:|:------------:|:------------:|:------------:|
|        _Nginx_        |     515749     |       4299.58        |        2.05MB        | 13.954846ms | 0s (Cached)  |  410.6972ms  |      0       |
|      _RoadSign_       |    8905230     |       76626.70       |       30.98MB        |  783.016µs  |   28.542µs   | 46.773083ms  |      0       |
| _RoadSign w/ Prefork_ |    4784308     |       40170.41       |       16.24MB        | 1.493636ms  |   34.291µs   |  8.727666ms  |      0       |

As result, roadsign undoubtedly is the fastest one.

It can be found that the prefork feature makes RoadSign more stable in concurrency. We can see this from the **Slowest
Time**. At the same time, the **Fastest Time** is affected because reusing ports requires some extra steps to handle
load balancing. Enable this feature at your own discretion depending on your use case.

More details can be found at benchmark's [README.md](./test/README.md)

## Installation

We strongly recommend you install RoadSign via docker compose.

```yaml
version: "3"
services:
  roadsign:
    image: code.smartsheep.studio/goatworks/roadsign:nightly
    restart: always
    volumes:
      - "./certs:/certs" # Optional, use for storage certificates
      - "./config:/config"
      - "./wwwroot:/wwwroot" # Optional, use for storage web apps
      - "./settings.yml:/settings.yml"
    ports:
      - "80:80"
      - "443:443"
      - "81:81"
```

After that, you can manage your roadsign instance with RoadSign CLI aka. RDS CLI.
To install it, run this command. (Make sure you have golang toolchain on your computer)

```shell
go install code.smartsheep.studio/goatworks/roadsign/pkg/cmd/rds@latest
```

## Usage

To use roadsign, you need to add a configuration for it. Create a file locally.
Name whatever you like. And follow our [documentation](https://wiki.smartsheep.studio/roadsign/02-configuration.html) to
write it.

After configure, you need sync your config to remote server. Before that, add a connection between roadsign server and
rds cli with this command.

```shell
rds connect <id> <url> <password>
# ID will allow you find this server in after commands.
# URL is to your roadsign server sideload api.
# Password is your roadsign server credential.
# ======================================================================
# !WARNING! All these things will storage in your $HOME/.roadsignrc.yaml
# ======================================================================
```

Then, sync your local config to remote.

```shell
rds sync <server id> <site id> <config file>
# Server ID is your server added by last command.
# Site ID is your new site id or old site id if you need update it.
# Config File is your local config file path.
```

After a few seconds, your website is ready!