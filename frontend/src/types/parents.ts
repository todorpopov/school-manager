export interface Parent {
    [key: string]: unknown
    parent_id: number
    user_id: number
    first_name: string
    last_name: string
    email: string
    roles: string[]
}

export interface CreateParentPayload {
    first_name: string
    last_name: string
    email: string
    password: string
}

export interface UpdateParentPayload {
    first_name: string
    last_name: string
    email: string
}

