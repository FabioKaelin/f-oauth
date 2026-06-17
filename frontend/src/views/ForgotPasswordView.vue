<template>
    <div class="forgot-password">
        <h1>Passwort vergessen</h1>

        <div v-if="!submitted">
            <span class="text">Gib deine E-Mail-Adresse ein und wir senden dir einen Link zum Zurücksetzen deines Passworts.</span>
            <br />
            <br />
            <input v-model="email" class="textInput" type="email" placeholder="Email" :disabled="loading" />
            <br />
            <span v-if="error != ''" class="error">{{ error }}</span>
            <br v-if="error != ''" />
            <button v-if="emailValid" @click="submit" :disabled="loading">
                <span v-if="loading">Wird gesendet...</span>
                <span v-else>Reset-Link senden</span>
            </button>
        </div>

        <div v-else>
            <span class="text">✅ Falls diese E-Mail-Adresse registriert ist, erhältst du in Kürze einen Reset-Link.</span>
            <br />
            <br />
            <span class="text">Überprüfe deinen Posteingang und deinen Spam-Ordner.</span>
        </div>

        <br />
        <br />
        <router-link to="/login" class="text">Zurück zum Login</router-link>
    </div>
</template>

<script lang="ts">
import { defineComponent } from "vue"
import axios from "axios"
import { getAxiosConfigMethod } from "../func"

export default defineComponent({
    name: "ForgotPasswordView",
    data() {
        return {
            email: "",
            loading: false,
            submitted: false,
            error: ""
        }
    },
    computed: {
        emailValid() {
            const emailRegex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/
            return this.email.match(emailRegex) !== null
        }
    },
    methods: {
        submit() {
            if (!this.emailValid) return
            this.loading = true
            this.error = ""
            axios
                .request(getAxiosConfigMethod("/password/forgot", "POST", { email: this.email }))
                .then(() => {
                    this.submitted = true
                })
                .catch(() => {
                    this.error = "Ein Fehler ist aufgetreten. Bitte versuche es erneut."
                })
                .finally(() => {
                    this.loading = false
                })
        }
    }
})
</script>

<style scoped>
.forgot-password {
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
