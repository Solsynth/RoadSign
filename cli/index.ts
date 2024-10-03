import { Builtins, Cli } from "clipanion"
import figlet from "figlet"
import chalk from "chalk"

import { LoginCommand } from "./src/cmd/login.ts"
import { LogoutCommand } from "./src/cmd/logout.ts"
import { ListServerCommand } from "./src/cmd/list.ts"
import { StatusCommand } from "./src/cmd/status.ts"
import { InfoCommand } from "./src/cmd/info.ts"
import { ProcessCommand } from "./src/cmd/process-info.ts"
import { DeployCommand } from "./src/cmd/deploy.ts"
import { SyncCommand } from "./src/cmd/sync.ts"
import { ReloadCommand } from "./src/cmd/reload.ts"

const [node, app, ...args] = process.argv

const ENABLE_STARTUP_ASCII_ART = false

if (process.env["ENABLE_STARTUP_ASCII_ART"] || ENABLE_STARTUP_ASCII_ART) {
  console.log(
    chalk.yellow(figlet.textSync("RoadSign CLI", { horizontalLayout: "default", verticalLayout: "default" }))
  )
}

const cli = new Cli({
  binaryLabel: `RoadSign CLI`,
  binaryName: `${node} ${app}`,
  binaryVersion: `1.0.0`
})

cli.register(Builtins.VersionCommand)
cli.register(Builtins.HelpCommand)
cli.register(LoginCommand)
cli.register(LogoutCommand)
cli.register(ListServerCommand)
cli.register(StatusCommand)
cli.register(InfoCommand)
cli.register(ProcessCommand)
cli.register(DeployCommand)
cli.register(SyncCommand)
cli.register(ReloadCommand)
cli.runExit(args)