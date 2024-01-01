import { createRouter, createWebHistory } from "vue-router"

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      name: "layouts.main",
      component: () => import("@/layouts/main.vue"),
      children: [
        {
          path: "/",
          name: "dashboard",
          component: () => import("@/views/dashboard.vue")
        },
      ]
    },
  ]
})

export default router
