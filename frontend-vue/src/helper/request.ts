export function getAxiosConfigMethod(url: string, method: string, data: any) {
    const backendURL = import.meta.env.VITE_SERVER_ENDPOINT as string
    // const backendURL = process.env.VUE_APP_BACKEND_URL
    // const backendURL = "http://localhost:8000"
    // console.log(backendURL);

    const config = {
        url: "/api" + url,
        baseURL: backendURL,
        method: method,
        caches: "no-cache",
        credentials: "include",
        withCredentials: true,
        headers: {
            accept: "application/json"
        },
        // withCredentials: true,
        data: data
    }
    // console.log(config);

    return config
}



export function getAxiosConfig(url: string) {
    const backendURL = import.meta.env.VITE_SERVER_ENDPOINT as string
    // console.log(backendURL);

    // const backendURL = "http://localhost:8000"
    const config = {
        url: "/api" + url,
        baseURL: backendURL,
        method: "GET",
        caches: "no-cache",
        credentials: "include",
        withCredentials: true,
        headers: {
            accept: "application/json",
        },
        // withCredentials: true
    }
    return config
}