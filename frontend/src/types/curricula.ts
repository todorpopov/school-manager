export interface CurriculumClass {
    class_id: number
    grade_level: number
    class_name: string
}

export interface CurriculumSubject {
    subject_id: number
    subject_name: string
}

export interface CurriculumTerm {
    term_id: number
    name: string
}

export interface Curriculum {
    [key: string]: unknown
    curriculum_id: number
    class: CurriculumClass
    subject: CurriculumSubject
    teacher_id: number | null
    term: CurriculumTerm
}

export interface CreateCurriculumPayload {
    class_id: number
    subject_id: number
    teacher_id: number
    term_id: number
}

