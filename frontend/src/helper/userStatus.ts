import axios from "axios"

export function setLoggedin(loggedInNew: boolean) {
    sessionStorage.setItem("loggedin", loggedInNew.toString())
    loggedIn = loggedInNew
}

export function getLoggedin() {
    // to boolean

    const loggedIn = sessionStorage.getItem("loggedin")
    if (loggedIn === null) {
        return false
    }
    return loggedIn === "true"
}

let loggedIn = false
export { loggedIn }

export function loadLoggedIn() {
    const backendURL = import.meta.env.VITE_SERVER_ENDPOINT as string
    const config = {
        url: "/api/users/me",
        baseURL: backendURL,
        method: "GET",
        caches: "no-cache",
        credentials: "include",
        withCredentials: true,
        headers: {
            accept: "application/json"
        },
        validateStatus: function (status: number) {
            return status < 500 // Resolve only if the status code is less than 500
        }
    }

    axios
        .request(config)
        .then(res => {
            if (res.status === 200) {
                console.log("logged in")
                setLoggedin(true)
                sessionStorage.setItem("loggedin", true.toString())
            } else if (res.status === 401) {
                console.log("not logged in 401")
                setLoggedin(false)
                sessionStorage.setItem("loggedin", false.toString())
            } else {
                console.log("not logged in")
                setLoggedin(false)
                sessionStorage.setItem("loggedin", false.toString())
            }
        })
        .catch(err => {
            console.log("b")
            console.log(err.response.status)
        })
}
