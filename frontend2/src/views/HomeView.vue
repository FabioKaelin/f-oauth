<template>
    <div class="home">
        <h1>Profile</h1>
        <span> <img :src="me.photo" alt="Profilbild" width="100" height="100" /> <br /> </span>
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
        <br>
        <button
            type="button"
            value="menu"
            class="clickButton"
            @click="
                () => {
                    isShow = true
                    newGroupName = ''
                }
            ">
            Bearbeiten
            <Modal
                v-model="isShow"
                :close="
                    () => {
                        isShow = false
                    }
                ">
                <div class="modal">
                    Name:
                    <input v-model="newGroupName" type="text" />
                    <br />
                    <button class="clickButton" @click="isShow = false">Cancel</button>
                    &ensp;
                    <button
                    class="clickButton"
                        @click="
                            () => {
                                createGroup()
                                isShow = false
                            }
                        ">
                        Erstellen
                    </button>
                </div>
            </Modal>
        </button>
        <br>
        <hr>
        <h2>Applications</h2>
        <!-- link to https://tipp.fabkli.ch as button -->
        <a href="https://tipp.fabkli.ch" target="_blank" rel="noopener noreferrer">
            <button class="clickButton">Tippspiel</button>
        </a>
        <br>
        <br>
        <!-- link to https://tipp.dev.fabkli.ch as button -->
        <a href="https://tipp.dev.fabkli.ch" target="_blank" rel="noopener noreferrer">
            <button class="clickButton">Tippspiel-Dev (Nur für Entwicklung)</button>
        </a>
    </div>
</template>

<script lang="ts">
import { defineComponent } from "vue"
import axios from "axios"
import { User } from "../structs"
import { getAxiosConfig } from "../func"

export default defineComponent({
    name: "HomeView",
    data() {
        return {
            count: 1,
            me: {} as User,
            readableRole: "",
            readableProvider: "",
            isShow: false,
        }
    },
    mounted() {
        axios
            .request(getAxiosConfig("/users/me"))
            .then((response: any) => {
                let me = response.data
                this.me = me
            })
            .catch((error: any) => {
                let fromDirect = this.$route.query.from
                if (fromDirect == undefined || fromDirect == null) {
                    fromDirect = window.location.origin
                    // fromDirect = "http://localhost:5173/profile" // TODO: lh
                }
                let from = fromDirect.toString()
                console.log(from)
                // console.log(error);
                console.log(error.response.status, "not logged in")
                this.$router.push({ name: "login", query: { from: from } })
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
                default:
                    return "Unbekannt"
            }
        }
    }
})
</script>

<style lang="scss" scoped>
// img {
//     width: 100%;
//     max-width: 800px;
//     height: auto;
//     margin: 0 auto;
//     display: block;
// }

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

.clickButton{
    font-size: larger;
}


.modal {
    // width: 300px;
    padding: 30px;
    border-radius: 10px;
    box-sizing: border-box;
    background-color: rgb(15, 90, 77);
    font-size: 20px;
    text-align: center;
    color: var(--font-color);
    font-size: normal;
}

</style>
