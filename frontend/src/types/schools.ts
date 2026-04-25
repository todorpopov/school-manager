export interface School {
    [key: string]: unknown
    school_id: number
    school_name: string
    school_address: string
}

export interface CreateSchoolPayload {
    school_name: string
    school_address: string
}

export interface UpdateSchoolPayload {
    school_name: string
    school_address: string
}

