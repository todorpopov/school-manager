export interface Class {
    [key: string]: unknown
    class_id: number
    school_id: number
    grade_level: number
    class_name: string
}

export interface CreateClassPayload {
    school_id: number
    grade_level: number
    class_name: string
}

