import React from 'react'

interface FieldProps {
  label: string
  error?: string
  required?: boolean
  children: React.ReactNode
}

export const Field: React.FC<FieldProps> = ({ label, error, required, children }) => {
  return (
    <div className="flex flex-col gap-1">
      <label className="text-sm font-medium text-slate-700 dark:text-slate-200">
        {label}
        {required && <span className="text-red-500 ml-0.5">*</span>}
      </label>
      {children}
      {error && <span className="text-xs text-red-500">{error}</span>}
    </div>
  )
}

