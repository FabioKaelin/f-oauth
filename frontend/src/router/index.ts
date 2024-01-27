import { createRouter, createWebHistory, RouteRecordRaw } from "vue-router"
import HomeView from "../views/HomeView.vue"
import PageNotFoundView from "../views/404View.vue"
import LoginView from "@/views/LoginView.vue"

const routes: Array<RouteRecordRaw> = [
    {
        path: "/",
        name: "home",
        // component: DashboardView
        component: HomeView
    },
    {
        path: "/login",
        name: "login",
        component: LoginView
        // component: HomeView
    },
    {
        path: "/:pathMatch(.*)",
        name: "404",
        component: PageNotFoundView
    },
]

const router = createRouter({
    // history: createWebHistory(import.meta.env.BASE_URL),
    history: createWebHistory(),
    routes
})

export default router
