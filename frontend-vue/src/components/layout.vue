<template>
    <div>
        <!-- router link to / -->
        <router-link to="/">Home</router-link> |
        <!-- router link to /profile -->
        <router-link to="/profile">Profile</router-link> |
        <!-- router link to /login -->
        <router-link to="/login">Login</router-link> |
        <!-- router link to /register -->
        <router-link to="/register">Register</router-link>
    </div>
</template>

<script lang="ts">

import { defineComponent } from 'vue'

export default defineComponent({
    name: 'Layout'
})
</script>