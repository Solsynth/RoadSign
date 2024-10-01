import { Command, type Usage } from "clipanion"
import { RsConfig } from "../utils/config.ts"
import chalk from "chalk"

export class ListServerCommand extends Command {
  static paths = [[`list`], [`ls`]]
  static usage: Usage = {
    category: `Networking`,
    description: `List all connected RoadSign Sideload Services`,
    details: `Listing all servers that already saved in RoadSign CLI configuration file`,
    examples: [["List all", `list`]]
  }

  async execute() {
    const config = await RsConfig.getInstance()

    for (let idx = 0; idx < config.config.servers.length; idx++) {
      const server = config.config.servers[idx]
      this.context.stdout.write(`${idx + 1}. ${chalk.bold(server.label)} ${chalk.gray(`(${server.url})`)}\n`)
    }

    this.context.stdout.write("\n" + chalk.cyan(`Connected ${config.config.servers.length} server(s) in total.`) + "\n")

    process.exit(0)
  }
}