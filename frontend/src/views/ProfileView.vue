<template>
    <div class="profile">
        <h1>Profil</h1>
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
                        <br>
                        <!-- <span class="hint">Bitte kleine Bilder hochladen, ansonsten wird es nicht aktuallisiert. <br>
                            Falls du es bereits versucht hat und es nicht funktioniert hat versuche das Bild
                            zuzuschneiden oder </span>
                        <br>
                        -->
                        <div>
                            <!-- 1 -->
                            <input type="file" accept="image/*" @change="onFileChanged($event)" />
                            <!-- <br> -->
                            <!-- 2<input type="file" accept="image/*;capture=camera" @change="onFileChanged($event)" /> -->
                            <!-- <br> -->
                            <!-- 3<input type="file" accept="image/*,capture=camera" @change="onFileChanged($event)" /> -->
                            <!-- <br> -->
                            <!-- 4<input type="file" accept="image/*" capture @change="onFileChanged($event)" /> -->
                        </div>
                        <!-- <span v-if="uploadStatus!=null">{{ uploadStatus }}</span> -->
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

function dataURLToBlob(dataURL) {
    var BASE64_MARKER = ';base64,';
    if (dataURL.indexOf(BASE64_MARKER) == -1) {
        var parts = dataURL.split(',');
        var contentType = parts[0].split(':')[1];
        var raw = parts[1];

        return new Blob([raw], { type: contentType });
    }

    parts = dataURL.split(BASE64_MARKER);
    contentType = parts[0].split(':')[1];
    raw = window.atob(parts[1]);
    var rawLength = raw.length;

    var uInt8Array = new Uint8Array(rawLength);

    for (var i = 0; i < rawLength; ++i) {
        uInt8Array[i] = raw.charCodeAt(i);
    }

    return new Blob([uInt8Array], { type: contentType });
}

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
            error: "",
            uploadStatus: null as any
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

            this.uploadStatus = target
            if (target && target.files) {
                this.file = target.files[0]
            }
            this.uploadStatus = this.file
        },
        saveImage() {
            this.uploadStatus = "saveImage"
            if (this.file) {
                this.uploadStatus = "file exists"
                if (this.file.type.match(/image.*/)) {
                    console.log('An image has been loaded');
                    // Load the image
                    let fakeThis = this
                    var reader = new FileReader();
                    reader.onload = function (readerEvent) {
                        var image = new Image();
                        console.debug("4 reader loaded")
                        image.onloadeddata = function (imageEvent) {
                            console.debug("10", imageEvent);
                        }
                        image.onload = function (imageEvent) {
                            console.debug("3 image loaded");
                            console.debug("5", imageEvent);

                            // Resize the image
                            var canvas = document.createElement('canvas'),
                                max_size = 200,// TODO : pull max size from a site config
                                width = image.width,
                                height = image.height;
                            if (width > height) {
                                if (width > max_size) {
                                    height *= max_size / width;
                                    width = max_size;
                                }
                            } else {
                                if (height > max_size) {
                                    width *= max_size / height;
                                    height = max_size;
                                }
                            }
                            canvas.width = width;
                            canvas.height = height;
                            canvas.getContext('2d')?.drawImage(image, 0, 0, width, height);
                            var dataUrl = canvas.toDataURL('image/jpeg');
                            var resizedImage = dataURLToBlob(dataUrl);
                            console.debug("6",resizedImage);
                            console.debug("7",resizedImage.size);
                            let formData = new FormData()
                            formData.append("image", resizedImage)
                            // formData.append("image", file)
                            axios.request(getAxiosConfigMethod("/users/me/image", "post", formData)).then((response: any) => {
                                console.debug("8",response)
                                const userId = fakeThis.me.id
                                const backendUrl = import.meta.env.VITE_SERVER_ENDPOINT
                                if (userId) {
                                    fakeThis.imageUrl = `${backendUrl}/api/users/${userId}/image?date=${Date.now()}`
                                } else {
                                    fakeThis.imageUrl = `${backendUrl}/api/users/nouser/image`
                                }
                            })
                        }
                        console.debug("1", readerEvent.target?.result)
                        image.src = readerEvent.target?.result as string;
                        console.debug("9", image.complete)
                    }
                    console.debug("2", this.file)
                    reader.readAsDataURL(this.file);
                }
            } else {
                this.uploadStatus = "file does not exist"

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

.hint {
    font-size: smaller;
    color: var(--font-color);
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
