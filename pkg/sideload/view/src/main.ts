import "./assets/main.css"

import "virtual:uno.css"

import { createApp } from "vue"
import { createPinia } from "pinia"

import root from "./root.vue"
import router from "./router"

const app = createApp(root)

app.use(createPinia())
app.use(router)

app.mount("#app")
