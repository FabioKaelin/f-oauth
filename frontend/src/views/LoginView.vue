<template>
    <div class="login">
        <h1>Login</h1>
        <span v-if="errorHeader != ''" class="error">
            {{ errorHeader }}
            <br>
        </span>

        <button type="button" class="login-with-google-button" @click="goToGoogle">Sign in with Google</button>
        <br />
        <button type="button" class="login-with-github-button" @click="goToGitHub">Sign in with GitHub</button>
        <br />
        <br />
        <span class="text">oder mit Email und Passwort</span>
        <br />
        <br />
        <input v-model="email" class="textInput" type="email" placeholder="Email" />
        <br v-if="!emailValid && email.length != 0" />
        <span v-if="!emailValid && email.length != 0" style="color: red">Die Email ist nicht gültig</span>
        <br />
        <input v-model="password" class="textInput" type="password" placeholder="Password" />
        <br />
        <span v-if="error != ''" class="error">{{ error }}</span>
        <br />
        <button v-if="password != '' && email != ''" @click="login">Login</button>
        <br />
        <hr />
        <br />
        <span class="text">oder wenn du noch kein Account hast und dich nicht mit Google oder GitHub einloggen möchtest
            (was ohne Account funktioniert) kannst du dich hier registrieren</span>
        <br />
        <input v-model="rname" class="textInput" type="text" placeholder="Name" />
        <br v-if="!usernameValid && rname.length != 0" />
        <span v-if="!usernameValid && rname.length != 0" style="color: red">Der Name darf nur Buchstaben, Zahlen,
            Leerzeichen, - und _ enthalten und mindestens 3 und maximal 20 Zeichen lang sein</span>
        <br />
        <input v-model="remail" class="textInput" type="email" placeholder="Email" />
        <br v-if="!emailrValid && remail.length != 0" />
        <span v-if="!emailrValid && remail.length != 0" style="color: red">Die Email ist nicht gültig</span>
        <br />
        <input v-model="rpassword" class="textInput" type="password" placeholder="Password" />
        <br v-if="!passwordrValid && rpassword.length != 0" />
        <span v-if="!passwordrValid && rpassword.length != 0" style="color: red">Das Passwort muss mindestens 8 Zeichen
            lang sein und mindestens eine Großbuchstabe, eine Kleinbuchstabe, eine Zahl und ein Sonderzeichen
            enthalten.</span>
        <br />
        <input v-model="rpasswordConfirm" class="textInput" type="password" placeholder="Password confirm" />
        <br v-if="!passwordMatch && rpasswordConfirm.length != 0" />
        <span v-if="!passwordMatch && rpasswordConfirm.length != 0" style="color: red">Die Passwörter stimmen nicht
            überein</span>
        <br />
        <button v-if="showRegister && usernameValid && emailrValid && passwordrValid" class="textInput"
            @click="register">Register</button>
    </div>
</template>

<script lang="ts">
import { defineComponent } from "vue"
import { getGithubUrl, getGoogleUrl } from "../func"
import axios from "axios"
import { getAxiosConfigMethod } from "../func"
import { store } from "../store"

