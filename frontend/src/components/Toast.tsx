import React from 'react'
import type { ToastVariant } from '../hooks/useToast'

interface ToastProps {
  message: string
  variant?: ToastVariant
  onDismiss: () => void
}

const VARIANT_STYLES: Record<ToastVariant, string> = {
  error:   'bg-red-50 border border-red-200 text-red-700 dark:bg-red-900/30 dark:border-red-800 dark:text-red-400',
  success: 'bg-emerald-50 border border-emerald-200 text-emerald-700 dark:bg-emerald-900/30 dark:border-emerald-800 dark:text-emerald-400',
  warning: 'bg-amber-50 border border-amber-200 text-amber-700 dark:bg-amber-900/30 dark:border-amber-800 dark:text-amber-400',
  info:    'bg-indigo-50 border border-indigo-200 text-indigo-700 dark:bg-indigo-900/30 dark:border-indigo-800 dark:text-indigo-400',
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
