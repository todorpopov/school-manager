import axios from 'axios'

interface ApiErrorResponse {
    error: boolean
    message: string
    data?: Record<string, string>[] | null
}

export function parseApiError(err: unknown): string {
    if (!axios.isAxiosError(err)) {
        return err instanceof Error ? err.message : 'An unexpected error occurred'
    }

    const response = err.response?.data as ApiErrorResponse | undefined

    if (!response) {
        return err.message ?? 'Network error — please check your connection'
    }

    const fieldErrors = response.data
    if (Array.isArray(fieldErrors) && fieldErrors.length > 0) {
        const messages = fieldErrors
            .flatMap((obj) => Object.entries(obj))
            .map(([field, msg]) => `${humaniseField(field)}: ${msg}`)
            .join(' · ')
        if (messages) return messages
    }

    return response.message || 'An unexpected error occurred'
}

function humaniseField(field: string): string {
    return field.replace(/_/g, ' ').replace(/\b\w/g, (c) => c.toUpperCase())
}