export default defineComponent({
    name: "LoginView",
    // props: {
    //     name: String,
    //     msg: { type: String, required: true }
    // },
    data() {
        return {
            count: 1,
            from: "",
            email: "",
            password: "",
            error: "",
            errorHeader: "",
            rname: "",
            remail: "",
            rpassword: "",
            rpasswordConfirm: ""
        }
    },
    computed: {
        passwordMatch() {
            return this.rpassword == this.rpasswordConfirm
        },
        showRegister() {
            return this.passwordMatch && this.rname.length > 0 && this.remail.length > 0 && this.rpassword.length > 0
        },
        usernameValid() {
            // validate guessgroup name Max 20 Char Only a-z A-Z 0-9 "-" "_" öäüÖÄÜêÊéàèÉÀÈç and " "

            const regex = /^[a-zA-Z0-9_äöüÄÖÜêÊéàèÉÀÈç -]{3,20}$/;
            if (this.rname.match(regex)) {
                return true
            } else {
                return false
            }
        },
        emailrValid() {
            // validate email
            const emailRegex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
            if (this.remail.match(emailRegex)) {
                return true
            } else {
                return false
            }
        },
        emailValid() {
            // validate email
            const emailRegex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
            console.log(".............")
            console.log(this.email)
            if (this.email.match(emailRegex)) {
                console.log("true")
                return true
            } else {
                console.log("false")
                return false
            }
        },
        passwordValid() {
            // validate password
            const passwordRegex = /^(?=.*?[A-Z])(?=(.*[a-z]){1,})(?=(.*[\d]){1,})(?=(.*[\W]){1,})(?!.*\s).{8,}$/;
            if (this.password.match(passwordRegex)) {
                return true
            } else {
                return false
            }
        },
        passwordrValid() {
            // validate password
            const passwordrRegex = /^(?=.*?[A-Z])(?=(.*[a-z]){1,})(?=(.*[\d]){1,})(?=(.*[\W]){1,})(?!.*\s).{8,}$/;
            if (this.rpassword.match(passwordrRegex)) {
                return true
            } else {
                return false
            }
        }
    },
    mounted() {
        let fromDirect = this.$route.query.from
        if (fromDirect == undefined || fromDirect == null) {
            fromDirect = window.location.origin
        }
        this.from = fromDirect.toString()
        console.log(this.from)
        if (store.loggedIn) {
            document.location.href = window.location.origin + "?from=" + this.from
        }

        let errorHeader = this.$route.query.error
        if (errorHeader != undefined) {
            let errorHeaderStr = errorHeader.toString()
            if (errorHeaderStr == "already_signed_up_with_different_method") {
                this.errorHeader = "Du hast dich bereits mit einer anderen Methode angemeldet. Bitte melde dich mit der gleichen Methode an."
            } else if (errorHeaderStr == "error_occured") {
                this.errorHeader = "Ein Fehler ist aufgetreten. Bitte versuche es erneut."
            }
            this.$router.replace({ query: {} })
        }
    },
    methods: {
        getGoogleUrl() {
            return getGoogleUrl(this.from)
        },
        goToGoogle() {
            window.location.href = this.getGoogleUrl()
        },
        goToGitHub() {
            window.location.href = getGithubUrl(this.from)
        },
        getFrom() {
            if (this.from == "") {
                return ""
            }
            return "?from=" + this.from
        },
        register() {
            if (this.rpassword != this.rpasswordConfirm) {
                return
            }
            if (!this.usernameValid) {
                return
            }
            if (!this.emailrValid) {
                return
            }
            if (!this.passwordrValid) {
                return
            }
            let data = {
                name: this.rname,
                email: this.remail,
                password: this.rpassword
            }

            let data2 = JSON.stringify({
                name: this.rname,
                email: this.remail,
                password: this.rpassword
            })

            console.log(data)
            console.log(data2)

            axios
                .request(getAxiosConfigMethod("/auth/register", "POST", data))
                .then(res => {
                    console.log(res)
                    if (res.status == 201) {
                        console.log("success")
                        this.email = this.remail
                        this.password = this.rpassword
                        this.login()
                        // this.$router.push("/login")
                    }
                })
                .catch(err => {
                    console.error(err.response.data)
                })
        },
        login() {
            axios
                .request(
                    getAxiosConfigMethod("/auth/login", "POST", {
                        email: this.email,
                        password: this.password
                    })
                )
                .then(res => {
                    console.log(res.data)
                    if (res.status == 200) {
                        console.log("success")
                        document.location.href = this.from
                    } else {
                        console.log("failed")
                        console.log(res.data)
                    }
                })
                .catch(err => {
                    if (err.response.status == 401) {
                        console.log("failed")
                        console.log(err.response.data)
                        this.error = err.response.data.message
                    } else {
                        this.error = err.response.data.message
                        console.log("catch")
                        console.log(err)
                    }
                })
        }
    }
})
</script>

<style scoped>
.login {
    background-color: var(--color);
    border-radius: 10px;
    padding: 5px;
}

