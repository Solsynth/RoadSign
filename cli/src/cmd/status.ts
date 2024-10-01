import { Command, Option, type Usage } from "clipanion"
import { RsConfig } from "../utils/config.ts"
import { createAuthHeader } from "../utils/auth.ts"
import chalk from "chalk"
import ora from "ora"

export class StatusCommand extends Command {
  static paths = [[`status`]]
  static usage: Usage = {
    category: `Networking`,
    description: `Check the status of RoadSign Sideload Service`,
    details: `Check the running status of a connected server`,
    examples: [["Check the status of labeled server", `status <label>`]]
  }

  label = Option.String({ required: true })

  async execute() {
    const config = await RsConfig.getInstance()

    const server = config.config.servers.find(item => item.label === this.label)
    if (server == null) {
      this.context.stdout.write(chalk.red(`Server with label ${chalk.bold(this.label)} was not found.\n`))
      return
    }

    const spinner = ora(`Checking status of ${this.label}...`).start()

    try {
      const res = await fetch(`${server.url}/cgi/metadata`, {
        headers: {
          Authorization: createAuthHeader(server.credential)
        }
      })
      if (res.status !== 200) {
        throw new Error(await res.text())
      }
      spinner.succeed(`Server with label ${chalk.bold(this.label)} is up and running! 🎉`)
    } catch (e) {
      spinner.fail(`Server with label ${chalk.bold(this.label)} is not running! 😢`)
      return
    }

    process.exit(0)
  }
}