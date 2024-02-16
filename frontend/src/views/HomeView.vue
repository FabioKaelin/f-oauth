<template>
    <div class="home">
        <h1>F-OAuth</h1>
        <p>Login Modul für fabkli.ch Webseiten</p>
        <br />
        <p>
            Logge die hier ein um beim <a href="https://tipp.fabkli.ch">Tippspiel</a> teilzunehmen. <br />
            Das Tippspiel ist nur zum Spass und zur Unterhaltung für die kommende Fussball EM in Deutschland gedacht. <br />
            Es werden weder Wetten, Gewinne noch Bezahlungen im Hintergrund abgewichelt und auch nicht akzeptiert.
        </p>
        <br />
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
.home {
    background-color: var(--color);
    border-radius: 10px;
    padding: 5px;
}

a {
    font-weight: bold;
    color: var(--color-half);
}

a.router-link-exact-active {
    color: var(--color-full);
}
</style>