.login-with-google-button {
    cursor: pointer;

    padding: 12px 16px 12px 42px;

    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Oxygen, Ubuntu, Cantarell, "Fira Sans", "Droid Sans", "Helvetica Neue", sans-serif;

    background-image: url(data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMTgiIGhlaWdodD0iMTgiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+PGcgZmlsbD0ibm9uZSIgZmlsbC1ydWxlPSJldmVub2RkIj48cGF0aCBkPSJNMTcuNiA5LjJsLS4xLTEuOEg5djMuNGg0LjhDMTMuNiAxMiAxMyAxMyAxMiAxMy42djIuMmgzYTguOCA4LjggMCAwIDAgMi42LTYuNnoiIGZpbGw9IiM0Mjg1RjQiIGZpbGwtcnVsZT0ibm9uemVybyIvPjxwYXRoIGQ9Ik05IDE4YzIuNCAwIDQuNS0uOCA2LTIuMmwtMy0yLjJhNS40IDUuNCAwIDAgMS04LTIuOUgxVjEzYTkgOSAwIDAgMCA4IDV6IiBmaWxsPSIjMzRBODUzIiBmaWxsLXJ1bGU9Im5vbnplcm8iLz48cGF0aCBkPSJNNCAxMC43YTUuNCA1LjQgMCAwIDEgMC0zLjRWNUgxYTkgOSAwIDAgMCAwIDhsMy0yLjN6IiBmaWxsPSIjRkJCQzA1IiBmaWxsLXJ1bGU9Im5vbnplcm8iLz48cGF0aCBkPSJNOSAzLjZjMS4zIDAgMi41LjQgMy40IDEuM0wxNSAyLjNBOSA5IDAgMCAwIDEgNWwzIDIuNGE1LjQgNS40IDAgMCAxIDUtMy43eiIgZmlsbD0iI0VBNDMzNSIgZmlsbC1ydWxlPSJub256ZXJvIi8+PHBhdGggZD0iTTAgMGgxOHYxOEgweiIvPjwvZz48L3N2Zz4=);
    background-repeat: no-repeat;
    background-position: 12px 11px;
    margin-bottom: 10px;
}

.login-with-github-button {
    cursor: pointer;

    padding: 12px 16px 12px 42px;

    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Oxygen, Ubuntu, Cantarell, "Fira Sans", "Droid Sans", "Helvetica Neue", sans-serif;

    background-image: url(data:image/svg+xml;base64,PD94bWwgdmVyc2lvbj0iMS4wIj8+PHN2ZyBmaWxsPSIjMDAwMDAwIiB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciICB2aWV3Qm94PSIwIDAgMzAgMzAiPiAgICA8cGF0aCBkPSJNMTUsM0M4LjM3MywzLDMsOC4zNzMsMywxNWMwLDUuNjIzLDMuODcyLDEwLjMyOCw5LjA5MiwxMS42M0MxMi4wMzYsMjYuNDY4LDEyLDI2LjI4LDEyLDI2LjA0N3YtMi4wNTEgYy0wLjQ4NywwLTEuMzAzLDAtMS41MDgsMGMtMC44MjEsMC0xLjU1MS0wLjM1My0xLjkwNS0xLjAwOWMtMC4zOTMtMC43MjktMC40NjEtMS44NDQtMS40MzUtMi41MjYgYy0wLjI4OS0wLjIyNy0wLjA2OS0wLjQ4NiwwLjI2NC0wLjQ1MWMwLjYxNSwwLjE3NCwxLjEyNSwwLjU5NiwxLjYwNSwxLjIyMmMwLjQ3OCwwLjYyNywwLjcwMywwLjc2OSwxLjU5NiwwLjc2OSBjMC40MzMsMCwxLjA4MS0wLjAyNSwxLjY5MS0wLjEyMWMwLjMyOC0wLjgzMywwLjg5NS0xLjYsMS41ODgtMS45NjJjLTMuOTk2LTAuNDExLTUuOTAzLTIuMzk5LTUuOTAzLTUuMDk4IGMwLTEuMTYyLDAuNDk1LTIuMjg2LDEuMzM2LTMuMjMzQzkuMDUzLDEwLjY0Nyw4LjcwNiw4LjczLDkuNDM1LDhjMS43OTgsMCwyLjg4NSwxLjE2NiwzLjE0NiwxLjQ4MUMxMy40NzcsOS4xNzQsMTQuNDYxLDksMTUuNDk1LDkgYzEuMDM2LDAsMi4wMjQsMC4xNzQsMi45MjIsMC40ODNDMTguNjc1LDkuMTcsMTkuNzYzLDgsMjEuNTY1LDhjMC43MzIsMC43MzEsMC4zODEsMi42NTYsMC4xMDIsMy41OTQgYzAuODM2LDAuOTQ1LDEuMzI4LDIuMDY2LDEuMzI4LDMuMjI2YzAsMi42OTctMS45MDQsNC42ODQtNS44OTQsNS4wOTdDMTguMTk5LDIwLjQ5LDE5LDIyLjEsMTksMjMuMzEzdjIuNzM0IGMwLDAuMTA0LTAuMDIzLDAuMTc5LTAuMDM1LDAuMjY4QzIzLjY0MSwyNC42NzYsMjcsMjAuMjM2LDI3LDE1QzI3LDguMzczLDIxLjYyNywzLDE1LDN6Ii8+PC9zdmc+);
    background-repeat: no-repeat;
    background-position: 1px 1px;
}

.error {
    color: red;
}

.text {
    color: var(--font-color);
    font-size: large;
}
</style>
