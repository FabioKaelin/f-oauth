import { createApp } from "vue"
import App from "@/App.vue"
import router from "@/router"

import VueDatePicker  from "@vuepic/vue-datepicker"
import "@vuepic/vue-datepicker/dist/main.css"
import "vue-universal-modal/dist/index.css"
import { createVuetify } from "vuetify"
import * as components from "vuetify/components"
import * as directives from "vuetify/directives"
// import 'vuetify/styles';
// import "../node_modules/vue3-toggle-button/dist/style.css"


const vuetify = createVuetify({
    components,
    directives
})

const app = createApp(App)

import VueUniversalModal from "vue-universal-modal"

app.use(VueUniversalModal, {
    teleportTarget: "#modals"
})

// app.use(vuetify)
app.use(router).component("VueDatePicker", VueDatePicker ).use(vuetify).mount("#app")
