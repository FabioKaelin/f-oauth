<template>
    <div class="home">
        <h1>F-OAuth</h1>
        <p>Login Modul für fabkli.ch Webseiten</p>
        <br>
        <router-link to="/profile">
            <span>Profil/Login</span>
            <br />
        </router-link>
    </div>
</template>

<script lang="ts">
import { defineComponent } from "vue"

export default defineComponent({
    name: "HomeView"
})
</script>

<style scoped lang="scss">

a {
    font-weight: bold;
    color: var(--color-half);
}

a.router-link-exact-active {
    color: var(--color-full);
}
</style>