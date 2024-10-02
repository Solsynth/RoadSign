import { RsConfig } from "../utils/config.ts"
import { Command, Option, type Usage } from "clipanion"
import chalk from "chalk"
import ora from "ora"
import * as fs from "node:fs"
import { createAuthHeader } from "../utils/auth.ts"
import { RsLocalConfig } from "../utils/config-local.ts"

export class ReloadCommand extends Command {
  static paths = [[`reload`]]
  static usage: Usage = {
    category: `Building`,
    description: `Reload configuration on RoadSign`,
    details: `Reload configuration on remote RoadSign to make changes applied.`,
    examples: [
      ["Reload an connected server", `reload <server>`],
    ]
  }

  server = Option.String({ required: true })

  async execute() {
    const cfg = await RsConfig.getInstance()
    const server = cfg.config.servers.find(item => item.label === this.server)
    if (server == null) {
      this.context.stdout.write(chalk.red(`Server with label ${chalk.bold(this.server)} was not found.\n`))
      return
    }

    const spinner = ora(`Reloading server ${chalk.bold(this.server)}...`).start()

    const prefStart = performance.now()

    try {
      const res = await fetch(`${server.url}/cgi/reload`, {
        method: "POST",
        headers: {
          Authorization: createAuthHeader(server.credential)
        }
      })
      if (res.status !== 200) {
        throw new Error(await res.text())
      }
      const prefTook = performance.now() - prefStart
      spinner.succeed(`Reloading completed in ${(prefTook / 1000).toFixed(2)}s 🎉`)
    } catch (e) {
      this.context.stdout.write(`Failed to reload remote: ${e}\n`)
      spinner.fail(`Server with label ${chalk.bold(this.server)} is not running! 😢`)
    }

    process.exit(0)
  }
}