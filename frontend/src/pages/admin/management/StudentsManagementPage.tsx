import React, { useEffect, useState } from 'react'
import { ResourceManager } from '../../../components/ResourceManager'
import type { FieldConfig } from '../../../components/ResourceManager'
import { Toast } from '../../../components/Toast'
import type { Student } from '../../../types/students.ts'
import { useGetStudents, useCreateStudent, useUpdateStudent, useDeleteStudent } from '../../../hooks/useStudents.ts'
import { useGetSchools } from '../../../hooks/useSchools.ts'
import { useGetClasses } from '../../../hooks/useClasses.ts'
import { useGetParents } from '../../../hooks/useParents.ts'
import { useToast } from '../../../hooks/useToast'
import axiosInstance from '../../../utils/axiosConfig'
import { validateName, validateEmail, validatePassword } from '../../../utils/validators'

const API_URL = import.meta.env.VITE_API_URL as string + '/api'

const StudentsManagementPage: React.FC = () => {
    const { data = [], isLoading, error } = useGetStudents()
    const { data: schools = [] } = useGetSchools()
    const { data: classes = [] } = useGetClasses()
    const { data: allParents = [] } = useGetParents()
    const createMutation = useCreateStudent()
    const updateMutation = useUpdateStudent()
    const deleteMutation = useDeleteStudent()

    const { toast, show, dismiss } = useToast()
    const [studentParentsMap, setStudentParentsMap] = useState<Record<number, { parent_id: number; first_name: string; last_name: string }[]>>({})
    const [parentsRefetchKey, setParentsRefetchKey] = useState(0)

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
    }, [data, parentsRefetchKey])

    const parentOptions = allParents.map((p) => ({
        label: `${p.first_name as string} ${p.last_name as string}`,
        value: String(p.parent_id),
    }))

    const syncParents = async (studentId: number, newParentIds: string[]) => {
        const current = (studentParentsMap[studentId] ?? []).map(p => String(p.parent_id))
        const toUnlink = current.filter(id => !newParentIds.includes(id))
        const toLink   = newParentIds.filter(id => !current.includes(id))
        await Promise.all([
            ...toUnlink.map(id => axiosInstance.delete(`${API_URL}/student-parent/student/${studentId}/parent/${id}`)),
            ...toLink.map(id => axiosInstance.post(`${API_URL}/student-parent/student/${studentId}/parent/${id}`)),
        ])
        setParentsRefetchKey(k => k + 1)
    }

    const schoolOptions = schools.map((s) => ({
        label: s.school_name as string,
        value: s.school_id as number,
    }))

    const gradeOptions = Array.from(new Set(classes.map((c) => c.grade_level)))
        .sort((a, b) => a - b)
        .map((g) => ({ label: String(g), value: g }))

    const classNameOptions = Array.from(new Set(classes.map((c) => c.class_name)))
        .sort()
        .map((n) => ({ label: n, value: n }))

    const STUDENT_FIELDS: FieldConfig<Student>[] = [
        { key: 'student_id', label: 'ID', type: 'number', hideInForm: true, hideInTable: true },
        { key: 'user_id', label: 'User ID', type: 'number', hideInForm: true, hideInTable: true },
        { key: 'school', label: 'School', type: 'text', hideInForm: true,
            renderCell: (value) => {
                const s = value as { school_name: string } | undefined
                return s?.school_name ?? '—'
            },
        },
        { key: 'school_id', label: 'School', type: 'select', required: true, options: schoolOptions, hideInTable: true },
        { key: 'first_name', label: 'First Name', type: 'text', required: true, placeholder: 'First name', validate: validateName },
        { key: 'last_name', label: 'Last Name', type: 'text', required: true, placeholder: 'Last name', validate: validateName },
        { key: 'email', label: 'Email', type: 'email', required: true, placeholder: 'Email', validate: validateEmail },
        { key: 'password', label: 'Password', type: 'password', required: true, placeholder: 'Password', hideInTable: true, hideInEditForm: true, validate: validatePassword },
        { key: 'class', label: 'Class', type: 'text', hideInForm: true,
            renderCell: (value) => {
                const c = value as { grade_level: number; class_name: string } | null | undefined
                return c ? `${c.grade_level}${c.class_name}` : '—'
            },
        },
        { key: 'grade_level', label: 'Grade', type: 'select', options: gradeOptions, hideInTable: true },
        { key: 'class_name', label: 'Class Letter', type: 'select', options: classNameOptions, hideInTable: true },
        { key: 'parent_ids', label: 'Parents', type: 'multiselect', options: parentOptions,
            renderCell: (_value, row) => {
                const parents = studentParentsMap[(row as Student).student_id] ?? []
                if (!parents.length) return '—'
                return parents.map(p => `${p.first_name} ${p.last_name}`).join(', ')
            },
        },
        { key: 'roles', label: 'Roles', type: 'text', hideInForm: true, hideInTable: true },
    ]

    useEffect(() => { if (error) show(error.message, 'error') }, [error])
    useEffect(() => { if (createMutation.error) show((createMutation.error as Error).message, 'error') }, [createMutation.error])
    useEffect(() => { if (updateMutation.error) show((updateMutation.error as Error).message, 'error') }, [updateMutation.error])
    useEffect(() => { if (deleteMutation.error) show((deleteMutation.error as Error).message, 'error') }, [deleteMutation.error])

    const isMutating = createMutation.isPending || updateMutation.isPending || deleteMutation.isPending

    const enrichedData = data.map((s) => ({
        ...s,
        school_id: (s.school as { school_id: number } | undefined)?.school_id,
        grade_level: (s.class as { grade_level: number } | null | undefined)?.grade_level ?? '',
        class_name: (s.class as { class_name: string } | null | undefined)?.class_name ?? '',
        parent_ids: (studentParentsMap[s.student_id] ?? []).map(p => String(p.parent_id)),
    }))

    const resolveClassId = (gradeLevel: unknown, className: unknown): number | null => {
        if (!gradeLevel || !className) return null
        const match = classes.find(
            (c) => c.grade_level === Number(gradeLevel) && c.class_name === String(className)
        )
        return match?.class_id ?? null
    }

    return (
        <main className="max-w-5xl mx-auto px-4 py-10">
            {toast && <Toast message={toast.message} variant={toast.variant} onDismiss={dismiss} />}
            <ResourceManager<Student>
                title="Students"
                data={enrichedData}
                fields={STUDENT_FIELDS}
                idKey="student_id"
                isLoading={isLoading || isMutating}
                onCreate={async (values) => {
                    const student = await createMutation.mutateAsync({
                        school_id: Number(values.school_id),
                        first_name: String(values.first_name ?? ''),
                        last_name: String(values.last_name ?? ''),
                        email: String(values.email ?? ''),
                        password: String(values.password ?? ''),
                        class_id: resolveClassId(values.grade_level, values.class_name),
                    })
                    const parentIds = (values.parent_ids as string[]) ?? []
                    if (student && parentIds.length) {
                        await syncParents(student.student_id, parentIds)
                    }
                }}
                onUpdate={async (id, values) => {
                    await updateMutation.mutateAsync({
                        id: id as number,
                        payload: {
                            school_id: Number(values.school_id),
                            first_name: String(values.first_name ?? ''),
                            last_name: String(values.last_name ?? ''),
                            email: String(values.email ?? ''),
                            class_id: resolveClassId(values.grade_level, values.class_name),
                        },
                    })
                    await syncParents(id as number, (values.parent_ids as string[]) ?? [])
                }}
                onDelete={async (id) => { await deleteMutation.mutateAsync(id as number) }}
            />
        </main>
    )
}

export default StudentsManagementPage
