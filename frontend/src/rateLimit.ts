import axios from "axios"
import { reactive } from "vue"

export const requestPause = reactive({ until: 0, kind: "", remaining: 0 })
let timer: ReturnType<typeof setInterval> | undefined

export function retrySeconds(value: unknown, fallback = 30): number {
    const numeric = Number(value)
    const seconds = Number.isFinite(numeric) && numeric > 0 ? numeric : (typeof value === "string" ? (Date.parse(value) - Date.now()) / 1000 : fallback)
    return Math.max(1, Math.min(86400, Math.ceil(Number.isFinite(seconds) ? seconds : fallback)))
}

export function showRequestPause(error: any): string | undefined {
    const response = error?.response
    if (response?.status !== 429 && response?.status !== 503) return undefined
    const seconds = retrySeconds(response.headers?.["retry-after"] ?? response.data?.retry_after)
    const until = Date.now() + seconds * 1000
    if (until >= requestPause.until) {
        requestPause.until = until
        requestPause.kind = response.status === 429 ? "limit" : "unavailable"
    }
    requestPause.remaining = Math.ceil((requestPause.until - Date.now()) / 1000)
    if (timer === undefined) timer = setInterval(() => {
        requestPause.remaining = Math.max(0, Math.ceil((requestPause.until - Date.now()) / 1000))
        if (requestPause.remaining === 0 && timer !== undefined) { clearInterval(timer); timer = undefined }
    }, 1000)
    return response.status === 429 ? `Zu viele Anfragen. Bitte in ${seconds} Sekunden erneut versuchen.` : `Der Dienst ist vorübergehend nicht verfügbar. Bitte in ${seconds} Sekunden erneut versuchen.`
}

axios.interceptors.response.use(response => response, error => {
    const message = showRequestPause(error)
    if (message && error.response?.data && typeof error.response.data === "object") error.response.data.message = message
    return Promise.reject(error)
})
