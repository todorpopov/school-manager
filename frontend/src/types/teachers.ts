export interface Teacher {
    [key: string]: unknown
    teacher_id: number
    user_id: number
    first_name: string
    last_name: string
    email: string
    roles: string[]
}

export interface CreateTeacherPayload {
    first_name: string
    last_name: string
    email: string
    password: string
}

export interface UpdateTeacherPayload {
    first_name: string
    last_name: string
    email: string
}

