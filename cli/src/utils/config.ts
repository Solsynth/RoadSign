import * as os from "node:os"
import * as path from "node:path"
import * as fs from "node:fs/promises"

interface RsConfigData {
  servers: RsConfigServerData[]
}

interface RsConfigServerData {
  label: string
  url: string
  credential: string
}

class RsConfig {
  private static instance: RsConfig

  public config: RsConfigData = {
    servers: []
  }

  private constructor() {
  }

  public static async getInstance(): Promise<RsConfig> {
    if (!RsConfig.instance) {
      RsConfig.instance = new RsConfig()
      await RsConfig.instance.readConfig()
    }
    return RsConfig.instance
  }

  public async readConfig() {
    const basepath = os.homedir()
    const filepath = path.join(basepath, ".roadsignrc")
    if (!await fs.exists(filepath)) {
      await fs.writeFile(filepath, JSON.stringify(this.config))
    }

    const data = await fs.readFile(filepath, "utf8")
    this.config = JSON.parse(data)
  }

  public async writeConfig() {
    const basepath = os.homedir()
    const filepath = path.join(basepath, ".roadsignrc")
    await fs.writeFile(filepath, JSON.stringify(this.config))
  }
}

export { RsConfig, type RsConfigData, type RsConfigServerData }