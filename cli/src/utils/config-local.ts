import * as path from "node:path"
import * as fs from "node:fs/promises"

interface RsLocalConfigData {
  sync?: RsLocalConfigSyncData
  deployments?: RsLocalConfigDeploymentData[]
}

interface RsLocalConfigSyncData {
  configPath: string
  region: string
}

interface RsLocalConfigDeploymentData {
  path: string
  region: string
  site: string
  autoBuild?: RsLocalConfigDeploymentAutoBuildData
}

interface RsLocalConfigDeploymentAutoBuildData {
  command: string
  environment?: string[]
}

class RsLocalConfig {
  private static instance: RsLocalConfig

  public config: RsLocalConfigData = {}

  private constructor() {
  }

  public static async getInstance(): Promise<RsLocalConfig> {
    if (!RsLocalConfig.instance) {
      RsLocalConfig.instance = new RsLocalConfig()
      await RsLocalConfig.instance.readConfig()
    }
    return RsLocalConfig.instance
  }

  public async readConfig() {
    const basepath = process.cwd()
    const filepath = path.join(basepath, ".roadsignrc")
    if (!await fs.exists(filepath)) {
      throw new Error(`.roadsignrc file was not found at ${filepath}`)
    }

    const data = await fs.readFile(filepath, "utf8")
    this.config = JSON.parse(data)
  }

  public async writeConfig() {
    const basepath = process.cwd()
    const filepath = path.join(basepath, ".roadsignrc")
    await fs.writeFile(filepath, JSON.stringify(this.config))
  }
}

export { RsLocalConfig, type RsLocalConfigData }