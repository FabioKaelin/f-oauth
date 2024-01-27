export interface User {
    id: string
    name: string
    email: string
    role: string
    privileges: number
    provider: string
    photo: string
    verified: boolean
}

export interface UserMinimal {
    id: string
    name: string
    privileges: number
}
