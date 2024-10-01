import { Command, Option, type Usage } from "clipanion"
import { RsConfig } from "../utils/config.ts"
import { oraPromise } from "ora"
import chalk from "chalk"

export class LogoutCommand extends Command {
  static paths = [[`logout`]]
  static usage: Usage = {
    category: `Networking`,
    description: `Logout from RoadSign Sideload Service`,
    details: `Logout from RoadSign Server`,
    examples: [["Logout with server label", `logout <label>`]]
  }

  label = Option.String({ required: true })

  async execute() {
    const config = await RsConfig.getInstance()

    const server = config.config.servers.findIndex(item => item.label === this.label)
    if (server === -1) {
      this.context.stdout.write(chalk.red(`Server with label ${chalk.bold(this.label)} was not found.\n`))
    } else {
      config.config.servers.splice(server, 1)
      this.context.stdout.write(chalk.green(`Server with label ${chalk.bold(this.label)} was successfully removed.\n`))
      await oraPromise(config.writeConfig(), { text: "Saving changes..." })
    }

    process.exit(0)
  }
}