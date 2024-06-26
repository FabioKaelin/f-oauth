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
        path: "/profile",
        name: "profile",
        component: () => import("../views/ProfileView.vue")
    },
    {
        path: "/login",
        name: "login",
        component: LoginView
        // component: HomeView
    },
    {
        path: "/dsg",
        name: "dsg",
        component: () => import("../views/DSGView.vue")
    },
    {
        path: "/about",
        name: "about",
        component: () => import("../views/AboutView.vue")
    },
    {
        path: "/:pathMatch(.*)",
        name: "404",
        component: PageNotFoundView
    }
]

const router = createRouter({
    // history: createWebHistory(import.meta.env.BASE_URL),
    history: createWebHistory(),
    routes
})

export default router
