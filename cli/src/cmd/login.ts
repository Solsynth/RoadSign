import { Command, Option, type Usage } from "clipanion"
import { createAuthHeader } from "../utils/auth.ts"
import { RsConfig } from "../utils/config.ts"
import ora, { oraPromise } from "ora"

export class LoginCommand extends Command {
  static paths = [[`login`]]
  static usage: Usage = {
    category: `Networking`,
    description: `Login to RoadSign Sideload Service`,
    details: `Login to RoadSign Server`,
    examples: [["Login with credentials", `login <label> <host> <password>`]]
  }

  label = Option.String({ required: true })
  host = Option.String({ required: true })
  credentials = Option.String({ required: true })

  async execute() {
    const config = await RsConfig.getInstance()
    const spinner = ora(`Connecting to ${this.host}...`).start()

    if (!this.host.includes(":")) {
      this.host += ":81"
    }
    if (!this.host.startsWith("http")) {
      this.host = "http://" + this.host
    }

    try {
      const pingRes = await fetch(`${this.host}/cgi/metadata`, {
        headers: {
          Authorization: createAuthHeader("RoadSign CLI", this.credentials)
        }
      })
      if (pingRes.status !== 200) {
        throw new Error(await pingRes.text())
      } else {
        const info = await pingRes.json()
        spinner.succeed(`Connected to ${this.host}, remote version ${info["version"]}`)

        config.config.servers.push({
          label: this.label,
          url: this.host,
          credential: this.credentials
        })
        await oraPromise(config.writeConfig(), { text: "Saving changes..." })
      }
    } catch (e) {
      spinner.fail(`Unable connect to remote: ${e}`)
    }

    process.exit(0)
  }
}