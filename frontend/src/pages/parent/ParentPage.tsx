import React, { useEffect, useState } from 'react'
import { Toast } from '../../components/Toast'
import { useToast } from '../../hooks/useToast'
import { useGetGradesForStudent, useGetAbsencesForStudent } from '../../hooks/useGradesAbsences'
import type { Grade, Absence } from '../../types/gradesAbsences'
import axiosInstance from '../../utils/axiosConfig'

const API_URL = import.meta.env.VITE_API_URL as string + '/api'

const MOCK_PARENT_ID = 1

interface Child {
    student_id: number
    first_name: string
    last_name: string
}

type Tab = 'grades' | 'absences'

const ChildGrades: React.FC<{ studentId: number }> = ({ studentId }) => {
    const { data: grades = [], isLoading } = useGetGradesForStudent(studentId)

    if (isLoading) return <div className="flex justify-center py-8"><div className="w-6 h-6 border-4 border-slate-200 border-t-indigo-500 rounded-full animate-spin" /></div>
    if (!grades.length) return <p className="text-center text-slate-400 py-8 text-sm">No grades recorded.</p>

    return (
        <table className="w-full border-collapse text-sm">
            <thead className="bg-slate-100 dark:bg-slate-900">
                <tr>
                    {['Subject', 'Term', 'Grade', 'Date'].map(h => (
                        <th key={h} className="px-4 py-3 text-left text-xs font-semibold text-slate-500 dark:text-slate-400 uppercase tracking-wide">{h}</th>
                    ))}
                </tr>
            </thead>
            <tbody>
                {grades.map((g: Grade) => (
                    <tr key={g.grade_id} className="border-t border-slate-100 dark:border-slate-700 hover:bg-slate-50 dark:hover:bg-slate-900/50">
                        <td className="px-4 py-3 text-slate-700 dark:text-slate-300">{g.curriculum.subject.subject_name}</td>
                        <td className="px-4 py-3 text-slate-700 dark:text-slate-300">{g.curriculum.term.name}</td>
                        <td className="px-4 py-3">
                            <span className={`font-semibold ${g.grade_value >= 5 ? 'text-green-600' : g.grade_value >= 3.5 ? 'text-yellow-600' : 'text-red-500'}`}>
                                {g.grade_value.toFixed(2)}
                            </span>
                        </td>
                        <td className="px-4 py-3 text-slate-500 dark:text-slate-400">{new Date(g.grade_date).toLocaleDateString()}</td>
                    </tr>
                ))}
            </tbody>
        </table>
    )
}

