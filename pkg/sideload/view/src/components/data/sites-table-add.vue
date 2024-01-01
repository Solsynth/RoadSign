<template>
  <div>
    <n-button circle size="small" type="primary" @click="creating = true">
      <template #icon>
        <n-icon :component="Add" />
      </template>
    </n-button>

    <n-modal
      v-model:show="creating"
      class="w-[720px]"
      content-style="padding: 0"
      preset="card"
      title="Create Site"
      segmented
      closable
    >
      <div class="py-4 px-5 border border-solid border-b border-[#eee]">
        <n-input
          v-model:value="data.id"
          placeholder="Will be the file name of this file"
        />
      </div>

      <div class="relative mt-[4px] h-[540px]">
        <vue-monaco-editor
          v-model:value="data.content"
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
                @click="submit()"
              >
                <template #icon>
                  <n-icon :component="Checkmark" />
                </template>
              </n-button>
            </template>
            This operation will publish this site right away.
          </n-tooltip>
        </div>
      </div>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { NButton, NIcon, NInput, NModal, NTooltip, useMessage } from "naive-ui"
import { Add, Checkmark } from "@vicons/carbon"
import { VueMonacoEditor } from "@guolao/vue-monaco-editor"
import { ref } from "vue"
import * as yaml from "js-yaml"

const message = useMessage()

const emits = defineEmits(["reload"])

const submitting = ref(false)
const creating = ref(false)

const data = ref<any>({})

async function submit() {
  let content
  try {
    content = yaml.load(data.value.content)
  } catch (e: any) {
    message.warning(`Your configuration has some issue: ${e.message}`)
    return
  }

  submitting.value = true
  const resp = await fetch(`/webhooks/sync/${data.value.id}`, {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(content)
  })
  if (resp.status != 200) {
    message.error(`Something went wrong... ${await resp.text()}`)
  } else {
    reset()
    emits("reload")
    message.success("Your site has been created! 🎉")
    creating.value = false
  }
  submitting.value = false
}

function reset() {
  data.value.id = ""
  data.value.content = ""
}
</script>

<style scoped>
.fab {
  position: absolute;
  bottom: 16px;
  right: 24px;
}
</style>