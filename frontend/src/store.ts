import { reactive } from "vue"
import { User } from "./structs"

export const store = reactive({
    loggedIn: false,
    user: {} as User,
    userLoaded: false
    // increment() {
    //     this.count++
    // }
})
