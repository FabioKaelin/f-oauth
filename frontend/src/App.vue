

<template>
    <layout></layout>
    <router-view />

</template>

<script lang="ts">

import { defineComponent } from 'vue'

import Layout from './components/layout.vue';
import { loadLoggedIn } from './helper/userStatus';
export default defineComponent({
    name: "App",
    components: { Layout },
    data() {
        return {
            interval: {} as  ReturnType<typeof setTimeout>

        }
    },
    methods: {
        loadLoggedIn() {
            loadLoggedIn()
        },
    },

    beforeDestroy() {
        clearInterval(this.interval)
    },
    created() {
        this.loadLoggedIn()
        this.interval = setInterval(() => this.loadLoggedIn(), 60*1000);
    },
})
// import HelloWorld from './components/HelloWorld.vue'
</script>

<style scoped>
</style>
