export interface Principals {
    [key: string]: unknown
    director_id: number
    user_id: number
    first_name: string
    last_name: string
    email: string
    roles: string[]
}

export interface CreatePrincipalPayload {
    first_name: string
    last_name: string
    email: string
    password: string
}

export interface UpdatePrincipalPayload {
    first_name: string
    last_name: string
    email: string
}
