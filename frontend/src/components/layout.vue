<template>
    <div class="linkbuttons">

        <Slide :closeOnNavigation="true" noOverlay >
            <router-link v-if="loggedIn" to="/">
                <span>Home</span>
            </router-link>
            <!-- <router-link to="/profile">
                <span>Profile</span>
            </router-link> -->
            <router-link v-if="!loggedIn" to="/login">
                <span>Login</span>
            </router-link>
            <router-link v-if="!loggedIn" to="/register">
                <span>Register</span>
            </router-link>
            <a v-if="loggedIn" @click="logout">
                <span>
                    Logout
                </span>
            </a>
        </Slide>
    </div>
</template>

<script lang="ts">
import axios from 'axios';
import { defineComponent } from 'vue'
import { getAxiosConfig } from '../helper/request';
import { Slide } from 'vue3-burger-menu'  // import the CSS transitions you wish to use, in this case we are using `Slide`
import { getLoggedin, loggedIn } from '../helper/userStatus';
// const Slide = require('vue-burger-menu').Slide
export default defineComponent({
    name: 'Layout',
    components: {
        Slide
    },
    data() {
        return {
            loggedIn: loggedIn,
            interval: {} as ReturnType<typeof setTimeout>
        }
    },
    created() {
        this.interval = setInterval(() => this.getLoggedin(), 3000);
    },
    beforeDestroy() {
        clearInterval(this.interval)
    },
    methods: {
        logout() {
            axios.request(getAxiosConfig('/auth/logout')).then((response) => {
                console.log(response);
                localStorage.removeItem("loggedin");
                // this.$router.push('/')
                window.location.reload()
            })
        },
        getLoggedin(): boolean {
            let loggedIn = getLoggedin()
            this.loggedIn = loggedIn
            return loggedIn
        }

    }
})
</script>

<style scoped>
.linkbreaker_2 {
    display: none;
}

a {
    cursor: pointer;
}

/* .linkbuttons { */
/* padding-top: 30px; */
/* padding-bottom: 5px; */
/* background-color: rgb(0, 50, 63); */
/* margin-bottom: 5px;
/* }

/* @media only screen and (max-width: 600px) {
    .linkbuttons {
        background-color: rgb(63, 0, 0);
    }

    .linkebreak {
        display: none;
    }

    .linkbreaker_2 {
        display: inline;
    }
} */
</style>