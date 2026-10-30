<template>
    <aside v-if="requestPause.until" class="request-pause" aria-label="Anfragestatus">
        <p role="status" aria-live="polite">
            {{ requestPause.remaining > 0 ? (requestPause.kind === 'limit' ? 'Zu viele Anfragen in kurzer Zeit.' : 'Der Dienst ist vorübergehend nicht verfügbar.') : 'Die Wartezeit ist vorbei.' }}
        </p>
        <p v-if="requestPause.remaining > 0">Bitte in <strong>{{ requestPause.remaining }} Sekunden</strong> erneut versuchen.</p>
        <p v-else>Sie können Ihre Aktion jetzt erneut versuchen.</p>
        <p>Ihre Eingaben bleiben erhalten. Aktionen werden nicht automatisch wiederholt.</p>
        <button type="button" @click="requestPause.until = 0">Schliessen</button>
    </aside>
</template>
<script setup lang="ts">
import { requestPause } from "../rateLimit"
</script>
<style scoped>
.request-pause{position:fixed;bottom:1rem;left:50%;transform:translateX(-50%);z-index:10000;box-sizing:border-box;width:min(36rem,calc(100% - 2rem));max-height:70vh;overflow:auto;padding:1rem;background:#fff4ce;color:#342a00;border:2px solid #9a7000;border-radius:12px;box-shadow:0 4px 24px #0005;text-align:left;overflow-wrap:anywhere}
.request-pause p{margin:.35rem 0}.request-pause button{min-height:44px;margin-top:.5rem;padding:.4rem 1rem;color:#342a00;background:#fff;border:1px solid #9a7000;border-radius:6px}
</style>
