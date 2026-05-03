import React, { useEffect, useState } from 'react'
import { ResourceManager } from '../../components/ResourceManager'
import type { FieldConfig } from '../../components/ResourceManager'
import { Toast } from '../../components/Toast'
import type { Parent } from '../../types/parents.ts'
import { useGetParents } from '../../hooks/useParents.ts'
import { useToast } from '../../hooks/useToast'
import axiosInstance from '../../utils/axiosConfig'

const API_URL = import.meta.env.VITE_API_URL as string + '/api'

const PrincipalParentsPage: React.FC = () => {
    const { data = [], isLoading, error } = useGetParents()
    const { toast, show, dismiss } = useToast()
    const [parentStudentsMap, setParentStudentsMap] = useState<Record<number, { first_name: string; last_name: string }[]>>({})

    useEffect(() => {
        if (!data.length) return
        const fetchAll = async () => {
            const entries = await Promise.all(
                data.map(async (p) => {
                    try {
                        const res = await axiosInstance.get<{ data: { student_id: number; first_name: string; last_name: string }[] }>(
                            `${API_URL}/student-parent/parent/${p.parent_id}/students`
                        )
                        return [p.parent_id, res.data.data ?? []] as [number, typeof res.data.data]
                    } catch {
                        return [p.parent_id, []] as [number, []]
                    }
                })
            )
            setParentStudentsMap(Object.fromEntries(entries))
        }
        fetchAll()
    }, [data])

    const PARENT_FIELDS: FieldConfig<Parent>[] = [
        { key: 'parent_id', label: 'ID', type: 'number', hideInForm: true, hideInTable: true },
        { key: 'first_name', label: 'First Name', type: 'text' },
        { key: 'last_name', label: 'Last Name', type: 'text' },
        { key: 'email', label: 'Email', type: 'email' },
        { key: 'student_ids', label: 'Children', type: 'text',
            renderCell: (_value, row) => {
                const students = parentStudentsMap[(row as Parent).parent_id] ?? []
                return students.length ? students.map(s => `${s.first_name} ${s.last_name}`).join(', ') : '—'
            },
        },
    ]

    useEffect(() => { if (error) show(error.message, 'error') }, [error])

    return (
        <main className="max-w-5xl mx-auto px-4 py-10">
            {toast && <Toast message={toast.message} variant={toast.variant} onDismiss={dismiss} />}
            <ResourceManager<Parent>
                title="Parents"
                data={data}
                fields={PARENT_FIELDS}
                idKey="parent_id"
                isLoading={isLoading}
                readOnly
                onCreate={async () => {}}
                onUpdate={async () => {}}
                onDelete={async () => {}}
            />
        </main>
    )
}

export default PrincipalParentsPage


