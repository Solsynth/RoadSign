import { Command, Option, type Usage } from "clipanion"
import { RsConfig, type RsConfigServerData } from "../utils/config.ts"
import { createAuthHeader } from "../utils/auth.ts"
import chalk from "chalk"
import ora from "ora"

export class InfoCommand extends Command {
  static paths = [[`info`], [`if`]]
  static usage: Usage = {
    category: `Networking`,
    description: `Fetching the stats of RoadSign Server`,
    details: `Fetching the configured things amount and other things of a connected server`,
    examples: [["Fetch stats from labeled server", `info <label> [area]`]]
  }

  label = Option.String({ required: true })
  area = Option.String({ required: false })
  loop = Option.Boolean("--loop,--follow,-f", false, { description: "Keep updating the results" })

  private static formatUptime(ms: number): string {
    let seconds: number = Math.floor(ms / 1000)
    let minutes: number = Math.floor(seconds / 60)
    let hours: number = Math.floor(minutes / 60)
    let days: number = Math.floor(hours / 24)

    seconds = seconds % 60
    minutes = minutes % 60
    hours = hours % 24

    const uptimeParts: string[] = []

    if (days > 0) uptimeParts.push(`${days} day${days > 1 ? "s" : ""}`)
    if (hours > 0) uptimeParts.push(`${hours} hour${hours > 1 ? "s" : ""}`)
    if (minutes > 0) uptimeParts.push(`${minutes} minute${minutes > 1 ? "s" : ""}`)
    if (seconds > 0 || uptimeParts.length === 0) uptimeParts.push(`${seconds} second${seconds > 1 ? "s" : ""}`)

    return uptimeParts.join(", ")
  }

  async fetchOverview(server: RsConfigServerData) {
    try {
      const res = await fetch(`${server.url}/cgi/stats`, {
        headers: {
          Authorization: createAuthHeader(server.credential)
        }
      })
      if (res.status !== 200) {
        throw new Error(await res.text())
      }

      const data = await res.json()
      this.context.stdout.write(`\nServer stats of ${chalk.bold(this.label)}\n`)
      this.context.stdout.write(`Uptime: ${chalk.bold(InfoCommand.formatUptime(data["uptime"]))}\n`)
      this.context.stdout.write(`Traffic since last startup: ${chalk.bold(data["traffic"]["total"])}\n`)
      this.context.stdout.write(`Unique clients since last startup: ${chalk.bold(data["traffic"]["unique_client"])}\n`)
      this.context.stdout.write(`\nServer info of ${chalk.bold(this.label)}\n`)
      this.context.stdout.write(`Warden Applications: ${chalk.bold(data["applications"])}\n`)
      this.context.stdout.write(`Destinations: ${chalk.bold(data["destinations"])}\n`)
      this.context.stdout.write(`Locations: ${chalk.bold(data["locations"])}\n`)
      this.context.stdout.write(`Regions: ${chalk.bold(data["regions"])}\n`)
    } catch (e) {
      return
    }
  }

  async fetchTrace(server: RsConfigServerData) {
    const res = await fetch(`${server.url}/cgi/traces`, {
      headers: {
        Authorization: createAuthHeader(server.credential)
      }
    })
    if (res.status !== 200) {
      throw new Error(await res.text())
    }

    const data = await res.json()
    for (const trace of data) {
      const ts = new Date(trace["timestamp"]).toLocaleString()
      const path = [trace["region"], trace["location"], trace["destination"]].join(" ➜ ")
      const uri = trace["uri"].split("?").length == 1 ? trace["uri"] : trace["uri"].split("?")[0] + ` ${chalk.grey(`w/ query parameters`)}`
      this.context.stdout.write(`${chalk.bgGrey(`[${ts}]`)} ${chalk.bold(path)} ${chalk.cyan(trace["ip_address"])} ${uri}\n`)
    }
  }

  async fetchRegions(server: RsConfigServerData) {
    const res = await fetch(`${server.url}/cgi/regions`, {
      headers: {
        Authorization: createAuthHeader(server.credential)
      }
    })
    if (res.status !== 200) {
      throw new Error(await res.text())
    }

    const data = await res.json()
    this.context.stdout.write("\n\n")
    for (const region of data) {
      this.context.stdout.write(` • ${chalk.bgGrey('region#')}${chalk.bold(region.id)} ${chalk.gray(`(${region.locations.length} locations)`)}\n`)
      for (const location of region.locations) {
        this.context.stdout.write(`   • ${chalk.bgGrey('location#')} ${chalk.bold(location.id)} ${chalk.gray(`(${location.destinations.length} destinations)`)}\n`)
        for (const destination of location.destinations) {
          this.context.stdout.write(`     • ${chalk.bgGrey('destination#')}${chalk.bold(destination.id)}\n`)
        }
      }
      this.context.stdout.write("\n")
    }
  }

  async execute() {
    const config = await RsConfig.getInstance()

    const server = config.config.servers.find(item => item.label === this.label)
    if (server == null) {
      this.context.stdout.write(chalk.red(`Server with label ${chalk.bold(this.label)} was not found.\n`))
      return
    }

    if (this.area == null) {
      this.area = "overview"
    }

    const spinner = ora(`Fetching stats from server ${this.label}...`).start()
    const prefStart = performance.now()

    switch (this.area) {
      case "overview":
        try {
          await this.fetchOverview(server)
          const prefTook = performance.now() - prefStart
          spinner.succeed(`Fetching completed in ${(prefTook / 1000).toFixed(2)}s 🎉`)
        } catch (e) {
          spinner.fail(`Server with label ${chalk.bold(this.label)} is not running! 😢`)
        }
        break
      case "trace":
        while (true) {
          try {
            await this.fetchTrace(server)
            const prefTook = performance.now() - prefStart
            if (!this.loop) {
              spinner.succeed(`Fetching completed in ${(prefTook / 1000).toFixed(2)}s 🎉`)
            }
          } catch (e) {
            spinner.fail(`Server with label ${chalk.bold(this.label)} is not running! 😢`)
            return
          }

          if (!this.loop) {
            break
          } else {
            spinner.text = "Updating..."
            await new Promise(resolve => setTimeout(resolve, 3000))
            this.context.stdout.write("\x1Bc")
          }
        }
        break
      case "regions":
        try {
          await this.fetchRegions(server)
          const prefTook = performance.now() - prefStart
          spinner.succeed(`Fetching completed in ${(prefTook / 1000).toFixed(2)}s 🎉`)
        } catch (e) {
          spinner.fail(`Server with label ${chalk.bold(this.label)} is not running! 😢`)
          return
        }
        break
      default:
        spinner.fail(chalk.red(`Info area was not exists ${chalk.bold(this.area)}...`))
    }

    process.exit(0)
  }
}