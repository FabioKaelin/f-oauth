<template>
    <div>
        Profile {{ count }}
    </div>
</template>

<script lang="ts">

import { defineComponent } from 'vue'

export default defineComponent({
    name: 'Profile',
    // props: {
    //     name: String,
    //     msg: { type: String, required: true }
    // },
    data() {
        return {
            count: 1
        }
    },
    // mounted() {
    //     this.name // type: string | undefined
    //     this.msg // type: string
    //     this.count // type: number
    // }
})

</script>