export interface Director {
    [key: string]: unknown
    director_id: number
    user_id: number
    first_name: string
    last_name: string
    email: string
    roles: string[]
}

export interface CreateDirectorPayload {
    first_name: string
    last_name: string
    email: string
    password: string
}

export interface UpdateDirectorPayload {
    first_name: string
    last_name: string
    email: string
}
