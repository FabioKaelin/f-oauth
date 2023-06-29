<template>
    <div>
        Register
        <br>
        <input type="text" placeholder="Name" v-model="name">
        <br>
        <input type="email" placeholder="Email" v-model="email">
        <br>
        <input type="password" placeholder="Password" v-model="password">
        <br>
        <input type="password" placeholder="Password confirm"  v-model="passwordConfirm">
        <br>
        <button @click="register">Register</button>
    </div>
</template>

<script lang="ts">

import { defineComponent } from 'vue'
import { getAxiosConfigMethod } from '../helper/request'
import axios from 'axios'

export default defineComponent({
    name: 'Register',
    // props: {
    //     name: String,
    //     msg: { type: String, required: true }
    // },
    data() {
        return {
            name: "",
            email: "",
            password: "",
            passwordConfirm: ""
        }
    },
    methods: {
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

            axios.request(getAxiosConfigMethod("/auth/register","POST", data)).then((res) => {
                console.log(res)
            }).catch((err) => {
                console.error(err)
            })
        }
    },
    // mounted() {
    //     this.name // type: string | undefined
    //     this.msg // type: string
    //     this.count // type: number
    // }
})

</script>