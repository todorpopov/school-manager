import React, { useEffect, useState } from 'react'
import { ResourceManager } from '../../components/ResourceManager'
import type { FieldConfig } from '../../components/ResourceManager'
import { Toast } from '../../components/Toast'
import type { Teacher } from '../../types/teachers.ts'
import type { Subject } from '../../types/subjects.ts'
import { useGetTeachers } from '../../hooks/useTeachers.ts'
import { useToast } from '../../hooks/useToast'
import axiosInstance from '../../utils/axiosConfig'

const API_URL = import.meta.env.VITE_API_URL as string + '/api'

const PrincipalTeachersPage: React.FC = () => {
    const { data = [], isLoading, error } = useGetTeachers()
    const { toast, show, dismiss } = useToast()
    const [teacherSubjectsMap, setTeacherSubjectsMap] = useState<Record<number, Subject[]>>({})

    useEffect(() => {
        if (!data.length) return
        const fetchAll = async () => {
            const entries = await Promise.all(
                data.map(async (t) => {
                    try {
                        const res = await axiosInstance.get<{ data: Subject[] }>(`${API_URL}/teacher-subject/teacher/${t.teacher_id}/subjects`)
                        return [t.teacher_id, res.data.data ?? []] as [number, Subject[]]
                    } catch {
                        return [t.teacher_id, []] as [number, Subject[]]
                    }
                })
            )
            setTeacherSubjectsMap(Object.fromEntries(entries))
        }
        fetchAll()
    }, [data])

    const TEACHER_FIELDS: FieldConfig<Teacher>[] = [
        { key: 'teacher_id', label: 'ID', type: 'number', hideInForm: true, hideInTable: true },
        { key: 'school', label: 'School', type: 'text', hideInForm: true,
            renderCell: (value) => (value as { school_name: string } | undefined)?.school_name ?? '—',
        },
        { key: 'first_name', label: 'First Name', type: 'text' },
        { key: 'last_name', label: 'Last Name', type: 'text' },
        { key: 'email', label: 'Email', type: 'email' },
        { key: 'subject_ids', label: 'Subjects', type: 'text',
            renderCell: (_value, row) => {
                const subjects = teacherSubjectsMap[(row as Teacher).teacher_id] ?? []
                return subjects.length ? subjects.map(s => s.subject_name).join(', ') : '—'
            },
        },
    ]

    useEffect(() => { if (error) show(error.message, 'error') }, [error])

    return (
        <main className="max-w-5xl mx-auto px-4 py-10">
            {toast && <Toast message={toast.message} variant={toast.variant} onDismiss={dismiss} />}
            <ResourceManager<Teacher>
                title="Teachers"
                data={data}
                fields={TEACHER_FIELDS}
                idKey="teacher_id"
                isLoading={isLoading}
                readOnly
                onCreate={async () => {}}
                onUpdate={async () => {}}
                onDelete={async () => {}}
            />
        </main>
    )
}

export default PrincipalTeachersPage


