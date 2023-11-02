import { createRouter, createWebHistory } from "vue-router"
// import Home from './views/Home.vue'

export default createRouter({
    history: createWebHistory(),
    routes: [
        {
            path: "/",
            name: "Home",
            component: () => import("./views/ProfileView.vue")
        },
        // {
        //     path: '/',
        //     name: 'Home',
        //     component: Home,
        // },
        {
            path: "/profile",
            name: "Profile",
            component: () => import("./views/ProfileView.vue")
        },
        {
            path: "/login",
            name: "login",
            component: () => import("./views/LoginView.vue")
            // query: {
            //     from: '',
            // }
        },
        {
            path: "/register",
            name: "register",
            component: () => import("./views/RegisterView.vue")
        }
    ]
})
