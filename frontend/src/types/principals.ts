export interface DirectorSchool {
    school_id: number
    school_name: string
    school_address: string
}

export interface Principals {
    [key: string]: unknown
    director_id: number
    user_id: number
    school: DirectorSchool
    first_name: string
    last_name: string
    email: string
    roles: string[]
}

export interface CreatePrincipalPayload {
    school_id: number
    first_name: string
    last_name: string
    email: string
    password: string
}

export interface UpdatePrincipalPayload {
    school_id: number
    first_name: string
    last_name: string
    email: string
}
