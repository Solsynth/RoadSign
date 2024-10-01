import { Cli } from "clipanion"
import figlet from "figlet"
import chalk from "chalk"

import { LoginCommand } from "./src/cmd/login.ts"
import { LogoutCommand } from "./src/cmd/logout.ts"
import { ListServerCommand } from "./src/cmd/list.ts"

const [node, app, ...args] = process.argv

console.log(
  chalk.yellow(figlet.textSync("RoadSign CLI", { horizontalLayout: "default", verticalLayout: "default" }))
)

const cli = new Cli({
  binaryLabel: `RoadSign CLI`,
  binaryName: `${node} ${app}`,
  binaryVersion: `1.0.0`
})

cli.register(LoginCommand)
cli.register(LogoutCommand)
cli.register(ListServerCommand)
cli.runExit(args)