import { RsConfig } from "../utils/config.ts"
import { Command, Option, type Usage } from "clipanion"
import chalk from "chalk"
import ora from "ora"
import * as fs from "node:fs"
import { createAuthHeader } from "../utils/auth.ts"
import { RsLocalConfig } from "../utils/config-local.ts"

export class SyncCommand extends Command {
  static paths = [[`sync`]]
  static usage: Usage = {
    category: `Building`,
    description: `Sync configuration to RoadSign over Sideload`,
    details: `Update remote RoadSign configuration with local ones.`,
    examples: [
      ["Sync to RoadSign", `sync <server> <region> <file>`],
      ["Sync to RoadSign with .roadsignrc file", `sync <server>`]
    ]
  }

  server = Option.String({ required: true })
  region = Option.String({ required: false })
  input = Option.String({ required: false })

  async sync(serverLabel: string, region: string, input: string) {
    const cfg = await RsConfig.getInstance()
    const server = cfg.config.servers.find(item => item.label === serverLabel)
    if (server == null) {
      this.context.stdout.write(chalk.red(`Server with label ${chalk.bold(this.server)} was not found.\n`))
      return
    }

    if (!fs.existsSync(input)) {
      this.context.stdout.write(chalk.red(`Input file ${chalk.bold(this.input)} was not found.\n`))
      return
    }
    if (!fs.statSync(input).isFile()) {
      this.context.stdout.write(chalk.red(`Input file ${chalk.bold(this.input)} is not a file.\n`))
      return
    }

    const spinner = ora(`Syncing ${chalk.bold(region)} to ${chalk.bold(this.server)}...`).start()

    const prefStart = performance.now()

    try {
      const res = await fetch(`${server.url}/webhooks/sync/${region}`, {
        method: "PUT",
        body: fs.readFileSync(input, "utf8"),
        headers: {
          Authorization: createAuthHeader(server.credential)
        }
      })
      if (res.status !== 200) {
        throw new Error(await res.text())
      }
      const prefTook = performance.now() - prefStart
      spinner.succeed(`Syncing completed in ${(prefTook / 1000).toFixed(2)}s 🎉`)
    } catch (e) {
      this.context.stdout.write(`Failed to sync to remote: ${e}\n`)
      spinner.fail(`Server with label ${chalk.bold(this.server)} is not running! 😢`)
    }
  }

  async execute() {
    if (this.region && this.input) {
      await this.sync(this.server, this.region, this.input)
    } else {
      let localCfg: RsLocalConfig
      try {
        localCfg = await RsLocalConfig.getInstance()
      } catch (e) {
        this.context.stdout.write(chalk.red(`Unable to load .roadsignrc: ${e}\n`))
        return
      }

      if (!localCfg.config.sync) {
        this.context.stdout.write(chalk.red(`No sync configuration found in .roadsignrc, exiting...\n`))
        return
      }

      await this.sync(this.server, localCfg.config.sync.region, localCfg.config.sync.configPath)
    }

    process.exit(0)
  }
}