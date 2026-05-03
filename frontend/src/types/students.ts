export interface StudentSchool {
    school_id: number
    school_name: string
    school_address: string
}

export interface StudentClass {
    class_id: number
    grade_level: number
    class_name: string
}

export interface Student {
    [key: string]: unknown
    student_id: number
    user_id: number
    first_name: string
    last_name: string
    email: string
    school: StudentSchool
    class: StudentClass | null
    roles: string[]
}

export interface CreateStudentPayload {
    school_id: number
    first_name: string
    last_name: string
    email: string
    password: string
    class_id?: number | null
}

export interface UpdateStudentPayload {
    school_id?: number | null
    first_name: string
    last_name: string
    email: string
    class_id?: number | null
}


