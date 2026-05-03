import React, { useEffect, useState } from 'react'
import { ResourceManager } from '../../components/ResourceManager'
import type { FieldConfig } from '../../components/ResourceManager'
import { Toast } from '../../components/Toast'
import type { Student } from '../../types/students.ts'
import { useGetStudents } from '../../hooks/useStudents.ts'
import { useToast } from '../../hooks/useToast'
import axiosInstance from '../../utils/axiosConfig'

const API_URL = import.meta.env.VITE_API_URL as string + '/api'

const PrincipalStudentsPage: React.FC = () => {
    const { data = [], isLoading, error } = useGetStudents()
    const { toast, show, dismiss } = useToast()
    const [studentParentsMap, setStudentParentsMap] = useState<Record<number, { first_name: string; last_name: string }[]>>({})

    useEffect(() => {
        if (!data.length) return
        const fetchAll = async () => {
            const entries = await Promise.all(
                data.map(async (s) => {
                    try {
                        const res = await axiosInstance.get<{ data: { parent_id: number; first_name: string; last_name: string }[] }>(
                            `${API_URL}/student-parent/student/${s.student_id}/parents`
                        )
                        return [s.student_id, res.data.data ?? []] as [number, typeof res.data.data]
                    } catch {
                        return [s.student_id, []] as [number, []]
                    }
                })
            )
            setStudentParentsMap(Object.fromEntries(entries))
        }
        fetchAll()
    }, [data])

    const STUDENT_FIELDS: FieldConfig<Student>[] = [
        { key: 'student_id', label: 'ID', type: 'number', hideInForm: true, hideInTable: true },
        { key: 'school', label: 'School', type: 'text', hideInForm: true,
            renderCell: (value) => (value as { school_name: string } | undefined)?.school_name ?? '—',
        },
        { key: 'first_name', label: 'First Name', type: 'text' },
        { key: 'last_name', label: 'Last Name', type: 'text' },
        { key: 'email', label: 'Email', type: 'email' },
        { key: 'class', label: 'Class', type: 'text',
            renderCell: (value) => {
                const c = value as { grade_level: number; class_name: string } | null | undefined
                return c ? `${c.grade_level}${c.class_name}` : '—'
            },
        },
        { key: 'parent_ids', label: 'Parents', type: 'text',
            renderCell: (_value, row) => {
                const parents = studentParentsMap[(row as Student).student_id] ?? []
                return parents.length ? parents.map(p => `${p.first_name} ${p.last_name}`).join(', ') : '—'
            },
        },
    ]

    useEffect(() => { if (error) show(error.message, 'error') }, [error])

    return (
        <main className="max-w-5xl mx-auto px-4 py-10">
            {toast && <Toast message={toast.message} variant={toast.variant} onDismiss={dismiss} />}
            <ResourceManager<Student>
                title="Students"
                data={data}
                fields={STUDENT_FIELDS}
                idKey="student_id"
                isLoading={isLoading}
                readOnly
                onCreate={async () => {}}
                onUpdate={async () => {}}
                onDelete={async () => {}}
            />
        </main>
    )
}

export default PrincipalStudentsPage


