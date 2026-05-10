export interface Term {
    [key: string]: unknown
    term_id: number
    name: string
    start_date: string
    end_date: string
}

export interface CreateTermPayload {
    name: string
    start_date: string
    end_date: string
}

