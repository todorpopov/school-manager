import { useState, useCallback, useRef } from 'react'

export type ToastVariant = 'error' | 'success' | 'warning' | 'info'

export interface ToastState {
    message: string
    variant: ToastVariant
}

export function useToast(durationMs = 4000) {
    const [toast, setToast] = useState<ToastState | null>(null)
    const timerRef = useRef<ReturnType<typeof setTimeout> | null>(null)

    const show = useCallback((message: string, variant: ToastVariant = 'error') => {
        if (timerRef.current) {
            clearTimeout(timerRef.current)
        }

        setToast({ message, variant })

        timerRef.current = setTimeout(() => {
            setToast(null)
        }, durationMs)
    }, [durationMs])

    const dismiss = useCallback(() => {
        if (timerRef.current) {
            clearTimeout(timerRef.current)
        }
        setToast(null)
    }, [])

    return { toast, show, dismiss }
}
