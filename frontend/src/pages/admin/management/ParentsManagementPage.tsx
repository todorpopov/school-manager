import React, { useEffect, useState } from 'react'
import { ResourceManager } from '../../../components/ResourceManager'
import type { FieldConfig } from '../../../components/ResourceManager'
import { Toast } from '../../../components/Toast'
import type { Parent } from '../../../types/parents.ts'
import { useGetParents, useCreateParent, useUpdateParent, useDeleteParent, linkParentToStudent, unlinkParentFromStudent } from '../../../hooks/useParents.ts'
import { useGetStudents } from '../../../hooks/useStudents.ts'
import { useToast } from '../../../hooks/useToast'
import axiosInstance from '../../../utils/axiosConfig'
import { validateName, validateEmail, validatePassword } from '../../../utils/validators'

const API_URL = import.meta.env.VITE_API_URL as string + '/api'

const ParentsManagementPage: React.FC = () => {
    const { data = [], isLoading, error } = useGetParents()
    const { data: allStudents = [] } = useGetStudents()
    const createMutation = useCreateParent()
    const updateMutation = useUpdateParent()
    const deleteMutation = useDeleteParent()

    const { toast, show, dismiss } = useToast()
    const [parentStudentsMap, setParentStudentsMap] = useState<Record<number, { student_id: number; first_name: string; last_name: string }[]>>({})
    const [studentRefetchKey, setStudentRefetchKey] = useState(0)

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
    }, [data, studentRefetchKey])

    const studentOptions = allStudents.map((s) => ({
        label: `${s.first_name as string} ${s.last_name as string}`,
        value: String(s.student_id),
    }))

    const syncStudents = async (parentId: number, newStudentIds: string[]) => {
        const current = (parentStudentsMap[parentId] ?? []).map(s => String(s.student_id))
        const toUnlink = current.filter(id => !newStudentIds.includes(id))
        const toLink   = newStudentIds.filter(id => !current.includes(id))
        await Promise.all([
            ...toUnlink.map(id => unlinkParentFromStudent(Number(id), parentId)),
            ...toLink.map(id => linkParentToStudent(Number(id), parentId)),
        ])
        setStudentRefetchKey(k => k + 1)
    }

    const PARENT_FIELDS: FieldConfig<Parent>[] = [
        { key: 'parent_id', label: 'ID', type: 'number', hideInForm: true, hideInTable: true },
        { key: 'user_id', label: 'User ID', type: 'number', hideInForm: true, hideInTable: true },
        { key: 'first_name', label: 'First Name', type: 'text', required: true, placeholder: 'First name', validate: validateName },
        { key: 'last_name', label: 'Last Name', type: 'text', required: true, placeholder: 'Last name', validate: validateName },
        { key: 'email', label: 'Email', type: 'email', required: true, placeholder: 'Email', validate: validateEmail },
        { key: 'password', label: 'Password', type: 'password', required: true, placeholder: 'Password', hideInTable: true, hideInEditForm: true, validate: validatePassword },
        { key: 'student_ids', label: 'Children', type: 'multiselect', options: studentOptions,
            renderCell: (_value, row) => {
                const students = parentStudentsMap[(row as Parent).parent_id] ?? []
                if (!students.length) return '—'
                return students.map(s => `${s.first_name} ${s.last_name}`).join(', ')
            },
        },
        { key: 'roles', label: 'Roles', type: 'text', hideInForm: true, hideInTable: true },
    ]

    useEffect(() => { if (error) show(error.message, 'error') }, [error])
    useEffect(() => { if (createMutation.error) show((createMutation.error as Error).message, 'error') }, [createMutation.error])
    useEffect(() => { if (updateMutation.error) show((updateMutation.error as Error).message, 'error') }, [updateMutation.error])
    useEffect(() => { if (deleteMutation.error) show((deleteMutation.error as Error).message, 'error') }, [deleteMutation.error])

    const isMutating = createMutation.isPending || updateMutation.isPending || deleteMutation.isPending

    const enrichedData = data.map((p) => ({
        ...p,
        student_ids: (parentStudentsMap[p.parent_id] ?? []).map(s => String(s.student_id)),
    }))

    return (
        <main className="max-w-5xl mx-auto px-4 py-10">
            {toast && <Toast message={toast.message} variant={toast.variant} onDismiss={dismiss} />}
            <ResourceManager<Parent>
                title="Parents"
                data={enrichedData}
                fields={PARENT_FIELDS}
                idKey="parent_id"
                isLoading={isLoading || isMutating}
                onCreate={async (values) => {
                    const parent = await createMutation.mutateAsync({
                        first_name: String(values.first_name ?? ''),
                        last_name: String(values.last_name ?? ''),
                        email: String(values.email ?? ''),
                        password: String(values.password ?? ''),
                    })
                    const studentIds = (values.student_ids as string[]) ?? []
                    if (parent && studentIds.length) {
                        await syncStudents(parent.parent_id, studentIds)
                    }
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
                    await syncStudents(id as number, (values.student_ids as string[]) ?? [])
                }}
                onDelete={async (id) => { await deleteMutation.mutateAsync(id as number) }}
            />
        </main>
    )
}

export default ParentsManagementPage
