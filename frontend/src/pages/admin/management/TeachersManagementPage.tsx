import { useEffect } from 'react'
import { ResourceManager } from '../../../components/ResourceManager'
import type { FieldConfig } from '../../../components/ResourceManager'
import { Toast } from '../../../components/Toast'
import type { Teacher } from '../../../types/teachers.ts'
import { useGetTeachers, useCreateTeacher, useUpdateTeacher, useDeleteTeacher } from '../../../hooks/useTeachers.ts'
import { useToast } from '../../../hooks/useToast'

const TEACHER_FIELDS: FieldConfig<Teacher>[] = [
    { key: 'teacher_id', label: 'ID', type: 'number', hideInForm: true, hideInTable: true },
    { key: 'user_id', label: 'User ID', type: 'number', hideInForm: true, hideInTable: true },
    { key: 'first_name', label: 'First Name', type: 'text', required: true, placeholder: 'First name' },
    { key: 'last_name', label: 'Last Name', type: 'text', required: true, placeholder: 'Last name' },
    { key: 'email', label: 'Email', type: 'email', required: true, placeholder: 'Email' },
    { key: 'password', label: 'Password', type: 'password', required: true, placeholder: 'Password', hideInTable: true },
    { key: 'roles', label: 'Roles', type: 'text', hideInForm: true },
]

export default function TeachersManagementPage() {
    const { data = [], isLoading, error } = useGetTeachers()
    const createMutation = useCreateTeacher()
    const updateMutation = useUpdateTeacher()
    const deleteMutation = useDeleteTeacher()

    const { toast, show, dismiss } = useToast()

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

    return (
        <main className="max-w-5xl mx-auto px-4 py-10">
            {toast && <Toast message={toast.message} variant={toast.variant} onDismiss={dismiss} />}
            <ResourceManager<Teacher>
                title="Teachers"
                data={data}
                fields={TEACHER_FIELDS}
                idKey="teacher_id"
                isLoading={isLoading || isMutating}
                onCreate={async (values) => {
                    await createMutation.mutateAsync({
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
