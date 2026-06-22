import React, { useState } from 'react'
import { useMutation, useQuery } from '@tanstack/react-query'
import { useAuth } from '../../hooks/useAuth'
import axiosInstance from '../../utils/axiosConfig'
import { parseApiError } from '../../utils/parseApiError'
import type { Principals } from '../../types/principals'

const API_URL = import.meta.env.VITE_API_URL as string + '/api'

interface GradeDistribution {
    number_of_sixes: number
    number_of_fives: number
    number_of_fours: number
    number_of_threes: number
    number_of_twos: number
}
interface GradeSummary {
    number_of_grades: number
    grade_distribution: GradeDistribution
}
interface Subject { subject_id: number; subject_name: string }
interface Teacher { teacher_id: number; first_name: string; last_name: string; email: string }
interface SubjectGradeSummary { subject: Subject; absence_summary: GradeSummary }
interface TeacherSubjectsGrades { teacher: Teacher; subjects: SubjectGradeSummary[] }
interface SchoolGrades { school: { school_id: number; school_name: string; school_address: string }; teachers: TeacherSubjectsGrades[] }

type ReportData = SchoolGrades[] | null

type ReportBreakdown = 'all' | 'per_teacher' | 'per_subject'

interface ReportConfig {
    id: ReportBreakdown
    label: string
}

const REPORTS: ReportConfig[] = [
    { id: 'all',         label: 'Grades' },
    { id: 'per_teacher', label: 'Grades per teacher' },
    { id: 'per_subject', label: 'Grades per subject' },
]

const thCls = 'px-4 py-2 text-left text-xs font-semibold text-slate-500 dark:text-slate-400 uppercase tracking-wide'
const tdCls = 'px-4 py-2 text-sm text-slate-700 dark:text-slate-300'
const numCls = 'px-4 py-2 text-sm text-slate-700 dark:text-slate-300 text-left tabular-nums'

function aggregate(teachers: TeacherSubjectsGrades[]): { total: number; sixes: number; fives: number; fours: number; threes: number; twos: number } {
    return (teachers ?? []).reduce((acc, t) => {
        (t.subjects ?? []).forEach(s => {
            const gd = s.absence_summary.grade_distribution
            acc.total += s.absence_summary.number_of_grades
            acc.sixes += gd.number_of_sixes
            acc.fives += gd.number_of_fives
            acc.fours += gd.number_of_fours
            acc.threes += gd.number_of_threes
            acc.twos += gd.number_of_twos
        })
        return acc
    }, { total: 0, sixes: 0, fives: 0, fours: 0, threes: 0, twos: 0 })
}

const GradeRow: React.FC<{ label: string; sixes: number; fives: number; fours: number; threes: number; twos: number; total: number }> = ({ label, sixes, fives, fours, threes, twos, total }) => (
    <tr className="border-t border-slate-100 dark:border-slate-700">
        <td className={tdCls}>{label}</td>
        <td className={numCls}><span className="font-semibold">{total}</span></td>
        <td className={numCls}><span className="text-green-600">{sixes}</span></td>
        <td className={numCls}><span className="text-emerald-600">{fives}</span></td>
        <td className={numCls}><span className="text-yellow-600">{fours}</span></td>
        <td className={numCls}><span className="text-orange-500">{threes}</span></td>
        <td className={numCls}><span className="text-red-500">{twos}</span></td>
    </tr>
)

const ResultTable: React.FC<{ data: SchoolGrades[]; breakdown: ReportBreakdown }> = ({ data, breakdown }) => {
    if (!data.length) return <p className="text-center text-slate-400 py-8 text-sm">No grade data found.</p>

    const school = data[0]
    const headers = breakdown === 'per_teacher' ? ['Teacher', 'Total', '6', '5', '4', '3', '2']
        : breakdown === 'per_subject' ? ['Subject', 'Total', '6', '5', '4', '3', '2']
        : ['', 'Total', '6', '5', '4', '3', '2']

    return (
        <div className="bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl overflow-hidden">
            <div className="px-5 py-3 bg-slate-50 dark:bg-slate-900 border-b border-slate-200 dark:border-slate-700">
                <h3 className="font-semibold text-slate-800 dark:text-slate-100">{school.school.school_name}</h3>
            </div>
            <table className="w-full border-collapse">
                <thead className="bg-slate-50 dark:bg-slate-900/40">
                    <tr>{headers.map(h => <th key={h} className={thCls}>{h}</th>)}</tr>
                </thead>
                <tbody>
                    {breakdown === 'all' && (() => {
                        const t = aggregate(school.teachers)
                        return <GradeRow label="School total" {...t} />
                    })()}
                    {breakdown === 'per_teacher' && school.teachers.map(t => {
                        const totals = (t.subjects ?? []).reduce((acc, s) => {
                            const gd = s.absence_summary.grade_distribution
                            acc.total += s.absence_summary.number_of_grades
                            acc.sixes += gd.number_of_sixes
                            acc.fives += gd.number_of_fives
                            acc.fours += gd.number_of_fours
                            acc.threes += gd.number_of_threes
                            acc.twos += gd.number_of_twos
                            return acc
                        }, { total: 0, sixes: 0, fives: 0, fours: 0, threes: 0, twos: 0 })
                        return <GradeRow key={t.teacher.teacher_id} label={`${t.teacher.first_name} ${t.teacher.last_name}`} {...totals} />
                    })}
                    {breakdown === 'per_subject' && (() => {
                        const subjectMap = new Map<number, { subject: Subject } & { total: number; sixes: number; fives: number; fours: number; threes: number; twos: number }>()
                        school.teachers.forEach(t => {
                            (t.subjects ?? []).forEach(s => {
                                const gd = s.absence_summary.grade_distribution
                                const existing = subjectMap.get(s.subject.subject_id)
                                if (existing) {
                                    existing.total += s.absence_summary.number_of_grades
                                    existing.sixes += gd.number_of_sixes
                                    existing.fives += gd.number_of_fives
                                    existing.fours += gd.number_of_fours
                                    existing.threes += gd.number_of_threes
                                    existing.twos += gd.number_of_twos
                                } else {
                                    subjectMap.set(s.subject.subject_id, {
                                        subject: s.subject,
                                        total: s.absence_summary.number_of_grades,
                                        sixes: gd.number_of_sixes,
                                        fives: gd.number_of_fives,
                                        fours: gd.number_of_fours,
                                        threes: gd.number_of_threes,
                                        twos: gd.number_of_twos,
                                    })
                                }
                            })
                        })
                        return Array.from(subjectMap.values()).map(r => (
                            <GradeRow key={r.subject.subject_id} label={r.subject.subject_name} total={r.total} sixes={r.sixes} fives={r.fives} fours={r.fours} threes={r.threes} twos={r.twos} />
                        ))
                    })()}
                </tbody>
            </table>
        </div>
    )
}

