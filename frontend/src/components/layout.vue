<template>
    <div>
        <!-- router link to / -->
        <router-link to="/">Home</router-link> |
        <!-- router link to /profile -->
        <router-link to="/profile">Profile</router-link> |
        <!-- router link to /login -->
        <router-link to="/login">Login</router-link> |
        <!-- router link to /register -->
        <router-link to="/register">Register</router-link> |
        <button @click="logout">Logout</button>
    </div>
</template>

<script lang="ts">

import axios from 'axios';
import { defineComponent } from 'vue'
import { getAxiosConfig } from '../helper/request';

export default defineComponent({
    name: 'Layout',
    methods: {
        logout() {
            axios.request(getAxiosConfig('/auth/logout')).then((response) => {
                console.log(response);
                this.$router.push('/')
            })
        }
    }
})
</script>