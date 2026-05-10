import React, { useEffect } from 'react'
import { ResourceManager } from '../../../components/ResourceManager'
import type { FieldConfig } from '../../../components/ResourceManager'
import { Toast } from '../../../components/Toast'
import type { Class } from '../../../types/classes.ts'
import { useGetClasses, useCreateClass, useDeleteClass } from '../../../hooks/useClasses.ts'
import { useGetSchools } from '../../../hooks/useSchools.ts'
import { useToast } from '../../../hooks/useToast'
import { validateText } from '../../../utils/validators'

const ClassManagementPage: React.FC = () => {
    const { data = [], isLoading, error } = useGetClasses()
    const { data: schools = [] } = useGetSchools()
    const createMutation = useCreateClass()
    const deleteMutation = useDeleteClass()

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

    const schoolOptions = schools.map((s) => ({
        label: s.school_name as string,
        value: s.school_id as number,
    }))

    const gradeOptions = Array.from({ length: 12 }, (_, i) => ({
        label: String(i + 1),
        value: i + 1,
    }))

    const CLASS_FIELDS: FieldConfig<Class>[] = [
        { key: 'class_id', label: 'ID', type: 'number', hideInForm: true, hideInTable: true },
        {
            key: 'school_id',
            label: 'School',
            type: 'select',
            required: true,
            options: schoolOptions,
            renderCell: (value) => {
                const school = schools.find(s => s.school_id === value)
                return school?.school_name as string ?? '—'
            },
        },
        {
            key: 'grade_level',
            label: 'Grade',
            type: 'select',
            required: true,
            options: gradeOptions,
        },
        {
            key: 'class_name',
            label: 'Class Letter',
            type: 'text',
            required: true,
            placeholder: 'e.g., A, B, C',
            validate: validateText,
        },
    ]

    const isMutating = createMutation.isPending || deleteMutation.isPending

    return (
        <main className="max-w-5xl mx-auto px-4 py-10">
            {toast && <Toast message={toast.message} variant={toast.variant} onDismiss={dismiss} />}
            <ResourceManager<Class>
                title="Classes"
                data={data}
                fields={CLASS_FIELDS}
                idKey="class_id"
                isLoading={isLoading || isMutating}
                onCreate={async (values) => {
                    await createMutation.mutateAsync({
                        school_id: Number(values.school_id),
                        grade_level: Number(values.grade_level),
                        class_name: String(values.class_name ?? ''),
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

export default ClassManagementPage

