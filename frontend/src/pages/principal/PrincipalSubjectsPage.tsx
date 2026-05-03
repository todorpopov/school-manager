import React, { useEffect, useState } from 'react'
import { ResourceManager } from '../../components/ResourceManager'
import type { FieldConfig } from '../../components/ResourceManager'
import { Toast } from '../../components/Toast'
import type { Subject } from '../../types/subjects.ts'
import { useGetSubjects } from '../../hooks/useSubjects.ts'
import { useToast } from '../../hooks/useToast'
import axiosInstance from '../../utils/axiosConfig'

const API_URL = import.meta.env.VITE_API_URL as string + '/api'

const DirectorSubjectsPage: React.FC = () => {
    const { data = [], isLoading, error } = useGetSubjects()
    const { toast, show, dismiss } = useToast()
    const [subjectTeachersMap, setSubjectTeachersMap] = useState<Record<number, { first_name: string; last_name: string }[]>>({})

    useEffect(() => {
        if (!data.length) return
        const fetchAll = async () => {
            const entries = await Promise.all(
                data.map(async (s) => {
                    try {
                        const res = await axiosInstance.get<{ data: { teacher_id: number; first_name: string; last_name: string }[] }>(
                            `${API_URL}/teacher-subject/subject/${s.subject_id}/teachers`
                        )
                        return [s.subject_id, res.data.data ?? []] as [number, typeof res.data.data]
                    } catch {
                        return [s.subject_id, []] as [number, []]
                    }
                })
            )
            setSubjectTeachersMap(Object.fromEntries(entries))
        }
        fetchAll()
    }, [data])

    const SUBJECT_FIELDS: FieldConfig<Subject>[] = [
        { key: 'subject_id', label: 'ID', type: 'number', hideInForm: true, hideInTable: true },
        { key: 'subject_name', label: 'Subject Name', type: 'text' },
        { key: 'teacher_names', label: 'Teachers', type: 'text',
            renderCell: (value) => (value as string) || '—',
        },
    ]

    useEffect(() => { if (error) show(error.message, 'error') }, [error])

    const enrichedData = data.map((s) => ({
        ...s,
        teacher_names: (subjectTeachersMap[s.subject_id] ?? [])
            .map(t => `${t.first_name} ${t.last_name}`)
            .join(', '),
    }))

    return (
        <main className="max-w-5xl mx-auto px-4 py-10">
            {toast && <Toast message={toast.message} variant={toast.variant} onDismiss={dismiss} />}
            <ResourceManager<Subject>
                title="Subjects"
                data={enrichedData}
                fields={SUBJECT_FIELDS}
                idKey="subject_id"
                isLoading={isLoading}
                readOnly
                onCreate={async () => {}}
                onUpdate={async () => {}}
                onDelete={async () => {}}
            />
        </main>
    )
}

export default DirectorSubjectsPage

