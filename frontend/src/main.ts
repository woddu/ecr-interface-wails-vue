import "./assets/main.css"
import { createApp } from "vue"
import { createPinia } from 'pinia'
import { createRouter, createWebHashHistory } from "vue-router"
import uiPlugin from "@nuxt/ui/vue-plugin"

import App from "./App.vue"

const app = createApp(App)
const pinia = createPinia()
const router = createRouter({
  routes: [
    {
      path: "/",
      component: () => import("./views/FileView.vue"),
      name: "File"
    },
    {
      path: "/scores",
      component: () => import("./views/ScoresView.vue"),
      name: "Scores"
    },
    {
      path: "/students",
      component: () => import("./views/StudentsView.vue"),
      name: "Students",
      children: [
        {
          path: "",
          component: () => import("./views/students/Students.vue"),
          name: "StudentsList"
        },
        {
          path: "scores",
          component: () => import("./views/students/StudentScores.vue"),
          name: "StudentScores"
        }
      ]
    }
  ],
  history: createWebHashHistory(),
})

app.use(uiPlugin)
app.use(router)
app.use(pinia)

app.mount("#app")
