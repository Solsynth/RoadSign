import { RsConfig } from "../utils/config.ts"
import { Command, Option, type Usage } from "clipanion"
import chalk from "chalk"
import ora from "ora"
import * as fs from "node:fs"
import * as child_process from "node:child_process"
import * as path from "node:path"
import { createAuthHeader } from "../utils/auth.ts"
import { RsLocalConfig, type RsLocalConfigDeploymentPostActionData } from "../utils/config-local.ts"
import * as os from "node:os"

export class DeployCommand extends Command {
  static paths = [[`deploy`]]
  static usage: Usage = {
    category: `Building`,
    description: `Deploying App / Static Site onto RoadSign`,
    details: `Deploying an application or hosting a static site via RoadSign, you need preconfigured the RoadSign, or sync the configurations via sync command.`,
    examples: [
      ["Deploying to RoadSign", `deploy <server> <region> <site> <file / directory>`],
      ["Deploying to RoadSign with .roadsignrc file", `deploy <server>`]
    ]
  }

  server = Option.String({ required: true })
  region = Option.String({ required: false })
  site = Option.String({ required: false })
  input = Option.String({ required: false })

  async deploy(serverLabel: string, region: string, site: string, input: string, postDeploy: RsLocalConfigDeploymentPostActionData | null = null) {
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

    let isDirectory = false
    if (fs.statSync(input).isDirectory()) {
      const compressPrefStart = performance.now()
      const compressSpinner = ora(`Compressing ${chalk.bold(input)}...`).start()
      const destName = path.join(os.tmpdir(), `${Date.now()}-roadsign-archive.zip`)
      child_process.execSync(`cd ${input} && zip -r ${destName} .`)
      const compressPrefTook = performance.now() - compressPrefStart
      compressSpinner.succeed(`Compressing completed in ${(compressPrefTook / 1000).toFixed(2)}s 🎉`)
      input = destName
      isDirectory = true
    }

    const destBreadcrumb = [region, site].join(" ➜ ")
    const spinner = ora(`Deploying ${chalk.bold(destBreadcrumb)} to ${chalk.bold(this.server)}...`).start()

    const prefStart = performance.now()

    try {
      const payload = new FormData()
      payload.set("attachments", await fs.openAsBlob(input), isDirectory ? "dist.zip" : path.basename(input))

      if(postDeploy) {
        if(postDeploy.command) {
          payload.set("post-deploy-script", postDeploy.command)
        } else if(postDeploy.scriptPath) {
          payload.set("post-deploy-script", fs.readFileSync(postDeploy.scriptPath, "utf8"))
        } else {
          this.context.stdout.write(chalk.yellow(`Configured post deploy action but no script provided, skip performing post deploy action...\n`))
        }
        payload.set("post-deploy-environment", postDeploy.environment?.join("\n") ?? "")
      }

      const res = await fetch(`${server.url}/webhooks/publish/${region}/${site}?mimetype=application/zip`, {
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
      if (isDirectory && input.endsWith(".zip")) {
        fs.unlinkSync(input)
      }
    }
  }

  async execute() {
    if (this.region && this.site && this.input) {
      await this.deploy(this.server, this.region, this.site, this.input)
    } else {
      let localCfg: RsLocalConfig
      try {
        localCfg = await RsLocalConfig.getInstance()
      } catch (e) {
        this.context.stdout.write(chalk.red(`Unable to load .roadsignrc: ${e}\n`))
        return
      }

      if (!localCfg.config.deployments) {
        this.context.stdout.write(chalk.red(`No deployments found in .roadsignrc, exiting...\n`))
        return
      }

      let idx = 0
      for (const deployment of localCfg.config.deployments ?? []) {
        this.context.stdout.write(chalk.cyan(`Deploying ${idx + 1} out of ${localCfg.config.deployments.length} deployments...\n`))
        await this.deploy(this.server, deployment.region, deployment.site, deployment.path)
      }

      this.context.stdout.write(chalk.green(`All deployments has been deployed!\n`))
    }

    process.exit(0)
  }
}