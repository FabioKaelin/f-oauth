import { createRouter, createWebHistory } from 'vue-router'
// import Home from '@/views/Home.vue'
import Home from './views/Home.vue'
// import Home from '@/views/Home.vue'

export default createRouter({
    history: createWebHistory(),
    routes: [
        {
            path: '/',
            name: 'Home',
            component: Home,
        },
        // {
        //     path: '/about',
        //     component: () => import('@/views/About.vue'),
        // },
        // {
        //     path: '/contact',
        //     component: () => import('@/views/Contact.vue'),
        // },
    ],
})
