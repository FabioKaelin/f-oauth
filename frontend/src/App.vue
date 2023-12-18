<template>
    <web-layout></web-layout>
    <router-view />
</template>

<script lang="ts">
import { defineComponent } from "vue"

import webLayout from "./components/webLayout.vue"
import { loadLoggedIn } from "./helper/userStatus"
export default defineComponent({
    name: "App",
    components: { webLayout: webLayout },
    data() {
        return {
            interval: {} as ReturnType<typeof setTimeout>
        }
    },

    beforeUnmount() {
        clearInterval(this.interval)
    },
    created() {
        this.loadLoggedIn()
        this.interval = setInterval(() => this.loadLoggedIn(), 60 * 1000)
    },
    methods: {
        loadLoggedIn() {
            loadLoggedIn()
        }
    }
})
// import HelloWorld from './components/HelloWorld.vue'
</script>

<style scoped></style>
