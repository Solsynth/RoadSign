import resolve from "@rollup/plugin-node-resolve"
import commonjs from "@rollup/plugin-commonjs"
import typescript from "@rollup/plugin-typescript"
import json from "@rollup/plugin-json"

export default {
  input: "index.ts",
  output: {
    banner: "#!/usr/bin/env node",
    file: "dist/index.cjs",
    format: "cjs",
    inlineDynamicImports: true,
  },
  plugins: [
    resolve(),
    commonjs(),
    json(),
    typescript({
      tsconfig: "./tsconfig.json"
    })
  ],
}