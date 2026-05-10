import React, { useEffect } from 'react'
import { ResourceManager } from '../../../components/ResourceManager'
import type { FieldConfig } from '../../../components/ResourceManager'
import { Toast } from '../../../components/Toast'
import type { Term } from '../../../types/terms.ts'
import { useGetTerms, useCreateTerm, useDeleteTerm } from '../../../hooks/useTerms.ts'
import { useToast } from '../../../hooks/useToast'
import { validateText } from '../../../utils/validators'

const TermManagementPage: React.FC = () => {
    const { data = [], isLoading, error } = useGetTerms()
    const createMutation = useCreateTerm()
    const deleteMutation = useDeleteTerm()

    const { toast, show, dismiss } = useToast()

    useEffect(() => {
        if (error) show(error.message, 'error')
    }, [error])

    useEffect(() => {
        if (createMutation.error) show((createMutation.error as Error).message, 'error')
    }, [createMutation.error])

    useEffect(() => {
        if (deleteMutation.error) show((deleteMutation.error as Error).message, 'error')
    }, [deleteMutation.error])

    const TERM_FIELDS: FieldConfig<Term>[] = [
        { key: 'term_id', label: 'ID', type: 'number', hideInForm: true, hideInTable: true },
        {
            key: 'name',
            label: 'Name',
            type: 'text',
            required: true,
            placeholder: 'e.g., Fall 2026, Spring 2026',
            validate: validateText,
        },
        {
            key: 'start_date',
            label: 'Start Date',
            type: 'date',
            required: true,
            renderCell: (value) => {
                // Format date for display
                try {
                    const date = new Date(String(value))
                    return date.toLocaleDateString()
                } catch {
                    return String(value)
                }
            },
        },
        {
            key: 'end_date',
            label: 'End Date',
            type: 'date',
            required: true,
            renderCell: (value) => {
                // Format date for display
                try {
                    const date = new Date(String(value))
                    return date.toLocaleDateString()
                } catch {
                    return String(value)
                }
            },
        },
    ]

    const isMutating = createMutation.isPending || deleteMutation.isPending

    return (
        <main className="max-w-5xl mx-auto px-4 py-10">
            {toast && <Toast message={toast.message} variant={toast.variant} onDismiss={dismiss} />}
            <ResourceManager<Term>
                title="Terms"
                data={data}
                fields={TERM_FIELDS}
                idKey="term_id"
                isLoading={isLoading || isMutating}
                onCreate={async (values) => {
                    await createMutation.mutateAsync({
                        name: String(values.name ?? ''),
                        start_date: String(values.start_date ?? ''),
                        end_date: String(values.end_date ?? ''),
                    })
                }}
                onUpdate={async () => { }}
                onDelete={async (id) => {
                    await deleteMutation.mutateAsync(id as number)
                }}
                hideEdit={true}
            />
        </main>
    )
}

export default TermManagementPage

