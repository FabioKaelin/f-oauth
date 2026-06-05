<template>
    <div class="profile">
        <h1>Profil</h1>
        <div v-if="error && loaded" class="error">{{ error }}</div>
        <div v-if="!loaded" class="loader"></div>
        <div v-if="loaded">
            <span> <img :src="imageUrl" alt="Profilbild" width="100" height="100" /> <br /> </span>
            <!-- <span> <img :src="me.photo" alt="Profilbild" width="100" height="100" /> <br /> </span> -->
            <table>
                <tbody>
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
                </tbody>
            </table>
            <br />
            <button type="button" value="menu" class="clickButton editButton" @click="openEditModal">
                Bearbeiten
            </button>
            <Modal v-model="isShow" :close="() => {
                if (isSaving) {
                    return
                }
                isShow = false
            }
                ">
                <div class="modal">
                    <div class="modalTitle">Profil bearbeiten</div>
                    <label class="field">
                        <span>Name</span>
                        <input v-model="newUsername" class="profileInput" type="text" />
                    </label>
                    <label class="field">
                        <span>Profilbild</span>
                        <input class="fileInput" type="file" accept="image/*" @change="onFileChanged($event)" />
                    </label>
                    <p class="hint">Bild auswählen, speichern und wir passen Vorschau und Profil automatisch an.</p>
                    <span v-if="uploadError" class="error">{{ uploadError }}</span>
                    <div class="actions">
                        <button class="secondaryButton" :disabled="isSaving" @click="isShow = false">Abbrechen</button>
                        <button class="primaryButton" :disabled="isSaving" @click="updateUser">
                            {{ isSaving ? "Speichern..." : "Aktualisieren" }}
                        </button>
                    </div>
                </div>
            </Modal>
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
import { defineComponent } from "vue"
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
            file: null as File | null,
            imageUrl: "",
            loaded: false,
            error: "",
            uploadStatus: null as any,
            isSaving: false,
            uploadError: ""
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
        openEditModal() {
            this.isShow = true
            this.newUsername = this.me.name
            this.file = null
            this.uploadError = ""
        },
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
        async updateUser() {
            this.isSaving = true
            this.uploadError = ""
            try {
                const response: any = await axios.request(getAxiosConfigMethod("/users/me", "put", { name: this.newUsername }))
                this.me = response.data
                if (this.file) {
                    try {
                        await this.saveImage()
                    } catch (uploadError: any) {
                        this.uploadError = uploadError?.response?.data?.message || uploadError?.message || "Bild konnte nicht gespeichert werden"
                        return
                    }
                }
                this.isShow = false
            } catch (error: any) {
                console.log(error)
                this.uploadError = error?.response?.data?.message || "Fehler beim Aktualisieren des Profils"
            } finally {
                this.isSaving = false
            }
        },
        onFileChanged($event: Event) {
            const target = $event.target as HTMLInputElement

            this.uploadStatus = target
            if (target && target.files) {
                this.file = target.files[0]
            }
            this.uploadStatus = this.file
        },
        async saveImage() {
            this.uploadStatus = "saveImage"
            if (!this.file) {
                return
            }

            if (!this.file.type.match(/image.*/)) {
                throw new Error("Bitte ein gültiges Bild auswählen")
            }

            // Keep uploads user-friendly: cap huge camera images before sending
            const resizedImage = await this.resizeImage(this.file, 1200, 0.85)
            const formData = new FormData()
            formData.append("image", resizedImage, "profile.jpg")

            await axios.request(getAxiosConfigMethod("/users/me/image", "post", formData))

            const userId = this.me.id
            const backendUrl = import.meta.env.VITE_SERVER_ENDPOINT
            if (userId) {
                this.imageUrl = `${backendUrl}/api/users/${userId}/image?date=${Date.now()}`
            } else {
                this.imageUrl = `${backendUrl}/api/users/default/image`
            }
        },

        resizeImage(file: File, maxSize: number, quality: number): Promise<Blob> {
            return new Promise((resolve, reject) => {
                const reader = new FileReader()
                reader.onload = event => {
                    const image = new Image()
                    image.onload = () => {
                        let width = image.width
                        let height = image.height

                        if (width > height) {
                            if (width > maxSize) {
                                height *= maxSize / width
                                width = maxSize
                            }
                        } else if (height > maxSize) {
                            width *= maxSize / height
                            height = maxSize
                        }

                        const canvas = document.createElement("canvas")
                        canvas.width = width
                        canvas.height = height
                        canvas.getContext("2d")?.drawImage(image, 0, 0, width, height)

                        canvas.toBlob(
                            blob => {
                                if (!blob) {
                                    reject(new Error("Bild konnte nicht verarbeitet werden"))
                                    return
                                }
                                resolve(blob)
                            },
                            "image/jpeg",
                            quality
                        )
                    }
                    image.onerror = () => reject(new Error("Bild konnte nicht gelesen werden"))
                    image.src = event.target?.result as string
                }
                reader.onerror = () => reject(new Error("Datei konnte nicht gelesen werden"))
                reader.readAsDataURL(file)
            })
        }
    }
})
</script>

