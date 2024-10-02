import { RsConfig } from "../utils/config.ts"
import { Command, Option, type Usage } from "clipanion"
import chalk from "chalk"
import ora from "ora"
import * as fs from "node:fs"
import * as child_process from "node:child_process"
import * as path from "node:path"
import { createAuthHeader } from "../utils/auth.ts"

export class DeployCommand extends Command {
  static paths = [[`deploy`]]
  static usage: Usage = {
    category: `Building`,
    description: `Deploying App / Static Site onto RoadSign`,
    details: `Deploying an application or hosting a static site via RoadSign, you need preconfigured the RoadSign, or sync the configurations via sync command.`,
    examples: [["Deploying to RoadSign", `deploy <server> <site> <slug> <file / directory>`]]
  }

  server = Option.String({ required: true })
  site = Option.String({ required: true })
  upstream = Option.String({ required: true })
  input = Option.String({ required: true })

  async execute() {
    const config = await RsConfig.getInstance()

    const server = config.config.servers.find(item => item.label === this.server)
    if (server == null) {
      this.context.stdout.write(chalk.red(`Server with label ${chalk.bold(this.server)} was not found.\n`))
      return
    }

    if (!fs.existsSync(this.input)) {
      this.context.stdout.write(chalk.red(`Input file ${chalk.bold(this.input)} was not found.\n`))
      return
    }

    let isDirectory = false
    if (fs.statSync(this.input).isDirectory()) {
      if (this.input.endsWith("/")) {
        this.input = this.input.slice(0, -1)
      }
      this.input += "/*"

      const compressPrefStart = performance.now()
      const compressSpinner = ora(`Compressing ${chalk.bold(this.input)}...`).start()
      const destName = `${Date.now()}-roadsign-archive.zip`
      child_process.execSync(`zip -rj ${destName} ${this.input}`)
      const compressPrefTook = performance.now() - compressPrefStart
      compressSpinner.succeed(`Compressing completed in ${(compressPrefTook / 1000).toFixed(2)}s 🎉`)
      this.input = destName
      isDirectory = true
    }

    const destBreadcrumb = [this.site, this.upstream].join(" ➜ ")
    const spinner = ora(`Deploying ${chalk.bold(destBreadcrumb)} to ${chalk.bold(this.server)}...`).start()

    const prefStart = performance.now()

    try {
      const payload = new FormData()
      payload.set("attachments", await fs.openAsBlob(this.input), isDirectory ? "dist.zip" : path.basename(this.input))
      const res = await fetch(`${server.url}/webhooks/publish/${this.site}/${this.upstream}?mimetype=application/zip`, {
        method: "PUT",
        body: payload,
        headers: {
          Authorization: createAuthHeader(server.credential)
        }
      })
      if (res.status !== 200) {
        throw new Error(await res.text())
      }
      const prefTook = performance.now() - prefStart
      spinner.succeed(`Deploying completed in ${(prefTook / 1000).toFixed(2)}s 🎉`)
    } catch (e) {
      this.context.stdout.write(`Failed to deploy to remote: ${e}\n`)
      spinner.fail(`Server with label ${chalk.bold(this.server)} is not running! 😢`)
    } finally {
      if (isDirectory && this.input.endsWith(".zip")) {
        fs.unlinkSync(this.input)
      }
    }

    process.exit(0)
  }
}