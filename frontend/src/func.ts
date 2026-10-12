import axios from "axios"
import { User } from "@/structs"

export function getAxiosConfigMethod(url: string, method: string, data: any) {
    const backendURL = import.meta.env.VITE_SERVER_ENDPOINT
    const config = {
        url: "/api" + url,
        baseURL: backendURL,
        method: method,
        caches: "no-cache",
        withCredentials: true,
        headers: {
            accept: "application/json",
            token: localStorage.getItem("token")
        },
        data: data
    }

    return config
}

export function getAxiosConfig(url: string) {
    const backendURL = import.meta.env.VITE_SERVER_ENDPOINT
    const config = {
        url: "/api" + url,
        baseURL: backendURL,
        method: "GET",
        caches: "no-cache",
        withCredentials: true,
        headers: {
            accept: "application/json",
            token: localStorage.getItem("token")
        }
    }
    return config
}

export function getLocation() {
    if (window.location.href != undefined) {
        return window.location.href
    } else {
        return ""
    }
}

import { store } from "./store"

function loadUser(): Promise<any> {
    return axios
        .request(getAxiosConfig("/users/me"))
        .then(response1 => {
            store.user = {
                ...response1.data,
                privileges: response1.data.privileges ?? (response1.data.role === "admin" || response1.data.role === "test-admin" ? 13 : 1)
            }
            store.loggedIn = true
            store.userLoaded = true
        })
        .catch(error => {
            if (error.response.status == 401) {
                console.debug("not logged in, (func.ts loadUser)")
                store.user = {} as User
                store.loggedIn = false
                store.userLoaded = true
                //     window.location.href = JSON.parse(error.response.data).redirect + "?from=" + getLocation()
            }
        })
}

export { loadUser }

// F-OAuth creates and validates a browser-bound OAuth state before provider login.
const getProviderUrl = (provider: string, from: string) => {
    const backend = import.meta.env.VITE_SERVER_ENDPOINT || "https://oauth.fabkli.ch"
    const qs = new URLSearchParams({ from })
    return `${backend}/api/sessions/oauth/start/${provider}?${qs.toString()}`
}
export const getGoogleUrl = (from: string) => getProviderUrl("google", from)
export const getGithubUrl = (from: string) => getProviderUrl("github", from)
