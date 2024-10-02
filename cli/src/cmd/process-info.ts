import { Command, Option, type Usage } from "clipanion"
import { RsConfig } from "../utils/config.ts"
import { createAuthHeader } from "../utils/auth.ts"
import Table from "cli-table3"
import chalk from "chalk"
import ora from "ora"

export class ProcessCommand extends Command {
  static paths = [[`process`], [`ps`]]
  static usage: Usage = {
    category: `Networking`,
    description: `Loading the application of RoadSign Server`,
    details: `Fetching the configured things amount and other things of a connected server`,
    examples: [
      ["Fetch app directory from labeled server", `ps <label>`],
      ["Fetch app logs from labeled server", `ps <label> <applicationId> logs`]
    ]
  }

  label = Option.String({ required: true })
  applicationId = Option.String({ required: false })
  subcommand = Option.String({ required: false })
  loop = Option.Boolean("--loop,--follow,-f", false, { description: "Keep updating the results" })

  async execute() {
    const config = await RsConfig.getInstance()

    const server = config.config.servers.find(item => item.label === this.label)
    if (server == null) {
      this.context.stdout.write(chalk.red(`Server with label ${chalk.bold(this.label)} was not found.\n`))
      return
    }

    const spinner = ora(`Fetching stats from server ${this.label}...`).start()
    const prefStart = performance.now()


    if (this.applicationId == null) {
      try {
        const res = await fetch(`${server.url}/cgi/applications`, {
          headers: {
            Authorization: createAuthHeader(server.credential)
          }
        })
        if (res.status !== 200) {
          throw new Error(await res.text())
        }
        const prefTook = performance.now() - prefStart
        if (!this.loop) {
          spinner.succeed(`Fetching completed in ${(prefTook / 1000).toFixed(2)}s 🎉`)
        }

        const table = new Table({
          head: ["ID", "Status", "Command"],
          colWidths: [20, 10, 48]
        })

        const statusMapping = ["Created", "Starting", "Started", "Exited", "Failed"]

        const data = await res.json()
        for (const app of data) {
          table.push([app["id"], statusMapping[app["status"]], app["command"].join(" ")])
        }

        this.context.stdout.write(table.toString())
      } catch (e) {
        spinner.fail(`Server with label ${chalk.bold(this.label)} is not running! 😢`)
        return
      }
    } else {
      switch (this.subcommand) {
        case "logs":
          while (true) {
            try {
              const res = await fetch(`${server.url}/cgi/applications/${this.applicationId}/logs`, {
                headers: {
                  Authorization: createAuthHeader(server.credential)
                }
              })
              if (res.status === 404) {
                spinner.fail(`App with id ${chalk.bold(this.applicationId)} was not found! 😢`)
                return
              }
              if (res.status !== 200) {
                throw new Error(await res.text())
              }
              const prefTook = performance.now() - prefStart
              if (!this.loop) {
                spinner.succeed(`Fetching completed in ${(prefTook / 1000).toFixed(2)}s 🎉`)
              }

              this.context.stdout.write(await res.text())
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
        case "start":
        case "stop":
        case "restart":
          try {
            const res = await fetch(`${server.url}/cgi/applications/${this.applicationId}/${this.subcommand}`, {
              method: "POST",
              headers: {
                Authorization: createAuthHeader(server.credential)
              }
            })
            if (res.status === 404) {
              spinner.fail(`App with id ${chalk.bold(this.applicationId)} was not found! 😢`)
              return
            }
            if (res.status === 500) {
              this.context.stdout.write(chalk.red(`Server failed to perform action for application: ${await res.text()}\n`))
              spinner.fail(`Failed to perform action ${chalk.bold(this.applicationId)}... 😢`)
              return
            }
            if (res.status !== 200) {
              throw new Error(await res.text())
            }
            const prefTook = performance.now() - prefStart
            if (!this.loop) {
              spinner.succeed(`Fetching completed in ${(prefTook / 1000).toFixed(2)}s 🎉`)
            }
          } catch (e) {
            spinner.fail(`Server with label ${chalk.bold(this.label)} is not running! 😢`)
            return
          }
          spinner.succeed(`Action for application ${chalk.bold(this.applicationId)} has been performed. 🎉`)
          break
        default:
          this.context.stdout.write(chalk.red(`Subcommand ${chalk.bold(this.subcommand)} was not found.\n`))
      }
    }

    process.exit(0)
  }
}