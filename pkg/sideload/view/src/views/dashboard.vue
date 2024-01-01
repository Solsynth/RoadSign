<template>
  <div class="flex flex-col gap-2">
    <div class="grid gap-2 grid-cols-2 lg:grid-cols-4">
      <n-card embedded>
        <n-statistic label="Status">{{ data?.status ? "Operational" : "Incident" }}</n-statistic>
      </n-card>
      <n-card embedded>
        <n-statistic label="Sites">{{ data?.sites }}</n-statistic>
      </n-card>
      <n-card embedded>
        <n-statistic label="Upstreams">{{ data?.upstreams }}</n-statistic>
      </n-card>
      <n-card embedded>
        <n-statistic label="Processes">{{ data?.processes }}</n-statistic>
      </n-card>
    </div>

    <sites-table />
  </div>
</template>

<script setup lang="ts">
import { NCard, NStatistic } from "naive-ui"
import { ref } from "vue"
import SitesTable from "@/components/data/sites-table.vue"

const data = ref<any>({})

async function readStatistics() {
  const resp = await fetch("/cgi/statistics")
  data.value = await resp.json()
}

readStatistics()
</script>