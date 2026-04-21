import React, { useState } from 'react'

interface ConfirmDeleteProps {
  label: string
  onConfirm: () => Promise<void>
  onCancel: () => void
}

export const ConfirmDelete: React.FC<ConfirmDeleteProps> = ({ label, onConfirm, onCancel }) => {
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
    <div className="flex flex-col gap-5 p-5">
      <p className="text-sm text-slate-700 dark:text-slate-200">
        Delete <strong>{label}</strong>? This action cannot be undone.
      </p>
      <div className="flex justify-end gap-2">
        <button
          onClick={onCancel}
          disabled={deleting}
          className="px-4 py-2 text-sm font-medium rounded-md border border-slate-300 dark:border-slate-600 text-slate-600 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-700 disabled:opacity-50 cursor-pointer bg-transparent transition-colors"
        >
          Cancel
        </button>
        <button
          onClick={handle}
          disabled={deleting}
          className="px-4 py-2 text-sm font-medium rounded-md bg-red-500 hover:bg-red-600 text-white disabled:opacity-50 cursor-pointer border-none transition-colors"
        >
          {deleting ? 'Deleting…' : 'Delete'}
        </button>
      </div>
    </div>
  )
}
