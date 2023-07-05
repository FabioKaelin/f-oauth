

<template>
    <!-- <div>
import Layout from './components/layout.vue';
        <a href="https://vitejs.dev" target="_blank">
            <img src="/vite.svg" class="logo" alt="Vite logo" />
        </a>
        <a href="https://vuejs.org/" target="_blank">
            <img src="./assets/vue.svg" class="logo vue" alt="Vue logo" />
        </a>
    </div>
    <HelloWorld msg="Vite + Vue" /> -->

    <layout></layout>
    <hr>

    <router-view />

</template>

<script lang="ts">

import { defineComponent } from 'vue'

import Layout from './components/layout.vue';
export default defineComponent({
    name: "App",
    components: { Layout }
})
// import HelloWorld from './components/HelloWorld.vue'
</script>

<style scoped>
.logo {
    height: 6em;
    padding: 1.5em;
    will-change: filter;
    transition: filter 300ms;
}

.logo:hover {
    filter: drop-shadow(0 0 2em #646cffaa);
}

.logo.vue:hover {
    filter: drop-shadow(0 0 2em #42b883aa);
}
</style>
