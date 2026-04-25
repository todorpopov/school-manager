import { useEffect } from 'react'
import { ResourceManager } from '../../../components/ResourceManager'
import type { FieldConfig } from '../../../components/ResourceManager'
import { Toast } from '../../../components/Toast'
import type { School } from '../../../types/schools.ts'
import { useGetSchools, useCreateSchool, useUpdateSchool, useDeleteSchool } from '../../../hooks/useSchools.ts'
import { useToast } from '../../../hooks/useToast'

const SCHOOL_FIELDS: FieldConfig<School>[] = [
    { key: 'school_id', label: 'ID', type: 'number', hideInForm: true, hideInTable: true },
    { key: 'school_name', label: 'Name', type: 'text', required: true, placeholder: 'School name' },
    { key: 'school_address', label: 'Address', type: 'text', required: true, placeholder: 'School address' },
]

export default function SchoolManagementPage() {
    const { data = [], isLoading, error } = useGetSchools()
    const createMutation = useCreateSchool()
    const updateMutation = useUpdateSchool()
    const deleteMutation = useDeleteSchool()

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
            <ResourceManager<School>
                title="Schools"
                data={data}
                fields={SCHOOL_FIELDS}
                idKey="school_id"
                isLoading={isLoading || isMutating}
                onCreate={async (values) => {
                    await createMutation.mutateAsync({
                        school_name: String(values.school_name ?? ''),
                        school_address: String(values.school_address ?? ''),
                    })
                }}
                onUpdate={async (id, values) => {
                    await updateMutation.mutateAsync({
                        id: id as number,
                        payload: {
                            school_name: String(values.school_name ?? ''),
                            school_address: String(values.school_address ?? ''),
                        },
                    })
                }}
                onDelete={async (id) => { await deleteMutation.mutateAsync(id as number) }}
            />
        </main>
    )
}
