export interface TeacherSchool {
    school_id: number
    school_name: string
    school_address: string
}

export interface Teacher {
    [key: string]: unknown
    teacher_id: number
    user_id: number
    school: TeacherSchool
    first_name: string
    last_name: string
    email: string
    roles: string[]
}

export interface CreateTeacherPayload {
    school_id: number
    first_name: string
    last_name: string
    email: string
    password: string
}

export interface UpdateTeacherPayload {
    school_id: number
    first_name: string
    last_name: string
    email: string
}

