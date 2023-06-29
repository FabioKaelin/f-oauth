<template>
    <div class="login<">
        <!-- Login {{ count }} -->
        login
        <hr>
        <a :href="getGoogleUrl()">Continue with Google</a>

    </div>
</template>

<script lang="ts">

import { defineComponent } from 'vue'
import {getGoogleUrl} from "../helper/getGoogleUrl"

export default defineComponent({
    name: 'Login',
    // props: {
    //     name: String,
    //     msg: { type: String, required: true }
    // },
    data() {
        return {
            count: 1,
            from: ""
        }
    },
    methods: {
        getGoogleUrl() {
            return getGoogleUrl(this.from)
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

<style scoped>



</style>