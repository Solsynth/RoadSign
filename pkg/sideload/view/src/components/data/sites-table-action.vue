<template>
  <div class="flex gap-[4px]">
    <n-button size="small" @click="publishing = true">
      <template #icon>
        <n-icon :component="CloudUpload" />
      </template>
    </n-button>
    <n-button size="small" @click="editConfig()">
      <template #icon>
        <n-icon :component="Edit" />
      </template>
    </n-button>

    <n-modal
      v-model:show="publishing"
      class="w-[720px]"
      preset="card"
      title="Publish Artifacts"
      segmented
      closable
    >
      We are sorry about this tool isn't completed yet. <br>
      For now, you can use our <b>Wonderful Command Line Tool —— RDS</b> <br>
      Learn more on our <a href="https://wiki.smartsheep.studio/roadsign/index.html" target="_blank">official wiki</a>.
      <br>
      <br>
      Install it by this command below
      <n-code code="go install code.smartsheep.studio/goatworks/roadsign/pkg/cmd/rds@latest" />
      <br>
      Then connect your rds client to this server
      <n-code :code="`rds connect <name> ${host} <credentials>`" />
      <br>
      After that you can publish your stuff (You need to compress them to zip archive before publish)
      <n-code :code="`rds deploy <name> ${props.id} <upstream id or process id>`" />
    </n-modal>

    <n-modal
      v-model:show="editing"
      class="w-[720px]"
      content-style="padding: 0"
      preset="card"
      title="Edit Configuration"
      segmented
      closable
    >
      <div class="relative h-[540px]">
        <vue-monaco-editor
          v-model:value="config"
          :options="{ automaticLayout: true, minimap: { enabled: false } }"
          language="yaml"
        />

        <div class="fab">
          <n-tooltip placement="left">
            <template #trigger>
              <n-button
                circle
                type="primary"
                size="large"
                class="shadow-lg"
                :loading="submitting"
                @click="syncConfig()"
              >
                <template #icon>
                  <n-icon :component="Save" />
                </template>
              </n-button>
            </template>
            This operation will restart all processes related. Service may interrupted for some while.
          </n-tooltip>
        </div>
      </div>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { NButton, NCode, NIcon, NModal, NTooltip, useMessage } from "naive-ui"
import { CloudUpload, Edit, Save } from "@vicons/carbon"
import { ref } from "vue"
import { VueMonacoEditor } from "@guolao/vue-monaco-editor"
import * as yaml from "js-yaml"

const message = useMessage()

const props = defineProps<{ id: string, rules: any[], upstreams: any[], processes: any[] }>()
const emits = defineEmits(["reload"])
const host = location.protocol + "//" + location.host

const submitting = ref(false)

const publishing = ref(false)
const editing = ref(false)

const config = ref<string | undefined>(undefined)

async function editConfig() {
  const resp = await fetch(`/cgi/sites/cfg/${props.id}`)
  config.value = await resp.text()
  editing.value = true
}

async function syncConfig() {
  if (config.value == null) return

  let content
  try {
    content = yaml.load(config.value)
  } catch (e: any) {
    message.warning(`Your configuration has some issue: ${e.message}`)
    return
  }

  submitting.value = true
  const resp = await fetch(`/webhooks/sync/${props.id}`, {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(content)
  })
  if (resp.status != 200) {
    message.error(`Something went wrong... ${await resp.text()}`)
  } else {
    emits("reload")
  }
  submitting.value = false
}
</script>

<style scoped>
.fab {
  position: absolute;
  bottom: 16px;
  right: 24px;
}
</style>