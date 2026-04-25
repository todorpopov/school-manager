import React, { useEffect } from 'react'
import { ResourceManager } from '../../../components/ResourceManager'
import type { FieldConfig } from '../../../components/ResourceManager'
import { Toast } from '../../../components/Toast'
import type { Curriculum } from '../../../types/curricula.ts'
import { useGetCurricula, useCreateCurriculum, useDeleteCurriculum } from '../../../hooks/useCurricula.ts'
import { useGetClasses } from '../../../hooks/useClasses.ts'
import { useGetSubjects } from '../../../hooks/useSubjects.ts'
import { useGetTeachers } from '../../../hooks/useTeachers.ts'
import { useGetTerms } from '../../../hooks/useTerms.ts'
import { useToast } from '../../../hooks/useToast'

const CurriculumManagementPage: React.FC = () => {
    const { data = [], isLoading, error } = useGetCurricula()
    const { data: classes = [] } = useGetClasses()
    const { data: subjects = [] } = useGetSubjects()
    const { data: teachers = [] } = useGetTeachers()
    const { data: terms = [] } = useGetTerms()

    const createMutation = useCreateCurriculum()
    const deleteMutation = useDeleteCurriculum()

    const { toast, show, dismiss } = useToast()

    const classOptions = classes.map((c) => ({
        label: `${c.grade_level}${c.class_name}`,
        value: c.class_id,
    }))
    const subjectOptions = subjects.map((s) => ({
        label: s.subject_name,
        value: s.subject_id,
    }))
    const teacherOptions = teachers.map((t) => ({
        label: `${t.first_name as string} ${t.last_name as string}`,
        value: t.teacher_id,
    }))
    const termOptions = terms.map((t) => ({
        label: t.name,
        value: t.term_id,
    }))

    const CURRICULUM_FIELDS: FieldConfig<Curriculum>[] = [
        { key: 'curriculum_id', label: 'ID', type: 'number', hideInForm: true, hideInTable: true },
        { key: 'class_id', label: 'Class', type: 'select', required: true, options: classOptions,
            renderCell: (_, row) => {
                const cls = row?.class as { class_name: string; grade_level: number } | undefined
                return cls ? `${cls.grade_level} ${cls.class_name}` : '—'
            } },
        { key: 'subject_id', label: 'Subject', type: 'select', required: true, options: subjectOptions,
            renderCell: (_, row) => (row?.subject as { subject_name: string } | undefined)?.subject_name ?? '—' },
        { key: 'teacher_id', label: 'Teacher', type: 'select', required: true, options: teacherOptions,
            renderCell: (value) => {
                const t = teachers.find(t => t.teacher_id === value)
                return t ? `${t.first_name as string} ${t.last_name as string}` : '—'
            },
        },
        { key: 'term_id', label: 'Term', type: 'select', required: true, options: termOptions,
            renderCell: (_, row) => (row?.term as { name: string } | undefined)?.name ?? '—' },
    ]

    useEffect(() => { if (error) show(error.message, 'error') }, [error])
    useEffect(() => { if (createMutation.error) show((createMutation.error as Error).message, 'error') }, [createMutation.error])
    useEffect(() => { if (deleteMutation.error) show((deleteMutation.error as Error).message, 'error') }, [deleteMutation.error])

    const enrichedData = data.map((c) => ({
        ...c,
        class_id:   (c.class   as { class_id:   number } | undefined)?.class_id,
        subject_id: (c.subject as { subject_id: number } | undefined)?.subject_id,
        term_id:    (c.term    as { term_id:    number } | undefined)?.term_id,
    }))

    const isMutating = createMutation.isPending || deleteMutation.isPending

    return (
        <main className="max-w-5xl mx-auto px-4 py-10">
            {toast && <Toast message={toast.message} variant={toast.variant} onDismiss={dismiss} />}
            <ResourceManager<Curriculum>
                title="Curricula"
                data={enrichedData}
                fields={CURRICULUM_FIELDS}
                idKey="curriculum_id"
                isLoading={isLoading || isMutating}
                onCreate={async (values) => {
                    await createMutation.mutateAsync({
                        class_id:   Number(values.class_id),
                        subject_id: Number(values.subject_id),
                        teacher_id: Number(values.teacher_id),
                        term_id:    Number(values.term_id),
                    })
                }}
                onUpdate={async () => { }}
                onDelete={async (id) => { await deleteMutation.mutateAsync(id as number) }}
            />
        </main>
    )
}

export default CurriculumManagementPage
