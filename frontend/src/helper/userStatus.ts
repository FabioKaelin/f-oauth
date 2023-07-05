import axios from "axios";

export function setLoggedin(loggedInNew: boolean) {
    sessionStorage.setItem("loggedin", loggedInNew.toString());
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
            accept: "application/json",
        },
    }


    axios.request(config).then((res) => {
        console.log(res);
        if (res.status === 200) {
            console.log("logged in")
            setLoggedin(true)
            sessionStorage.setItem("loggedin", true.toString());
        } else {
            console.log("not logged in")
            setLoggedin(false)
            sessionStorage.setItem("loggedin", false.toString());
        }
    }).catch(() => {
        console.log("not logged in")
        setLoggedin(false)
        sessionStorage.setItem("loggedin", false.toString());
    })
}