<style lang="scss" scoped>
.profile {
    background: linear-gradient(180deg, color-mix(in srgb, var(--color) 88%, #ffffff 12%), var(--color));
    border-radius: 18px;
    padding: 20px;
    box-shadow: 0 18px 50px rgba(0, 0, 0, 0.18);
}

.hint {
    font-size: smaller;
    color: var(--font-color);
}

.error {
    color: #ff8b8b;
    font-weight: 700;
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
    border: 3px solid color-mix(in srgb, var(--color-half) 55%, white 45%);
    box-shadow: 0 10px 24px rgba(0, 0, 0, 0.18);
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

.clickButton,
.primaryButton,
.secondaryButton {
    font-size: larger;
}

.editButton {
    margin-top: 12px;
}

.modal {
    width: min(92vw, 460px);
    padding: 28px;
    border-radius: 20px;
    box-sizing: border-box;
    background: linear-gradient(180deg, color-mix(in srgb, var(--color-dark) 94%, white 6%), var(--color-dark));
    border: 1px solid color-mix(in srgb, var(--color-half) 40%, white 60%);
    font-size: 20px;
    text-align: left;
    color: var(--font-color);
    font-size: normal;
    box-shadow: 0 24px 70px rgba(0, 0, 0, 0.28);
}

.modalTitle {
    font-size: 1.4rem;
    font-weight: 800;
    margin-bottom: 18px;
    text-align: center;
    letter-spacing: 0.03em;
}

.field {
    display: flex;
    flex-direction: column;
    gap: 8px;
    margin-bottom: 14px;
}

.field span {
    font-size: 0.95rem;
    opacity: 0.9;
}

.profileInput,
.fileInput {
    width: 100%;
    box-sizing: border-box;
    border-radius: 14px;
    border: 1px solid color-mix(in srgb, var(--color-half) 45%, white 55%);
    background-color: color-mix(in srgb, var(--color-dark) 78%, white 22%);
    color: var(--font-color);
    padding: 12px 14px;
    font-size: 1rem;
}

.profileInput:focus,
.fileInput:focus {
    outline: 2px solid color-mix(in srgb, var(--color-half) 65%, white 35%);
    outline-offset: 2px;
}

.fileInput {
    padding: 10px 12px;
}

.actions {
    display: flex;
    gap: 12px;
    justify-content: flex-end;
    margin-top: 18px;
}

.secondaryButton,
.primaryButton {
    min-width: 140px;
    padding: 12px 16px;
    border-radius: 14px;
    border: none;
}

.secondaryButton {
    background-color: color-mix(in srgb, var(--color-half) 50%, black 50%);
}

.primaryButton {
    background: linear-gradient(135deg, color-mix(in srgb, var(--color) 90%, white 10%), color-mix(in srgb, var(--color-half) 75%, white 25%));
    box-shadow: 0 10px 24px rgba(0, 0, 0, 0.2);
}


.applicationlink {
    text-decoration: none;
    font-size: larger;
    font-size: larger;
    color: var(--font-color);
}
</style>
