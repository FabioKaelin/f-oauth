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
            store.user = response1.data
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

export const getGoogleUrl = (from: string) => {
    const rootUrl = `https://accounts.google.com/o/oauth2/v2/auth`

    const options = {
        redirect_uri: import.meta.env.VITE_GOOGLE_OAUTH_REDIRECT as string,
        client_id: import.meta.env.VITE_GOOGLE_OAUTH_CLIENT_ID as string,
        access_type: "offline",
        response_type: "code",
        prompt: "consent",
        scope: ["https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"].join(" "),
        state: from
    }

    const qs = new URLSearchParams(options)

    return `${rootUrl}?${qs.toString()}`
}

export const getGithubUrl = (from: string) => {
    const rootURl = `https://github.com/login/oauth/authorize`

    console.log("import.meta.env.VITE_GITHUB_OAUTH_CLIENT_ID", import.meta.env.VITE_GITHUB_OAUTH_CLIENT_ID)
    console.log("import.meta.env.VITE_GITHUB_OAUTH_CLIENT_ID as string", import.meta.env.VITE_GITHUB_OAUTH_CLIENT_ID as string)
    console.log("import.meta.env.VITE_GITHUB_OAUTH_REDIRECT_URL", import.meta.env.VITE_GITHUB_OAUTH_REDIRECT_URL)

    const options = {
        client_id: import.meta.env.VITE_GITHUB_OAUTH_CLIENT_ID as string,
        redirect_uri: import.meta.env.VITE_GITHUB_OAUTH_REDIRECT_URL as string,
        scope: "user:email",
        state: from
    }

    const qs = new URLSearchParams(options)

    return `${rootURl}?${qs.toString()}`
}
