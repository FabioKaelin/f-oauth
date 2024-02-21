<template>
    <div class="Page404">
        <h1>404</h1>
        <p>Page not found</p>
    </div>
</template>

<script lang="ts">
import { defineComponent } from "vue"

export default defineComponent({
    name: "PageNotFoundView"
})
</script>

<style scoped lang="scss">
.Page404 {
    background-color: var(--color);
    border-radius: 10px;
    padding: 5px;
}
</style>