const ChildAbsences: React.FC<{ studentId: number }> = ({ studentId }) => {
    const { data: absences = [], isLoading } = useGetAbsencesForStudent(studentId)

    if (isLoading) return <div className="flex justify-center py-8"><div className="w-6 h-6 border-4 border-slate-200 border-t-indigo-500 rounded-full animate-spin" /></div>
    if (!absences.length) return <p className="text-center text-slate-400 py-8 text-sm">No absences recorded.</p>

    const excused = absences.filter(a => a.is_excused).length
    const unexcused = absences.length - excused

    return (
        <div className="flex flex-col gap-4">
            <div className="flex gap-4 px-4">
                <div className="flex items-center gap-2 text-sm text-slate-600 dark:text-slate-300">
                    <span className="w-2.5 h-2.5 rounded-full bg-yellow-400 inline-block" />
                    Excused: <strong>{excused}</strong>
                </div>
                <div className="flex items-center gap-2 text-sm text-slate-600 dark:text-slate-300">
                    <span className="w-2.5 h-2.5 rounded-full bg-red-500 inline-block" />
                    Unexcused: <strong>{unexcused}</strong>
                </div>
            </div>
            <table className="w-full border-collapse text-sm">
                <thead className="bg-slate-100 dark:bg-slate-900">
                    <tr>
                        {['Subject', 'Term', 'Date', 'Status'].map(h => (
                            <th key={h} className="px-4 py-3 text-left text-xs font-semibold text-slate-500 dark:text-slate-400 uppercase tracking-wide">{h}</th>
                        ))}
                    </tr>
                </thead>
                <tbody>
                    {absences.map((a: Absence) => (
                        <tr key={a.absence_id} className="border-t border-slate-100 dark:border-slate-700 hover:bg-slate-50 dark:hover:bg-slate-900/50">
                            <td className="px-4 py-3 text-slate-700 dark:text-slate-300">{a.curriculum.subject.subject_name}</td>
                            <td className="px-4 py-3 text-slate-700 dark:text-slate-300">{a.curriculum.term.name}</td>
                            <td className="px-4 py-3 text-slate-500 dark:text-slate-400">{new Date(a.absence_date).toLocaleDateString()}</td>
                            <td className="px-4 py-3">
                                <span className={`text-xs font-medium px-2 py-0.5 rounded-full ${a.is_excused ? 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900/30 dark:text-yellow-400' : 'bg-red-100 text-red-600 dark:bg-red-900/30 dark:text-red-400'}`}>
                                    {a.is_excused ? 'Excused' : 'Unexcused'}
                                </span>
                            </td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    )
}

const ParentPage: React.FC = () => {
    const [children, setChildren] = useState<Child[]>([])
    const [selectedChild, setSelectedChild] = useState<Child | null>(null)
    const [tab, setTab] = useState<Tab>('grades')
    const [loadingChildren, setLoadingChildren] = useState(true)
    const { toast, show, dismiss } = useToast()

    useEffect(() => {
        const fetch = async () => {
            try {
                const res = await axiosInstance.get<{ data: Child[] }>(
                    `${API_URL}/student-parent/parent/${MOCK_PARENT_ID}/students`
                )
                const kids = res.data.data ?? []
                setChildren(kids)
                if (kids.length) setSelectedChild(kids[0])
            } catch {
                show('Failed to load children', 'error')
            } finally {
                setLoadingChildren(false)
            }
        }
        fetch()
    }, [])

    return (
        <main className="max-w-4xl mx-auto px-4 py-10 flex flex-col gap-6">
            {toast && <Toast message={toast.message} variant={toast.variant} onDismiss={dismiss} />}

            <h1 className="text-2xl font-semibold text-slate-800 dark:text-slate-100">My Children</h1>
            {loadingChildren && (
                <div className="flex justify-center py-10">
                    <div className="w-7 h-7 border-4 border-slate-200 border-t-indigo-500 rounded-full animate-spin" />
                </div>
            )}
            {!loadingChildren && !children.length && (
                <p className="text-center text-slate-400 py-10 text-sm">No children linked to your account.</p>
            )}
            {!loadingChildren && children.length > 0 && (
                <>
                    {children.length > 1 && (
                        <div className="flex gap-2 flex-wrap">
                            {children.map(child => (
                                <button
                                    key={child.student_id}
                                    onClick={() => { setSelectedChild(child); setTab('grades') }}
                                    className={`px-4 py-2 rounded-full text-sm font-medium border transition-colors cursor-pointer ${
                                        selectedChild?.student_id === child.student_id
                                            ? 'bg-indigo-500 text-white border-indigo-500'
                                            : 'bg-white dark:bg-slate-800 text-slate-600 dark:text-slate-300 border-slate-300 dark:border-slate-600 hover:border-indigo-400'
                                    }`}
                                >
                                    {child.first_name} {child.last_name}
                                </button>
                            ))}
                        </div>
                    )}
                    {selectedChild && (
                        <div className="bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl overflow-hidden">
                            <div className="px-5 py-4 bg-slate-50 dark:bg-slate-900 border-b border-slate-200 dark:border-slate-700">
                                <h2 className="text-lg font-semibold text-slate-800 dark:text-slate-100">
                                    {selectedChild.first_name} {selectedChild.last_name}
                                </h2>
                            </div>
                            <div className="flex border-b border-slate-200 dark:border-slate-700">
                                {(['grades', 'absences'] as Tab[]).map(t => (
                                    <button
                                        key={t}
                                        onClick={() => setTab(t)}
                                        className={`px-6 py-3 text-sm font-medium capitalize border-b-2 transition-colors cursor-pointer bg-transparent ${
                                            tab === t
                                                ? 'border-indigo-500 text-indigo-600 dark:text-indigo-400'
                                                : 'border-transparent text-slate-500 hover:text-slate-700 dark:hover:text-slate-300'
                                        }`}
                                    >
                                        {t}
                                    </button>
                                ))}
                            </div>
                            <div className="overflow-x-auto">
                                {tab === 'grades'   && <ChildGrades   studentId={selectedChild.student_id} />}
                                {tab === 'absences' && <ChildAbsences studentId={selectedChild.student_id} />}
                            </div>
                        </div>
                    )}
                </>
            )}
        </main>
    )
}

export default ParentPage

