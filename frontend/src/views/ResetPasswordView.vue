<template>
    <div class="reset-password">
        <h1>Neues Passwort setzen</h1>

        <div v-if="tokenError != ''">
            <span class="error">{{ tokenError }}</span>
            <br />
            <br />
            <router-link to="/forgot-password" class="text">Neuen Reset-Link anfordern</router-link>
        </div>

        <div v-else-if="!success">
            <input v-model="newPassword" class="textInput" type="password" placeholder="Neues Passwort" />
            <br v-if="!passwordValid && newPassword.length != 0" />
            <span v-if="!passwordValid && newPassword.length != 0" style="color: red">
                Das Passwort muss mindestens 8 Zeichen lang sein und mindestens eine Großbuchstabe, eine Kleinbuchstabe, eine Zahl und ein Sonderzeichen enthalten.
            </span>
            <br />
            <input v-model="confirmPassword" class="textInput" type="password" placeholder="Passwort bestätigen" />
            <br v-if="!passwordMatch && confirmPassword.length != 0" />
            <span v-if="!passwordMatch && confirmPassword.length != 0" style="color: red">Die Passwörter stimmen nicht überein</span>
            <br />
            <span v-if="error != ''" class="error">{{ error }}</span>
            <br v-if="error != ''" />
            <button v-if="passwordValid && passwordMatch" :disabled="loading" @click="submit">
                <span v-if="loading">Wird gespeichert...</span>
                <span v-else>Passwort setzen</span>
            </button>
        </div>

        <div v-else>
            <span class="text">✅ Dein Passwort wurde erfolgreich zurückgesetzt.</span>
            <br />
            <br />
            <router-link to="/login" class="text">Zum Login</router-link>
        </div>
    </div>
</template>

<script lang="ts">
import { defineComponent } from "vue"
import axios from "axios"
import { getAxiosConfigMethod } from "../func"

export default defineComponent({
    name: "ResetPasswordView",
    data() {
        return {
            newPassword: "",
            confirmPassword: "",
            loading: false,
            success: false,
            error: "",
            tokenError: "",
            token: ""
        }
    },
    computed: {
        passwordValid() {
            const passwordRegex = /^(?=.*?[A-Z])(?=(.*[a-z]){1,})(?=(.*[\d]){1,})(?=(.*[\W]){1,})(?!.*\s).{8,}$/
            return this.newPassword.match(passwordRegex) !== null
        },
        passwordMatch() {
            return this.newPassword === this.confirmPassword
        }
    },
    mounted() {
        this.token = (this.$route.query.token as string) || ""
        if (!this.token) {
            this.tokenError = "Ungültiger Reset-Link. Bitte fordere einen neuen an."
        }
    },
    methods: {
        submit() {
            if (!this.passwordValid || !this.passwordMatch) return
            this.loading = true
            this.error = ""
            axios
                .request(getAxiosConfigMethod("/password/reset/" + this.token, "POST", { password: this.newPassword }))
                .then(() => {
                    this.success = true
                })
                .catch(err => {
                    const msg: string = err?.response?.data?.message || ""
                    if (msg === "token expired") {
                        this.tokenError = "Dieser Reset-Link ist abgelaufen. Bitte fordere einen neuen an."
                    } else if (msg === "token already used") {
                        this.tokenError = "Dieser Reset-Link wurde bereits verwendet. Bitte fordere einen neuen an."
                    } else if (msg === "invalid token") {
                        this.tokenError = "Dieser Reset-Link ist ungültig. Bitte fordere einen neuen an."
                    } else {
                        this.error = "Ein Fehler ist aufgetreten. Bitte versuche es erneut."
                    }
                })
                .finally(() => {
                    this.loading = false
                })
        }
    }
})
</script>

<style scoped>
.reset-password {
    background-color: var(--color);
    border-radius: 10px;
    padding: 5px;
}

.error {
    color: red;
}

.text {
    color: var(--font-color);
    font-size: large;
}
</style>
