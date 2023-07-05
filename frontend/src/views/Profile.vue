<template>
    <div>
        Profile
        <hr>
        Name: {{ me.name }} <br>
        Email: {{ me.email }} <br>
        Loginmethode: {{ me.provider }} <br>
        Rolle: {{ getReadableRole() }} <br>
        Image: <img :src="me.photo" alt="Profilbild" width="100px" height="100px"> <br>
        Image: {{ me.photo }} <br>
    </div>
</template>

<script lang="ts">

import { defineComponent } from 'vue'
import { getAxiosConfig } from '../helper/request';
import axios from 'axios'
import { user } from '../helper/types';

export default defineComponent({
    name: 'Profile',
    // props: {
    //     name: String,
    //     msg: { type: String, required: true }
    // },
    data() {
        return {
            count: 1,
            me : {} as user
        }
    },
    methods: {
        getReadableRole() {
            let role = this.me.role;
            switch (role) {
                case "admin":
                    return "Administrator";
                case "test-admin":
                    return "Test-Administrator";
                case "user":
                    return "Benutzer";
                case "test-user":
                    return "Test-Benutzer";
                default:
                    return "Unbekannt";
            }
        },
        getRealableProvider(){
            let provider = this.me.provider;
            switch (provider) {
                case "local":
                    return "Benutzername und Passwort";
                case "google":
                    return "Google";
                default:
                    return "Unbekannt";
            }
        }
    },
    mounted() {
        axios.request(getAxiosConfig("/users/me"))
            .then((response: any) => {
                let me = response.data;
                this.me = me;
            }).catch((error: any) => {
                console.log(error);
                this.$router.push({ name: 'login' });
            });
    }
})

</script>