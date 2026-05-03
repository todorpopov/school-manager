import type React from 'react';
import { useState } from 'react'

interface ConfirmDeleteProps {
    onConfirm: () => Promise<void>
    onCancel: () => void
}

export const ConfirmDelete: React.FC<ConfirmDeleteProps> = ({ onConfirm, onCancel }) => {
    const [deleting, setDeleting] = useState(false)

    const handle = async () => {
        setDeleting(true)
        try {
            await onConfirm()
        } finally {
            setDeleting(false)
        }
    }

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center">
            <div className="absolute inset-0 bg-black/40 backdrop-blur-sm" onClick={onCancel} />
            <div className="relative bg-white dark:bg-slate-800 rounded-xl shadow-xl border border-slate-200 dark:border-slate-700 w-full max-w-sm mx-4 p-6 flex flex-col gap-5">
                <div>
                    <p className="text-base font-semibold text-slate-800 dark:text-slate-100 mb-1">Delete record</p>
                    <p className="text-sm text-slate-500 dark:text-slate-400">Are you sure you want to delete this record?</p>
                </div>
                <div className="flex justify-end gap-2">
                    <button
                        onClick={onCancel}
                        disabled={deleting}
                        className="px-3 py-1.5 text-sm font-medium rounded-md border border-slate-300 dark:border-slate-600 text-slate-600 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-700 disabled:opacity-50 cursor-pointer bg-transparent transition-colors"
                    >
                        Cancel
                    </button>
                    <button
                        onClick={handle}
                        disabled={deleting}
                        className="px-3 py-1.5 text-sm font-medium rounded-md bg-red-500 hover:bg-red-600 text-white disabled:opacity-50 cursor-pointer border-none transition-colors"
                    >
                        {deleting ? 'Deleting…' : 'Delete'}
                    </button>
                </div>
            </div>
        </div>
    )
}
