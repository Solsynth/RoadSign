<template>
  <div class="flex flex-col gap-1">

    <div>
      <div class="font-bold">Rules</div>
      <n-code :hljs="hljs" :code="parseData(props.rules)" language="json" />
    </div>

    <div>
      <div class="font-bold">Upstreams</div>
      <n-code :hljs="hljs" :code="parseData(props.upstreams)" language="json" />
    </div>

    <div>
      <div class="font-bold">Processes</div>
      <n-code :hljs="hljs" :code="parseData(props.processes)" language="json" />
    </div>

  </div>
</template>

<script setup lang="ts">
import { NCode } from "naive-ui"
import hljs from "highlight.js/lib/core"
import json from "highlight.js/lib/languages/json"

hljs.registerLanguage("json", json)

const props = defineProps<{ rules: any[], upstreams: any[], processes: any[] }>()

function parseData(data: any): string {
  return JSON.stringify(data, null, 1)
    .replace(/  +/g, " ")
    .replace(/\n/g, "")
}
</script>