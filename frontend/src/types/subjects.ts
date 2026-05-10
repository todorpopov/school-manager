export interface Subject {
    [key: string]: unknown
    subject_id: number
    subject_name: string
}

export interface CreateSubjectPayload {
    subject_name: string
}
