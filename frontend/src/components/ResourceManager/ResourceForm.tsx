import type React from 'react';
import { useState, useEffect } from 'react'
import type { FieldConfig, SelectOption } from './types'
import { Field } from '../Field'
import { inputCls } from '../../styles/formStyles'

interface ResourceFormProps<T extends { [key: string]: unknown }> {
    fields: FieldConfig<T>[]
    initialValues?: Partial<T>
    onSubmit: (values: Partial<T>) => Promise<void>
    onCancel: () => void
    submitLabel: string
    isLoading: boolean
}

function getDefaultValues<T extends { [key: string]: unknown }>(
    fields: FieldConfig<T>[],
    initial?: Partial<T>
): Partial<T> {
    const defaults: Partial<T> = {}

    fields
        .filter((f) => !f.hideInForm)
        .forEach((f) => {
            if (initial && initial[f.key] !== undefined) {
                defaults[f.key] = initial[f.key]
            } else if (f.type === 'multiselect') {
                defaults[f.key] = [] as unknown as T[keyof T]
            } else {
                defaults[f.key] = '' as unknown as T[keyof T]
            }
        })

    return defaults
}

export function ResourceForm<T extends { [key: string]: unknown }>({
    fields,
    initialValues,
    onSubmit,
    onCancel,
    submitLabel,
    isLoading,
}: ResourceFormProps<T>) {
    const visibleFields = fields.filter((f) => !f.hideInForm)
    const [values, setValues] = useState<Partial<T>>(() => getDefaultValues(visibleFields, initialValues))
    const [errors, setErrors] = useState<Partial<Record<keyof T, string>>>({})
    const [submitting, setSubmitting] = useState(false)

    useEffect(() => {
        setValues(getDefaultValues(visibleFields, initialValues))
        setErrors({})
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [initialValues])

    const validate = (): boolean => {
        const newErrors: Partial<Record<keyof T, string>> = {}
        visibleFields.forEach((f) => {
            if (f.required) {
                const val = values[f.key]
                const isEmpty = val === undefined || val === '' || (Array.isArray(val) && val.length === 0)
                if (isEmpty) newErrors[f.key] = `${f.label} is required`
            }
        })
        setErrors(newErrors)
        return Object.keys(newErrors).length === 0
    }

    const handleChange = (key: keyof T, value: unknown) => {
        setValues((prev) => ({ ...prev, [key]: value as T[keyof T] }))
        if (errors[key]) setErrors((prev) => ({ ...prev, [key]: undefined }))
    }

    const handleMultiSelect = (key: keyof T, option: string, checked: boolean) => {
        const current = (values[key] as string[]) ?? []
        const updated = checked ? [...current, option] : current.filter((v) => v !== option)
        handleChange(key, updated)
    }

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault()
        if (!validate()) return
        setSubmitting(true)
        try {
            await onSubmit(values)
        } finally {
            setSubmitting(false)
        }
    }

    const renderField = (field: FieldConfig<T>) => {
        const value = values[field.key]
        const hasError = !!errors[field.key]

        switch (field.type) {
            case 'textarea':
                return (
                    <textarea
                        id={String(field.key)}
                        value={(value as string) ?? ''}
                        onChange={(e) => handleChange(field.key, e.target.value)}
                        placeholder={field.placeholder}
                        rows={3}
                        className={`${inputCls(hasError)} resize-y min-h-[72px]`}
                    />
                )
            case 'select':
                return (
                    <select
                        id={String(field.key)}
                        value={(value as string | number) ?? ''}
                        onChange={(e) => handleChange(field.key, e.target.value)}
                        className={inputCls(hasError)}
                    >
                        <option value="">— Select —</option>
                        {field.options?.map((opt: SelectOption) => (
                            <option key={String(opt.value)} value={opt.value}>{opt.label}</option>
                        ))}
                    </select>
                )
            case 'multiselect':
                return (
                    <div className={`flex flex-wrap gap-2 p-2 border rounded-md bg-white dark:bg-slate-900 justify-center ${hasError ? 'border-red-400' : 'border-slate-300 dark:border-slate-600'}`}>
                        {field.options?.map((opt: SelectOption) => {
                            const checked = ((value as string[]) ?? []).includes(String(opt.value))
                            return (
                                <label key={String(opt.value)} className="flex items-center gap-1.5 text-sm text-slate-700 dark:text-slate-200 bg-slate-100 dark:bg-slate-700 px-3 py-1 rounded-full cursor-pointer">
                                    <input
                                        type="checkbox"
                                        checked={checked}
                                        onChange={(e) => handleMultiSelect(field.key, String(opt.value), e.target.checked)}
                                        className="accent-indigo-500"
                                    />
                                    {opt.label}
                                </label>
                            )
                        })}
                    </div>
                )
            default:
                return (
                    <input
                        id={String(field.key)}
                        type={field.type}
                        value={(value as string | number) ?? ''}
                        onChange={(e) => handleChange(field.key, e.target.value)}
                        placeholder={field.placeholder}
                        className={inputCls(hasError)}
                    />
                )
        }
    }

    const busy = submitting || isLoading

    return (
        <form onSubmit={handleSubmit} noValidate className="flex flex-col gap-4 p-5">
            {visibleFields.map((field) => (
                <Field
                    key={String(field.key)}
                    label={field.label}
                    error={errors[field.key] as string | undefined}
                    required={field.required}
                >
                    {renderField(field)}
                </Field>
            ))}

            <div className="flex justify-end gap-2 pt-2">
                <button
                    type="button"
                    onClick={onCancel}
                    disabled={busy}
                    className="px-4 py-2 text-sm font-medium rounded-md border border-slate-300 dark:border-slate-600 text-slate-600 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-700 disabled:opacity-50 cursor-pointer bg-transparent transition-colors"
                >
                    Cancel
                </button>
                <button
                    type="submit"
                    disabled={busy}
                    className="px-4 py-2 text-sm font-medium rounded-md bg-indigo-500 hover:bg-indigo-600 text-white disabled:opacity-50 cursor-pointer border-none transition-colors"
                >
                    {busy ? 'Saving…' : submitLabel}
                </button>
            </div>
        </form>
    )
}
