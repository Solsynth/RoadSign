# 🚦 RoadSign

A blazing fast http server with a lot of shining features.

## Features

1. Reverse proxy
2. WebSocket Support
3. Static File Hosting
4. Low Configuration
5. Analytics and Metrics
6. Integrate with CI/CD
7. One-liner CLI
8. Open-source and free
9. **Blazing fast ⚡**

But, as we said, this isn't a reverse proxy. It's a http server.
So we didn't recommend you directly use it as a gateway to users, because it doesn't support all the features of http.
For example like doesn't support multiple site HTTPS, HTTP/2.0, HTTP/3 and Server-Side Events.

You can use RoadSign behind a reverse proxy like caddy, and let it handle the HTTP/2.0, HTTP/3 and all the other
cutting-edge stuff.

### How fast is it?

We use roadsign and nginx to host the same static file
and test them with [go-wrk](https://github.com/tsliwowicz/go-wrk).
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
    image: xsheep2010/roadsign:nightly
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

After that, you can manage your roadsign instance with RoadSign CLI aka. RDC.
To install it, run this command. (Make sure you have golang toolchain on your computer)

```shell
go install git.solsynth.dev/goatworks/roadsign/pkg/cmd/rdc@latest
```

## Usage

To use roadsign, you need to add a configuration for it. Create a file locally.
Name whatever you like. And follow our [documentation](https://wiki.smartsheep.studio/roadsign/configuration/index.html)
to
write it.

After configure, you need sync your config to remote server. Before that, add a connection between roadsign server and
rds cli with this command.

```shell
rdc connect <id> <url> <password>
# ID will allow you find this server in after commands.
# URL is to your roadsign server sideload api.
# Password is your roadsign server credential.
# ======================================================================
# !WARNING! All these things will storage in your $HOME/.roadsignrc.yaml
# ======================================================================
```

Then, sync your local config to remote.

```shell
rdc sync <server id> <region id> <config file>
# Server ID is your server added by last command.
# Site ID is your new site id or old site id if you need update it.
# Config File is your local config file path.
```

After a few seconds, your website is ready!