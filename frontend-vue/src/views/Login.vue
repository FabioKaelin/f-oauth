<template>
    <div class="login<">
        <!-- Login {{ count }} -->
        login
        <hr>
        <a :href="getGoogleUrl()">Continue with Google</a>
        <hr>
        <input type="email" placeholder="Email" v-model="email">
        <br>
        <input type="password" placeholder="Password" v-model="password">
        <br>
        <button @click="login">Login</button>
    </div>
</template>

<script lang="ts">

import { defineComponent } from 'vue'
import { getGoogleUrl } from "../helper/getGoogleUrl"
import axios from 'axios'
import { getAxiosConfigMethod } from '../helper/request'

export default defineComponent({
    name: 'Login',
    // props: {
    //     name: String,
    //     msg: { type: String, required: true }
    // },
    data() {
        return {
            count: 1,
            from: "",
            email: "",
            password: ""
        }
    },
    methods: {
        getGoogleUrl() {
            return getGoogleUrl(this.from)
        },
        login() {
            axios.request(getAxiosConfigMethod("/auth/login","POST", {
                email: this.email,
                password: this.password
            })).then((res) => {
                console.log(res)
            }).catch((err) => {
                console.error(err)
            })
        }
    },
    mounted() {
        let fromDirect = this.$route.query.from
        if (fromDirect == undefined || fromDirect == null) {
            fromDirect = "http://localhost:5173/profile"
        }
        this.from = fromDirect.toString()
        console.log(this.from)
        //     this.name // type: string | undefined
        //     this.msg // type: string
        //     this.count // type: number
    }
})

</script>

<style scoped></style>