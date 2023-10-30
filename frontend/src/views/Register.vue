<template>
    <div class="register">
        <h2>Register</h2>
        <br>
        <input type="text" placeholder="Name" v-model="name">
        <br>
        <input type="email" placeholder="Email" v-model="email">
        <br>
        <input type="password" placeholder="Password" v-model="password">
        <br>
        <input type="password" placeholder="Password confirm" v-model="passwordConfirm">
        <br>
        <button v-if="showRegister" @click="register">Register</button>
        <br>
        <span>
            <span>if you have a account you can login <router-link :to="'/login' + getFrom()">here</router-link></span>
            <br>
            <span>you can even login without an account with google <router-link
                    :to="'/login' + getFrom()">here</router-link></span>
        </span>
    </div>
</template>

<script lang="ts">

import { defineComponent } from 'vue'
import { getAxiosConfigMethod } from '../helper/request'
import axios from 'axios'
import { getLoggedin } from '../helper/userStatus'

export default defineComponent({
    name: 'Register',
    data() {
        return {
            name: "",
            email: "",
            password: "",
            passwordConfirm: "",
            from: ""
        }
    },
    computed: {
        passwordMatch() {
            return this.password == this.passwordConfirm
        },
        showRegister() {
            return this.passwordMatch && this.name.length > 0 && this.email.length > 0 && this.password.length > 0
        }
    },
    methods: {
        getLoggedin() {
            return getLoggedin()
        },
        getFrom() {
            if (this.from == "") {
                return ""
            }
            return "?from=" + this.from
        },
        register() {
            let data = {
                name: this.name,
                email: this.email,
                password: this.password,
            }

            let data2 = JSON.stringify({
                "name": this.name,
                "email": this.email,
                "password": this.password
            })

            console.log(data)
            console.log(data2)

            axios.request(getAxiosConfigMethod("/auth/register", "POST", data)).then((res) => {
                console.log(res)
                if (res.status == 201) {
                    console.log("success")
                    this.login()
                    // this.$router.push("/login")
                }
            }).catch((err) => {
                console.error(err.response.data)
            })
        },
        login() {
            axios.request(getAxiosConfigMethod("/auth/login", "POST", {
                email: this.email,
                password: this.password
            })).then((res) => {
                console.log(res)
                if (res.status == 200) {
                    console.log("success")
                    setTimeout(() => {
                        document.location.href = this.from
                    }, 50);
                }
            }).catch((err) => {
                console.error(err.response.data)
            })
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
            console.log(window.location.origin + "?from=" + this.from)
            console.log(this.from)
            document.location.href = window.location.origin + "?from=" + this.from
        }
    }
})

</script>