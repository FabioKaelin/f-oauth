<template>
    <div class="login<">
        <h1>Login</h1>
        <button type="button" class="login-with-google-button" @click="goToGoogle">Sign in with Google</button>
        <br />
        <span>you can even login without an account with google</span>
        <br />
        <span><b>or</b></span>
        <br />
        <input v-model="email" type="email" placeholder="Email" />
        <br />
        <input v-model="password" type="password" placeholder="Password" />
        <br />
        <span v-if="error != ''" class="error">{{ error }}</span>
        <br />
        <button v-if="password != '' && email != ''" @click="login">Login</button>
        <br />
        <span>if you don't have an account you can create one <router-link :to="'/register' + getFrom()">here</router-link></span>
    </div>
</template>

<script lang="ts">
import { defineComponent } from "vue"
import { getGoogleUrl } from "../helper/getGoogleUrl"
import axios from "axios"
import { getAxiosConfigMethod } from "../helper/request"
import { getLoggedin } from "../helper/userStatus"

export default defineComponent({
    name: "LoginView",
    // props: {
    //     name: String,
    //     msg: { type: String, required: true }
    // },
    data() {
        return {
            count: 1,
            from: "",
            email: "",
            password: "",
            error: ""
        }
    },
    mounted() {
        let fromDirect = this.$route.query.from
        if (fromDirect == undefined || fromDirect == null) {
            fromDirect = window.location.origin
            // fromDirect = "http://localhost:5173/profile" // TODO: lh
        }
        this.from = fromDirect.toString()
        console.log(this.from)
        if (this.getLoggedin()) {
            document.location.href = window.location.origin + "?from=" + this.from
        }
    },
    methods: {
        getGoogleUrl() {
            return getGoogleUrl(this.from)
        },
        getLoggedin() {
            return getLoggedin()
        },
        goToGoogle() {
            window.location.href = this.getGoogleUrl()
        },
        getFrom() {
            if (this.from == "") {
                return ""
            }
            return "?from=" + this.from
        },
        login() {
            axios
                .request(
                    getAxiosConfigMethod("/auth/login", "POST", {
                        email: this.email,
                        password: this.password
                    })
                )
                .then(res => {
                    console.log(res.data)
                    if (res.status == 200) {
                        console.log("success")
                        document.location.href = this.from
                    } else {
                        console.log("failed")
                        console.log(res.data)
                    }
                })
                .catch(err => {
                    if (err.response.status == 401) {
                        console.log("failed")
                        console.log(err.response.data)
                        this.error = err.response.data.message
                    } else {
                        console.log("catch")
                        console.log(err)
                    }
                })
        }
    }
})
</script>

<style scoped>
.login-with-google-button {
    cursor: pointer;

    padding: 12px 16px 12px 42px;

    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Oxygen, Ubuntu, Cantarell, "Fira Sans", "Droid Sans", "Helvetica Neue", sans-serif;

    background-image: url(data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMTgiIGhlaWdodD0iMTgiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+PGcgZmlsbD0ibm9uZSIgZmlsbC1ydWxlPSJldmVub2RkIj48cGF0aCBkPSJNMTcuNiA5LjJsLS4xLTEuOEg5djMuNGg0LjhDMTMuNiAxMiAxMyAxMyAxMiAxMy42djIuMmgzYTguOCA4LjggMCAwIDAgMi42LTYuNnoiIGZpbGw9IiM0Mjg1RjQiIGZpbGwtcnVsZT0ibm9uemVybyIvPjxwYXRoIGQ9Ik05IDE4YzIuNCAwIDQuNS0uOCA2LTIuMmwtMy0yLjJhNS40IDUuNCAwIDAgMS04LTIuOUgxVjEzYTkgOSAwIDAgMCA4IDV6IiBmaWxsPSIjMzRBODUzIiBmaWxsLXJ1bGU9Im5vbnplcm8iLz48cGF0aCBkPSJNNCAxMC43YTUuNCA1LjQgMCAwIDEgMC0zLjRWNUgxYTkgOSAwIDAgMCAwIDhsMy0yLjN6IiBmaWxsPSIjRkJCQzA1IiBmaWxsLXJ1bGU9Im5vbnplcm8iLz48cGF0aCBkPSJNOSAzLjZjMS4zIDAgMi41LjQgMy40IDEuM0wxNSAyLjNBOSA5IDAgMCAwIDEgNWwzIDIuNGE1LjQgNS40IDAgMCAxIDUtMy43eiIgZmlsbD0iI0VBNDMzNSIgZmlsbC1ydWxlPSJub256ZXJvIi8+PHBhdGggZD0iTTAgMGgxOHYxOEgweiIvPjwvZz48L3N2Zz4=);
    background-repeat: no-repeat;
    background-position: 12px 11px;
}

.error {
    color: red;
}
</style>
