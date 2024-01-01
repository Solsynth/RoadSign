<template>
  <div>
    <n-card title="Sites">
      <template #header-extra>
        <sites-table-add @reload="readSites()" />
      </template>

      <n-data-table
        :columns="columns"
        :data="data"
        :row-key="(row: any) => row.id"
      />
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { NCard, NDataTable, NTag } from "naive-ui"
import { h, ref } from "vue"
import SitesTableExpand from "@/components/data/sites-table-expand.vue"
import SitesTableAction from "@/components/data/sites-table-action.vue"
import SitesTableAdd from "@/components/data/sites-table-add.vue"

const columns = [
  {
    type: "expand",
    renderExpand(row: any) {
      return h(SitesTableExpand, { ...row, class: "pl-[38px]" })
    }
  },
  {
    title: "ID",
    key: "id",
    render(row: any) {
      return h(NTag, { type: "info", bordered: false, size: "small" }, row?.id)
    }
  },
  {
    title: "Rules",
    key: "rules",
    render(row: any) {
      return row?.rules?.length ?? 0
    }
  },
  {
    title: "Upstreams",
    key: "upstreams",
    render(row: any) {
      return row?.upstreams?.length ?? 0
    }
  },
  {
    title: "Processes",
    key: "processes",
    render(row: any) {
      return row?.processes?.length ?? 0
    }
  },
  {
    title: "Actions",
    key: "actions",
    render(row: any) {
      return h(SitesTableAction, { ...row, onReload: () => readSites() })
    }
  }
]

const data = ref<any[]>([])

async function readSites() {
  const resp = await fetch("/cgi/sites")
  data.value = await resp.json()
}

readSites()
</script>