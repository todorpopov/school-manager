import React, { useEffect } from 'react'
import { ResourceManager } from '../../../components/ResourceManager'
import type { FieldConfig } from '../../../components/ResourceManager'
import { Toast } from '../../../components/Toast'
import type { Subject } from '../../../types/subjects.ts'
import { useGetSubjects, useCreateSubject, useDeleteSubject } from '../../../hooks/useSubjects.ts'
import { useToast } from '../../../hooks/useToast'
import { validateText } from '../../../utils/validators'

const SubjectManagementPage: React.FC = () => {
    const { data = [], isLoading, error } = useGetSubjects()
    const createMutation = useCreateSubject()
    const deleteMutation = useDeleteSubject()

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

    const SUBJECT_FIELDS: FieldConfig<Subject>[] = [
        { key: 'subject_id', label: 'ID', type: 'number', hideInForm: true, hideInTable: true },
        {
            key: 'subject_name',
            label: 'Subject Name',
            type: 'text',
            required: true,
            placeholder: 'e.g., Mathematics, Biology, History',
            validate: validateText,
        },
    ]

    const isMutating = createMutation.isPending || deleteMutation.isPending

    return (
        <main className="max-w-5xl mx-auto px-4 py-10">
            {toast && <Toast message={toast.message} variant={toast.variant} onDismiss={dismiss} />}
            <ResourceManager<Subject>
                title="Subjects"
                data={data}
                fields={SUBJECT_FIELDS}
                idKey="subject_id"
                isLoading={isLoading || isMutating}
                onCreate={async (values) => {
                    await createMutation.mutateAsync({
                        subject_name: String(values.subject_name ?? ''),
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

export default SubjectManagementPage

