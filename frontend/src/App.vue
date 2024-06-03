<template>
    <Slide :close-on-navigation="false">
        <!-- <Slide class="plswork" :closeOnNavigation="true"> -->
        <span class="nav">
            <span v-if="store.loggedIn && store.user" class="userinformation">
                <!-- <img v-if="store.user.photo != '' && store.user.photo.includes('https://')" :src="store.user.photo"
                    alt="" width="30" height="30" @error="imageLoadError"> -->
                <span>{{ store.user.name }}</span>
                <br />
            </span>
            <br>

            <span>
                <router-link to="/">Home</router-link>
                <br />
            </span>

            <span v-if="store.loggedIn">
                <router-link to="/profile">Profile</router-link>
                <br>
            </span>

            <span v-if="!store.loggedIn">
                <router-link to="/login">Login</router-link>
                <br>
            </span>

            <span v-if="store.loggedIn">
                <span class="fakelink" @click="logout">Logout</span>
                <br />
            </span>

            <span>
                <router-link to="/dsg">Datenschutz</router-link>
                <br />
            </span>

            <span>
                <router-link to="/about">About</router-link>
                <br />
            </span>
            <span>
                <span class="fakelink" @click="toGuessGame">Zum Tippspiel</span>
                <br />
            </span>

            <br />
            <span>
                Theme:
                <ThemeSelect />
            </span>
            <span v-if="frontendVersion != '' || backendVersion != ''"> <br /><br /> </span>
            <span v-if="frontendVersion != ''">
                <br />
                <span>Frontend: {{ frontendVersion }}</span>
            </span>
            <span v-if="backendVersion != ''">
                <br />
                <span>Backend: {{ backendVersion }}</span>
            </span>
        </span>
    </Slide>
    <div class="content1">
        <br />
        <div class="content">
            <router-view />
            <br />
            <!-- <br>
            user: {{ user }}
            <br>
            loggedin: {{ loggedin }} -->
        </div>
    </div>
</template>

<script lang="ts">
import { Slide } from "vue3-burger-menu" //   import the CSS transitions you wish to use, in this case we are using `Slide`

import { getAxiosConfig, loadUser, getLocation } from "@/func"

import { defineComponent } from "vue"
import axios from "axios"
import { store } from "@/store"
import ThemeSelect from "./components/Select/ThemeSelect.vue"

export default defineComponent({
    name: "App",
    components: {
        Slide,
        ThemeSelect
    },
    data() {
        return {
            frontendVersion: "",
            backendVersion: "",
            store
            // theme: ""
        }
    },
    computed: {
        admin() {
            if (this.store.user && this.store.user.privileges && this.store.user.privileges > 8) {
                return true
            } else {
                return false
            }
        }
    },
    watch: {
        theme() {
            console.log("theme changed")
            this.setTheme()
        }
    },
    mounted() {
        loadUser()
        // if (localStorage.theme) {
        //     this.theme = localStorage.theme
        //     this.setTheme()
        // }
        axios
            .request(getAxiosConfig("/version"))
            .then(response => {
                if (response.data.version && response.data.version != "") {
                    let backendVersion = response.data.version
                    if (backendVersion != "latest") {
                        backendVersion = "v" + backendVersion
                    }
                    this.backendVersion = backendVersion
                }
            })
            .catch(error => {
                console.log("error", error)
            })

        let frontendVersion = import.meta.env.VITE_F_VERSION

        if (frontendVersion && frontendVersion != "") {
            if (frontendVersion != "latest") {
                frontendVersion = "v" + frontendVersion
            }
            console.log("frontendVersion: '" + frontendVersion + "'")
            this.frontendVersion = frontendVersion
        }
    },
    methods: {
        imageLoadError() {
            console.log("Image failed to load")
        },
        getCookie(name) {
            const value = `; ${document.cookie}`
            const parts = value.split(`; ${name}=`)
            if (parts && parts != undefined) {
                if (parts.length === 2) {
                    const last = parts.pop()
                    if (last && last != undefined) {
                        return last.split(";").shift()
                    }
                }
            }
            return ""
        },
        redirectToLogin() {
            window.location.href = "/login" + "?from=" + getLocation()
        },
        redirectToRegister() {
            window.location.href = "/register" + "?from=" + getLocation()
        },
        logout() {
            axios.request(getAxiosConfig("/auth/logout")).then(response => {
                console.log(response)
                localStorage.removeItem("loggedin")
                // this.$router.push('/')
                window.location.reload()
            })
        },
        toGuessGame() {
            // https://tipp.fabkli.ch
            window.location.href = "https://tipp.fabkli.ch"
        },
        setTheme() {
            // document.firstElementChild?.setAttribute("data-theme", this.theme)
            // localStorage.theme = this.theme
        }
    }
})
</script>

<style land="scss">
#app {
    font-family: Avenir, Helvetica, Arial, sans-serif;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
    text-align: center;
    /* color: #63819f; */
    color: var(--font-color);
}

/* .nav { */
/* padding: 30px; */
/* width: 110%; */
/* z-index: 1000; */
/* left: -10px; */
/* overflow: hidden; */
/* position: fixed; */
/* background-color: var(--main-background-color); */
/* top: 0px; */
/* padding-top: 5px; */
/* padding-bottom: 5px; */
/* } */

/* .plswork{
    background-color: chocolate;
} */

.bm-menu {
    display: block;
    /* background-color: #ff0000 !important; */
    /* background: #ff00ff !important; */
    background-color: color-mix(in hsl, var(--main-background-color), var(--color-full) 5%) !important;
}

/* .bm-item-list { */
/* background-color: #00ff00 !important; */
/* background: #0000ff !important; */
/* } */

.bm-burger-button {
    left: 20px !important;
    top: 20px !important;
}

.bm-burger-bars {
    background-color: var(--color-half) !important;
}

/* .customselect{ */
/* all:unset; */
/* all: initial !important; */
/* } */
.content {
    /* top: ; */
    max-width: 600px;
    margin: 0 auto;
    margin-top: 40px;
    /* margin-top: 70px; */
    z-index: 100;
    padding-left: 4px;
    padding-right: 4px;
}

.content1 {
    margin-top: 0px;
    height: 99vh;
    width: 99vw;
    margin: 0 auto;
    margin-left: 0px;
}

.nav {
    /* line-break: anywhere; */
    display: block !important;
    justify-content: left;
    text-align: left;
    padding: 0px;
    /* background-color: crimson; */
}

.nav a {
    font-weight: bold;
    color: var(--inactive-link-color);
    /* color: var(--inactive-link-color); */
}

.nav a.router-link-exact-active {
    color: var(--color-full);
    /* color: var(--active-link-color); */
    color: var(--active-link-color);
}

.fakelink {
    font-weight: bold;
    color: var(--inactive-link-color);
    /* color: var(--inactive-link-color); */
    text-decoration: underline;
    cursor: pointer;
}

/* html {
    background-color: var(--main-background-color);
} */

/* .htmlclass{
    background-color: yellow;
    background-color: var(--main-background-color);
} */

/* body {
    background-color: rgb(60, 255, 0);
} */

/* .bodyclass{
    background-color: var(--main-background-color);
} */

body {
    /* max-width: 600px;
    margin: 0 auto; */
    /* margin: 0 0; */
    max-width: 100vw;
    /* max-width: 99vw; */
    max-height: 100vh;
    /* max-height: 99vh; */

    margin: 0 auto;
}

.static {
    outline: 2px solid #d62506;
    outline-offset: 0px;
}

.spaceBottom {
    margin-bottom: 10px;
}
</style>
