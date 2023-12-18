<template>
    <div>
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
    </div>
</template>

<script lang="ts">
import { defineComponent } from "vue"
import { getAxiosConfig } from "../helper/request"
import axios from "axios"
import { user } from "../helper/types"

export default defineComponent({
    name: "ProfileView",
    // props: {
    //     name: String,
    //     msg: { type: String, required: true }
    // },
    data() {
        return {
            count: 1,
            me: {} as user,
            readableRole: "",
            readableProvider: ""
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
                // console.log(error);
                console.log(error.response.status, "not logged in")
                this.$router.push({ name: "login" })
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

<style scoped>
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
</style>
