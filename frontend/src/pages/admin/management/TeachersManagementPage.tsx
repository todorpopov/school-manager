import React, { useEffect, useState } from 'react'
import { ResourceManager } from '../../../components/ResourceManager'
import type { FieldConfig } from '../../../components/ResourceManager'
import { Toast } from '../../../components/Toast'
import type { Teacher } from '../../../types/teachers.ts'
import type { Subject } from '../../../types/subjects.ts'
import { useGetTeachers, useCreateTeacher, useUpdateTeacher, useDeleteTeacher } from '../../../hooks/useTeachers.ts'
import { useGetSchools } from '../../../hooks/useSchools.ts'
import { useGetSubjects } from '../../../hooks/useSubjects.ts'
import { useToast } from '../../../hooks/useToast'
import axiosInstance from '../../../utils/axiosConfig'
import { validateName, validateEmail, validatePassword } from '../../../utils/validators'

const API_URL = import.meta.env.VITE_API_URL as string + '/api'

const TeachersManagementPage: React.FC = () => {
    const { data = [], isLoading, error } = useGetTeachers()
    const { data: schools = [] } = useGetSchools()
    const { data: allSubjects = [] } = useGetSubjects()
    const createMutation = useCreateTeacher()
    const updateMutation = useUpdateTeacher()
    const deleteMutation = useDeleteTeacher()

    const { toast, show, dismiss } = useToast()
    const [teacherSubjectsMap, setTeacherSubjectsMap] = useState<Record<number, Subject[]>>({})
    const [subjectRefetchKey, setSubjectRefetchKey] = useState(0)

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
    }, [data, subjectRefetchKey])

    const schoolOptions = schools.map((s) => ({
        label: s.school_name as string,
        value: s.school_id as number,
    }))

    const subjectOptions = allSubjects.map((s) => ({
        label: s.subject_name,
        value: String(s.subject_id),
    }))

    const TEACHER_FIELDS: FieldConfig<Teacher>[] = [
        { key: 'teacher_id', label: 'ID', type: 'number', hideInForm: true, hideInTable: true },
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
        { key: 'password', label: 'Password', type: 'password', required: true, placeholder: 'Password', hideInTable: true, validate: validatePassword },
        { key: 'subject_ids', label: 'Subjects', type: 'multiselect', options: subjectOptions,
            renderCell: (_value, row) => {
                const subjects = teacherSubjectsMap[(row as Teacher).teacher_id] ?? []
                if (!subjects.length) return '—'
                return subjects.map(s => s.subject_name).join(', ')
            },
        },
        { key: 'roles', label: 'Roles', type: 'text', hideInForm: true, hideInTable: true },
    ]

    useEffect(() => { if (error) show(error.message, 'error') }, [error])
    useEffect(() => { if (createMutation.error) show((createMutation.error as Error).message, 'error') }, [createMutation.error])
    useEffect(() => { if (updateMutation.error) show((updateMutation.error as Error).message, 'error') }, [updateMutation.error])
    useEffect(() => { if (deleteMutation.error) show((deleteMutation.error as Error).message, 'error') }, [deleteMutation.error])

    const syncSubjects = async (teacherId: number, newSubjectIds: string[]) => {
        const current = (teacherSubjectsMap[teacherId] ?? []).map(s => String(s.subject_id))
        const toUnlink = current.filter(id => !newSubjectIds.includes(id))
        const toLink   = newSubjectIds.filter(id => !current.includes(id))
        await Promise.all([
            ...toUnlink.map(id => axiosInstance.delete(`${API_URL}/teacher-subject/teacher/${teacherId}/subject/${id}`)),
            ...toLink.map(id => axiosInstance.post(`${API_URL}/teacher-subject/teacher/${teacherId}/subject/${id}`)),
        ])
        setSubjectRefetchKey(k => k + 1)
    }

    const isMutating = createMutation.isPending || updateMutation.isPending || deleteMutation.isPending

    const enrichedData = data.map((t) => ({
        ...t,
        school_id: (t.school as { school_id: number } | undefined)?.school_id,
        subject_ids: (teacherSubjectsMap[t.teacher_id] ?? []).map(s => String(s.subject_id)),
    }))

    return (
        <main className="max-w-5xl mx-auto px-4 py-10">
            {toast && <Toast message={toast.message} variant={toast.variant} onDismiss={dismiss} />}
            <ResourceManager<Teacher>
                title="Teachers"
                data={enrichedData}
                fields={TEACHER_FIELDS}
                idKey="teacher_id"
                isLoading={isLoading || isMutating}
                onCreate={async (values) => {
                    const teacher = await createMutation.mutateAsync({
                        school_id: Number(values.school_id),
                        first_name: String(values.first_name ?? ''),
                        last_name: String(values.last_name ?? ''),
                        email: String(values.email ?? ''),
                        password: String(values.password ?? ''),
                    })
                    const subjectIds = (values.subject_ids as string[]) ?? []
                    if (teacher && subjectIds.length) {
                        await syncSubjects(teacher.teacher_id, subjectIds)
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
                        },
                    })
                    await syncSubjects(id as number, (values.subject_ids as string[]) ?? [])
                }}
                onDelete={async (id) => { await deleteMutation.mutateAsync(id as number) }}
            />
        </main>
    )
}

export default TeachersManagementPage
