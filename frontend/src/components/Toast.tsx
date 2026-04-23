import type React from 'react'
import type { ToastVariant } from '../hooks/useToast'

interface ToastProps {
    message: string
    variant?: ToastVariant
    onDismiss: () => void
}

const VARIANT_STYLES: Record<ToastVariant, string> = {
    error:   'bg-red-100 border border-red-300 text-red-800 dark:bg-red-900/80 dark:border-red-700 dark:text-red-200',
    success: 'bg-emerald-100 border border-emerald-300 text-emerald-800 dark:bg-emerald-900/80 dark:border-emerald-700 dark:text-emerald-200',
    warning: 'bg-amber-100 border border-amber-300 text-amber-800 dark:bg-amber-900/80 dark:border-amber-700 dark:text-amber-200',
    info:    'bg-indigo-100 border border-indigo-300 text-indigo-800 dark:bg-indigo-900/80 dark:border-indigo-700 dark:text-indigo-200',
}

export const Toast: React.FC<ToastProps> = ({ message, variant = 'error', onDismiss }) => {
    return (
        <div className={`fixed top-5 left-1/2 -translate-x-1/2 z-50 flex items-center gap-3 px-4 py-3 rounded-lg shadow-sm text-sm w-max max-w-sm ${VARIANT_STYLES[variant]}`}>
            <span className="flex-1">{message}</span>
            <button
                onClick={onDismiss}
                className="shrink-0 opacity-70 hover:opacity-100 bg-transparent border-none text-inherit cursor-pointer text-base leading-none"
            >
                ✕
            </button>
        </div>
    )
}
