export interface Grade {
    grade_id: number
    grade_value: number
    grade_date: string
    student: { student_id: number; first_name: string; last_name: string }
    curriculum: {
        curriculum_id: number
        subject: { subject_name: string }
        teacher_id: number | null
        term: { name: string }
    }
}

export interface Absence {
    absence_id: number
    absence_date: string
    is_excused: boolean
    student: { student_id: number; first_name: string; last_name: string }
    curriculum: {
        curriculum_id: number
        subject: { subject_name: string }
        term: { name: string }
    }
}

export interface CreateGradePayload {
    student_id: number
    curriculum_id: number
    grade_value: number
    grade_date: string
}

export interface CreateAbsencePayload {
    student_id: number
    curriculum_id: number
    absence_date: string
    is_excused: boolean
}

export interface BulkCreateGradesPayload {
    entries: CreateGradePayload[]
}

export interface BulkCreateAbsencesPayload {
    entries: CreateAbsencePayload[]
}
