import React, { useEffect } from 'react'
import { ResourceManager } from '../../../components/ResourceManager'
import type { FieldConfig } from '../../../components/ResourceManager'
import { Toast } from '../../../components/Toast'
import type { Principals } from '../../../types/principals.ts'
import { useGetDirectors, useCreateDirector, useUpdateDirector, useDeleteDirector } from '../../../hooks/usePrincipals.ts'
import { useGetSchools } from '../../../hooks/useSchools.ts'
import { useToast } from '../../../hooks/useToast'
import { validateName, validateEmail, validatePassword } from '../../../utils/validators'

const PrincipalManagementPage: React.FC = () => {
    const { data = [], isLoading, error } = useGetDirectors()
    const { data: schools = [] } = useGetSchools()
    const createMutation = useCreateDirector()
    const updateMutation = useUpdateDirector()
    const deleteMutation = useDeleteDirector()

    const { toast, show, dismiss } = useToast()

    const schoolOptions = schools.map((s) => ({
        label: s.school_name as string,
        value: s.school_id as number,
    }))

    const PRINCIPAL_FIELDS: FieldConfig<Principals>[] = [
        { key: 'director_id', label: 'ID', type: 'number', hideInForm: true, hideInTable: true },
        { key: 'user_id', label: 'User ID', type: 'number', hideInForm: true, hideInTable: true },
        { key: 'school', label: 'School', type: 'text', hideInForm: true,
            renderCell: (value) => {
                const s = value as { school_name: string } | undefined
                return s?.school_name ?? '—'
            },
        },
        { key: 'school_id', label: 'School', type: 'select', required: true, options: schoolOptions, hideInTable: true },
        { key: 'first_name', label: 'First Name', type: 'text', required: true, placeholder: 'First name', validate: validateName },
        { key: 'last_name', label: 'Last Name', type: 'text', required: true, placeholder: 'Last name', validate: validateName },
        { key: 'email', label: 'Email', type: 'email', required: true, placeholder: 'Email', validate: validateEmail },
        { key: 'password', label: 'Password', type: 'password', required: true, placeholder: 'Password', hideInTable: true, validate: validatePassword },
        { key: 'roles', label: 'Roles', type: 'text', hideInForm: true, hideInTable: true },
    ]

    useEffect(() => {
        if (error) show(error.message, 'error')
    }, [error])

    useEffect(() => {
        if (createMutation.error) show((createMutation.error as Error).message, 'error')
    }, [createMutation.error])

    useEffect(() => {
        if (updateMutation.error) show((updateMutation.error as Error).message, 'error')
    }, [updateMutation.error])

    useEffect(() => {
        if (deleteMutation.error) show((deleteMutation.error as Error).message, 'error')
    }, [deleteMutation.error])

    const isMutating = createMutation.isPending || updateMutation.isPending || deleteMutation.isPending

    const enrichedData = data.map((d) => ({
        ...d,
        school_id: (d.school as { school_id: number } | undefined)?.school_id,
    }))

    return (
        <main className="max-w-5xl mx-auto px-4 py-10">
            {toast && <Toast message={toast.message} variant={toast.variant} onDismiss={dismiss} />}
            <ResourceManager<Principals>
                title="Directors"
                data={enrichedData}
                fields={PRINCIPAL_FIELDS}
                idKey="director_id"
                isLoading={isLoading || isMutating}
                onCreate={async (values) => {
                    await createMutation.mutateAsync({
                        school_id: Number(values.school_id),
                        first_name: String(values.first_name ?? ''),
                        last_name: String(values.last_name ?? ''),
                        email: String(values.email ?? ''),
                        password: String(values.password ?? ''),
                    })
                }}
                onUpdate={async (id, values) => {
                    await updateMutation.mutateAsync({
                        id: id as number,
                        payload: {
                            school_id: Number(values.school_id),
                            first_name: String(values.first_name ?? ''),
                            last_name: String(values.last_name ?? ''),
                            email: String(values.email ?? ''),
                        },
                    })
                }}
                onDelete={async (id) => { await deleteMutation.mutateAsync(id as number) }}
            />
        </main>
    )
}

export default PrincipalManagementPage
