import React, { useEffect } from 'react'
import { ResourceManager } from '../../components/ResourceManager'
import type { FieldConfig } from '../../components/ResourceManager'
import { Toast } from '../../components/Toast'
import type { Class } from '../../types/classes'
import { useGetClassesBySchoolId } from '../../hooks/useClasses'
import { useDirectorSchoolId } from '../../hooks/useDirectorSchool'
import { useToast } from '../../hooks/useToast'

const PrincipalClassesPage: React.FC = () => {
    const schoolId = useDirectorSchoolId()
    const { data = [], isLoading, error } = useGetClassesBySchoolId(schoolId)
    const { toast, show, dismiss } = useToast()

    const CLASS_FIELDS: FieldConfig<Class>[] = [
        { key: 'class_id', label: 'ID', type: 'number', hideInForm: true, hideInTable: true },
        { key: 'grade_level', label: 'Grade', type: 'number' },
        { key: 'class_name', label: 'Class Name', type: 'text' },
        { key: 'school_id', label: 'School ID', type: 'number', hideInTable: true },
    ]

    useEffect(() => { if (error) show(error.message, 'error') }, [error])

    return (
        <main className="max-w-5xl mx-auto px-4 py-10">
            {toast && <Toast message={toast.message} variant={toast.variant} onDismiss={dismiss} />}
            <ResourceManager<Class>
                title="Classes"
                data={data}
                fields={CLASS_FIELDS}
                idKey="class_id"
                isLoading={isLoading}
                readOnly
                onCreate={async () => {}}
                onUpdate={async () => {}}
                onDelete={async () => {}}
            />
        </main>
    )
}

export default PrincipalClassesPage


