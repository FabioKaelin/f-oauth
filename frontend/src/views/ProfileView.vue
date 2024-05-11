<template>
    <div class="profile">
        <h1>Profile</h1>
        <div v-if="error && loaded" class="error">{{ error }}</div>
        <div v-if="!loaded" class="loader"></div>
        <div v-if="loaded">
            <span> <img :src="imageUrl" alt="Profilbild" width="100" height="100" /> <br /> </span>
            <!-- <span> <img :src="me.photo" alt="Profilbild" width="100" height="100" /> <br /> </span> -->
            <table>
                <tr>
                    <td>Name</td>
                    <td>{{ me.name }}</td>
                </tr>
                <tr>
                    <td>Email</td>
                    <td>{{ me.email }}</td>
                </tr>
                <tr>
                    <td>Loginmethode</td>
                    <td>{{ getRealableProvider() }}</td>
                </tr>
                <tr>
                    <td>Rolle</td>
                    <td>{{ getReadableRole() }}</td>
                </tr>
            </table>
            <br />
            <button type="button" value="menu" class="clickButton" @click="() => {
                    isShow = true
                    newUsername = me.name
                    file = null
                }
                ">
                Bearbeiten
                <Modal v-model="isShow" :close="() => {
                        isShow = false
                    }
                    ">
                    <div class="modal">
                        Name:
                        <input v-model="newUsername" type="text" />
                        <br />
                        <div>
                            <input type="file" accept="image/*" capture @change="onFileChanged($event)" />
                        </div>
                        <br />
                        <button class="clickButton" @click="isShow = false">Abbrechen</button>
                        &ensp;
                        <button class="clickButton" @click="() => {
                                updateUser()
                                isShow = false
                            }
                            ">
                            Aktualisieren
                        </button>
                    </div>
                </Modal>
            </button>
        </div>
        <br />
        <br />
        <hr />
        <h2>Applications</h2>
        <!-- link to https://tipp.fabkli.ch as button -->
        <a class="applicationlink" href="https://tipp.fabkli.ch" target="_blank" rel="noopener noreferrer">
            <button class="clickButton">
                <img src="https://tipp.fabkli.ch/favicon.png" alt="" />
                <br />
                Tippspiel
            </button>
        </a>
        <br />
        <br />
        <!-- link to https://tipp.dev.fabkli.ch as button -->
        <!-- <a class="applicationlink" href="https://tipp.dev.fabkli.ch" target="_blank" rel="noopener noreferrer">
            <button class="clickButton">Tippspiel-Dev (Nur für Entwicklung)</button>
        </a> -->
    </div>
</template>

<script lang="ts">
import { defineComponent, ref } from "vue"
import axios from "axios"
import { User } from "../structs"
import { getAxiosConfig, getAxiosConfigMethod } from "../func"

export default defineComponent({
    name: "ProfileView",
    data() {
        return {
            count: 1,
            me: {} as User,
            readableRole: "",
            readableProvider: "",
            isShow: false,
            newUsername: "",
            file: ref<File | null>(),
            imageUrl: "",
            loaded: false,
            error: ""
        }
    },
    mounted() {
        axios
            .request(getAxiosConfig("/users/me"))
            .then((response: any) => {
                let me = response.data
                this.me = me
                const userId = me.id
                this.loaded = true
                const backendUrl = import.meta.env.VITE_SERVER_ENDPOINT
                if (userId) {
                    this.imageUrl = `${backendUrl}/api/users/${userId}/image`
                } else {
                    this.imageUrl = `${backendUrl}/api/users/nouser/image`
                }
            })
            .catch((error: any) => {
                if (error.response.status == 401) {
                    let fromDirect = this.$route.query.from
                    if (fromDirect == undefined || fromDirect == null) {
                        fromDirect = window.location.href
                        console.log("1 fromDirect", fromDirect)
                    }
                    let from = fromDirect.toString()
                    console.log(from)
                    // console.log(error);
                    console.log(error.response.status, "not logged in")
                    this.$router.push({ name: "login", query: { from: from } })
                    return
                }
                this.error = "Fehler beim Laden der Daten. Bitte versuchen Sie es später erneut."
                this.loaded = true
            })
    },
    methods: {
        getReadableRole() {
            let role = this.me.role
            if (role == undefined) return "Unbekannt"
            switch (role.toLowerCase()) {
                case "admin":
                    return "Administrator"
                case "test-admin":
                    return "Test-Administrator"
                case "user":
                    return "Benutzer"
                case "test-user":
                    return "Test-Benutzer"
                default:
                    return "Unbekannt"
            }
        },
        getRealableProvider() {
            let provider = this.me.provider
            if (provider == undefined) return "Unbekannt"
            switch (provider.toLowerCase()) {
                case "local":
                    return "Benutzername und Passwort"
                case "google":
                    return "Google"
                case "github":
                    return "GitHub"
                default:
                    return "Unbekannt"
            }
        },
        updateUser() {
            this.saveImage()
            axios
                .request(getAxiosConfigMethod("/users/me", "put", { name: this.newUsername }))
                .then((response: any) => {
                    let me = response.data
                    this.me = me
                })
                .catch((error: any) => {
                    console.log(error)
                })
        },
        onFileChanged($event: Event) {
            const target = $event.target as HTMLInputElement
            if (target && target.files) {
                this.file = target.files[0]
            }
        },
        saveImage() {
            if (this.file) {
                try {
                    let formData = new FormData()
                    formData.append("image", this.file)
                    axios.request(getAxiosConfigMethod("/users/me/image", "post", formData)).then(() => {
                        const userId = this.me.id
                        const backendUrl = import.meta.env.VITE_SERVER_ENDPOINT
                        if (userId) {
                            this.imageUrl = `${backendUrl}/api/users/${userId}/image?date=${Date.now()}`
                        } else {
                            this.imageUrl = `${backendUrl}/api/users/nouser/image`
                        }
                    })
                } catch (error) {
                    console.error(error)
                    this.file = null
                } finally {
                    console.log("finally")
                }
            }
        }
    }
})
</script>

<style lang="scss" scoped>
.profile {
    background-color: var(--color);
    border-radius: 10px;
    padding: 5px;
}

.error {
    color: red;
}

// img {
//     width: 100%;
//     max-width: 800px;
//     height: auto;
//     margin: 0 auto;
//     display: block;
// }

img {
    border-radius: 40%;
    height: 100px;
}

table {
    font-family: arial, sans-serif;
    border-collapse: collapse;
    /* width: 100%; */
    margin: auto;
}

td,
th {
    border: 1px solid #dddddd;
    text-align: left;
    padding: 8px;
}

tr:nth-child(even) {
    background-color: #dddddd1f;
}

.clickButton {
    font-size: larger;
}

.modal {
    // width: 300px;
    padding: 30px;
    border-radius: 10px;
    box-sizing: border-box;
    background-color: var(--color-dark);
    font-size: 20px;
    text-align: center;
    color: var(--font-color);
    font-size: normal;
}


.applicationlink {
    text-decoration: none;
    font-size: larger;
    font-size: larger;
    color: var(--font-color);
}
</style>
