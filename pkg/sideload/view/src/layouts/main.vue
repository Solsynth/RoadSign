<template>
  <n-layout>
    <n-layout-header class="header py-[8px] px-[36px]" bordered>
      <div class="flex items-center gap-2">
        <router-link class="link" to="/">
          RoadSign<i>!</i>
        </router-link>
      </div>

      <div class="nav-menu">
        <div class="h-full flex items-center header-nav">
          <n-menu v-model:value="key" :options="options" mode="horizontal" />
        </div>
      </div>
    </n-layout-header>
    <n-layout-content class="h-[calc(100vh-70px)] container mx-auto" content-style="padding: 24px">
      <router-view />
    </n-layout-content>
  </n-layout>
</template>

<script setup lang="ts">
import { type MenuOption, NIcon, NLayout, NLayoutContent, NLayoutHeader, NMenu } from "naive-ui"
import { type Component, h, ref } from "vue"
import { Dashboard } from "@vicons/carbon"
import { RouterLink, useRoute, useRouter } from "vue-router"

const route = useRoute()
const router = useRouter()
const key = ref(route.name?.toString())

router.afterEach((to) => {
  key.value = to.name?.toString() ?? "index"
})

const options: MenuOption[] = [
  {
    label: () => h(RouterLink, { to: { name: "dashboard" } }, "Dashboard"),
    icon: renderIcon(Dashboard),
    key: "dashboard"
  }
]

function renderIcon(icon: Component) {
  return () => h(NIcon, null, { default: () => h(icon) })
}
</script>

<style scoped>
.header {
  display: grid;
  grid-template-columns: 1fr auto 1fr;
  gap: 40px;
}

.link {
  all: unset;
  cursor: pointer;
}
</style>