export default function PrincipalStatisticsPage() {
    const { user } = useAuth()
    const [selectedBreakdown, setSelectedBreakdown] = useState<ReportBreakdown>('all')
    const [reportData, setReportData] = useState<ReportData>(null)
    const [error, setError] = useState<string | null>(null)

    const { data: principal, isLoading: principalLoading } = useQuery<Principals>({
        queryKey: ['principal-by-user', user?.userId],
        queryFn: async () => {
            const { data } = await axiosInstance.get<{ data: Principals }>(`${API_URL}/director/user/${user!.userId}`)
            return data.data
        },
        enabled: !!user?.userId,
    })

    const mutation = useMutation({
        mutationFn: async (_breakdown: ReportBreakdown) => {
            const schoolId = principal?.school?.school_id
            if (!schoolId) throw new Error('School not found for this principal')
            const { data } = await axiosInstance.post<{ data: ReportData; error: boolean; message: string }>(
                `${API_URL}/reporting`,
                { report_type: 'grades', school_ids: [schoolId] }
            )
            if (data.error) throw new Error(data.message)
            return data.data
        },
        onSuccess: (data) => { setReportData(data); setError(null) },
        onError: (err) => { setError(parseApiError(err)); setReportData(null) },
    })

    const handleGenerate = () => mutation.mutate(selectedBreakdown)

    const inputCls = 'px-3 py-2 text-sm border border-slate-300 dark:border-slate-600 rounded-md bg-white dark:bg-slate-900 text-slate-800 dark:text-slate-100 focus:outline-none focus:ring-2 focus:ring-indigo-400'
    const btnPrimary = 'px-4 py-2 text-sm font-medium rounded-md bg-indigo-500 hover:bg-indigo-600 text-white border-none cursor-pointer transition-colors disabled:opacity-50'

    return (
        <main className="max-w-5xl mx-auto px-4 py-10 flex flex-col gap-6">
            <h1 className="text-2xl font-semibold text-slate-800 dark:text-slate-100">Statistics</h1>

            {principal && (
                <p className="text-sm text-slate-500 dark:text-slate-400">
                    School: <span className="font-medium text-slate-700 dark:text-slate-200">{principal.school.school_name}</span>
                </p>
            )}

            <div className="bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl p-5 flex flex-col gap-4">
                <h2 className="text-sm font-semibold text-slate-600 dark:text-slate-300 uppercase tracking-wide">Generate Report</h2>

                <div className="flex flex-col gap-1 max-w-sm">
                    <label className="text-xs font-medium text-slate-500 dark:text-slate-400">Report type</label>
                    <select
                        className={inputCls}
                        value={selectedBreakdown}
                        onChange={e => {
                            setSelectedBreakdown(e.target.value as ReportBreakdown)
                            setReportData(null)
                            setError(null)
                        }}
                    >
                        {REPORTS.map(r => (
                            <option key={r.id} value={r.id}>{r.label}</option>
                        ))}
                    </select>
                </div>

                <div className="flex justify-end">
                    <button
                        className={btnPrimary}
                        disabled={mutation.isPending || principalLoading || !principal}
                        onClick={handleGenerate}
                    >
                        {mutation.isPending ? 'Generating…' : 'Generate Report'}
                    </button>
                </div>
            </div>

            {error && (
                <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-700 rounded-xl px-5 py-4 text-sm text-red-700 dark:text-red-300">
                    {error}
                </div>
            )}

            {reportData !== null && (
                <div className="flex flex-col gap-3">
                    <h2 className="text-base font-semibold text-slate-700 dark:text-slate-200">
                        {REPORTS.find(r => r.id === selectedBreakdown)?.label}
                    </h2>
                    <ResultTable data={reportData} breakdown={selectedBreakdown} />
                </div>
            )}
        </main>
    )
}


