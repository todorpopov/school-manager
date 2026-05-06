import React, { useState } from 'react'
import { useMutation } from '@tanstack/react-query'
import { useGetSchools } from '../../hooks/useSchools'
import axiosInstance from '../../utils/axiosConfig'
import { parseApiError } from '../../utils/parseApiError'
import type { School } from '../../types/schools'

const API_URL = import.meta.env.VITE_API_URL as string + '/api'

interface AbsenceSummary {
    number_of_absences: number
    number_of_excuses: number
    number_of_non_excuses: number
}
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

interface SubjectAbsenceSummary { subject: Subject; absence_summary: AbsenceSummary }
interface SubjectGradeSummary   { subject: Subject; absence_summary: GradeSummary }

interface TeacherSubjectsAbsences { teacher: Teacher; subjects: SubjectAbsenceSummary[] }
interface TeacherSubjectsGrades   { teacher: Teacher; subjects: SubjectGradeSummary[] }

interface SchoolAbsences { school: School; teachers: TeacherSubjectsAbsences[] }
interface SchoolGrades   { school: School; teachers: TeacherSubjectsGrades[] }

type ReportData = SchoolAbsences[] | SchoolGrades[] | null

type ReportScope = 'all' | 'school'
type ReportType  = 'grades' | 'absences'
type ReportBreakdown = 'per_subject' | 'all'

interface ReportConfig {
    id: string
    label: string
    type: ReportType
    scope: ReportScope
    breakdown: ReportBreakdown
}

const REPORTS: ReportConfig[] = [
    { id: 'grades-subject-all',    label: 'Grades per subject — all schools',       type: 'grades',   scope: 'all',    breakdown: 'per_subject' },
    { id: 'absences-subject-all',  label: 'Absences per subject — all schools',     type: 'absences', scope: 'all',    breakdown: 'per_subject' },
    { id: 'grades-all',            label: 'Grades — all schools',                   type: 'grades',   scope: 'all',    breakdown: 'all' },
    { id: 'absences-all',          label: 'Absences — all schools',                 type: 'absences', scope: 'all',    breakdown: 'all' },
    { id: 'grades-subject-school', label: 'Grades per subject per school',   type: 'grades',   scope: 'school', breakdown: 'per_subject' },
    { id: 'absences-subject-school',label: 'Absences per subject per school',type: 'absences', scope: 'school', breakdown: 'per_subject' },
    { id: 'grades-school',         label: 'Grades per school',               type: 'grades',   scope: 'school', breakdown: 'all' },
    { id: 'absences-school',       label: 'Absences per school',             type: 'absences', scope: 'school', breakdown: 'all' },
]

const thCls = 'px-4 py-2 text-left text-xs font-semibold text-slate-500 dark:text-slate-400 uppercase tracking-wide'
const tdCls = 'px-4 py-2 text-sm text-slate-700 dark:text-slate-300'
const numCls = 'px-4 py-2 text-sm text-slate-700 dark:text-slate-300 text-left tabular-nums'

const GradesTable: React.FC<{ data: SchoolGrades[]; perSubject: boolean }> = ({ data, perSubject }) => {
    if (!data.length) return <p className="text-center text-slate-400 py-8 text-sm">No grade data found.</p>

    if (perSubject) {
        return (
            <div className="flex flex-col gap-6">
                {data.map(school => (
                    <div key={school.school.school_id} className="bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl overflow-hidden">
                        <div className="px-5 py-3 bg-slate-50 dark:bg-slate-900 border-b border-slate-200 dark:border-slate-700">
                            <h3 className="font-semibold text-slate-800 dark:text-slate-100">{school.school.school_name}</h3>
                            <p className="text-xs text-slate-400">{school.school.school_address}</p>
                        </div>
                        {school.teachers.map(t => (
                            <div key={t.teacher.teacher_id} className="border-b border-slate-100 dark:border-slate-700 last:border-0">
                                <div className="px-5 py-2 bg-indigo-50/50 dark:bg-indigo-900/10">
                                    <p className="text-sm font-medium text-indigo-700 dark:text-indigo-300">{t.teacher.first_name} {t.teacher.last_name}</p>
                                </div>
                                <table className="w-full border-collapse">
                                    <thead className="bg-slate-50 dark:bg-slate-900/40">
                                        <tr>
                                            {['Subject','Total','6','5','4','3','2'].map(h => <th key={h} className={thCls}>{h}</th>)}
                                        </tr>
                                    </thead>
                                    <tbody>
                                        {(t.subjects ?? []).map(s => {
                                            const gd = s.absence_summary.grade_distribution
                                            return (
                                                <tr key={s.subject.subject_id} className="border-t border-slate-100 dark:border-slate-700">
                                                    <td className={tdCls}>{s.subject.subject_name}</td>
                                                    <td className={numCls}><span className="font-semibold">{s.absence_summary.number_of_grades}</span></td>
                                                    <td className={numCls}><span className="text-green-600">{gd.number_of_sixes}</span></td>
                                                    <td className={numCls}><span className="text-emerald-600">{gd.number_of_fives}</span></td>
                                                    <td className={numCls}><span className="text-yellow-600">{gd.number_of_fours}</span></td>
                                                    <td className={numCls}><span className="text-orange-500">{gd.number_of_threes}</span></td>
                                                    <td className={numCls}><span className="text-red-500">{gd.number_of_twos}</span></td>
                                                </tr>
                                            )
                                        })}
                                    </tbody>
                                </table>
                            </div>
                        ))}
                    </div>
                ))}
            </div>
        )
    }

    return (
        <div className="flex flex-col gap-6">
            {data.map(school => {
                const rows = school.teachers.map(t => {
                    const totals = (t.subjects ?? []).reduce((acc, s) => {
                        const gd = s.absence_summary.grade_distribution
                        acc.total  += s.absence_summary.number_of_grades
                        acc.sixes  += gd.number_of_sixes
                        acc.fives  += gd.number_of_fives
                        acc.fours  += gd.number_of_fours
                        acc.threes += gd.number_of_threes
                        acc.twos   += gd.number_of_twos
                        return acc
                    }, { total:0, sixes:0, fives:0, fours:0, threes:0, twos:0 })
                    return { teacher: t.teacher, ...totals }
                })
                return (
                    <div key={school.school.school_id} className="bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl overflow-hidden">
                        <div className="px-5 py-3 bg-slate-50 dark:bg-slate-900 border-b border-slate-200 dark:border-slate-700">
                            <h3 className="font-semibold text-slate-800 dark:text-slate-100">{school.school.school_name}</h3>
                        </div>
                        <table className="w-full border-collapse">
                            <thead className="bg-slate-50 dark:bg-slate-900/40">
                                <tr>{['Teacher','Total','6','5','4','3','2'].map(h => <th key={h} className={thCls}>{h}</th>)}</tr>
                            </thead>
                            <tbody>
                                {rows.map(r => (
                                    <tr key={r.teacher.teacher_id} className="border-t border-slate-100 dark:border-slate-700">
                                        <td className={tdCls}>{r.teacher.first_name} {r.teacher.last_name}</td>
                                        <td className={numCls}><span className="font-semibold">{r.total}</span></td>
                                        <td className={numCls}><span className="text-green-600">{r.sixes}</span></td>
                                        <td className={numCls}><span className="text-emerald-600">{r.fives}</span></td>
                                        <td className={numCls}><span className="text-yellow-600">{r.fours}</span></td>
                                        <td className={numCls}><span className="text-orange-500">{r.threes}</span></td>
                                        <td className={numCls}><span className="text-red-500">{r.twos}</span></td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                )
            })}
        </div>
    )
}

const AbsencesTable: React.FC<{ data: SchoolAbsences[]; perSubject: boolean }> = ({ data, perSubject }) => {
    if (!data.length) return <p className="text-center text-slate-400 py-8 text-sm">No absence data found.</p>

    if (perSubject) {
        return (
            <div className="flex flex-col gap-6">
                {data.map(school => (
                    <div key={school.school.school_id} className="bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl overflow-hidden">
                        <div className="px-5 py-3 bg-slate-50 dark:bg-slate-900 border-b border-slate-200 dark:border-slate-700">
                            <h3 className="font-semibold text-slate-800 dark:text-slate-100">{school.school.school_name}</h3>
                        </div>
                        {school.teachers.map(t => (
                            <div key={t.teacher.teacher_id} className="border-b border-slate-100 dark:border-slate-700 last:border-0">
                                <div className="px-5 py-2 bg-indigo-50/50 dark:bg-indigo-900/10">
                                    <p className="text-sm font-medium text-indigo-700 dark:text-indigo-300">{t.teacher.first_name} {t.teacher.last_name}</p>
                                </div>
                                <table className="w-full border-collapse">
                                    <thead className="bg-slate-50 dark:bg-slate-900/40">
                                        <tr>{['Subject','Total','Excused','Unexcused'].map(h => <th key={h} className={thCls}>{h}</th>)}</tr>
                                    </thead>
                                    <tbody>
                                        {(t.subjects ?? []).map(s => (
                                            <tr key={s.subject.subject_id} className="border-t border-slate-100 dark:border-slate-700">
                                                <td className={tdCls}>{s.subject.subject_name}</td>
                                                <td className={numCls}><span className="font-semibold">{s.absence_summary.number_of_absences}</span></td>
                                                <td className={numCls}><span className="text-yellow-600">{s.absence_summary.number_of_excuses}</span></td>
                                                <td className={numCls}><span className="text-red-500">{s.absence_summary.number_of_non_excuses}</span></td>
                                            </tr>
                                        ))}
                                    </tbody>
                                </table>
                            </div>
                        ))}
                    </div>
                ))}
            </div>
        )
    }

    return (
        <div className="flex flex-col gap-6">
            {data.map(school => {
                const rows = school.teachers.map(t => {
                    const totals = (t.subjects ?? []).reduce((acc, s) => {
                        acc.total    += s.absence_summary.number_of_absences
                        acc.excused  += s.absence_summary.number_of_excuses
                        acc.unexcused += s.absence_summary.number_of_non_excuses
                        return acc
                    }, { total:0, excused:0, unexcused:0 })
                    return { teacher: t.teacher, ...totals }
                })
                return (
                    <div key={school.school.school_id} className="bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl overflow-hidden">
                        <div className="px-5 py-3 bg-slate-50 dark:bg-slate-900 border-b border-slate-200 dark:border-slate-700">
                            <h3 className="font-semibold text-slate-800 dark:text-slate-100">{school.school.school_name}</h3>
                        </div>
                        <table className="w-full border-collapse">
                            <thead className="bg-slate-50 dark:bg-slate-900/40">
                                <tr>{['Teacher','Total','Excused','Unexcused'].map(h => <th key={h} className={thCls}>{h}</th>)}</tr>
                            </thead>
                            <tbody>
                                {rows.map(r => (
                                    <tr key={r.teacher.teacher_id} className="border-t border-slate-100 dark:border-slate-700">
                                        <td className={tdCls}>{r.teacher.first_name} {r.teacher.last_name}</td>
                                        <td className={numCls}><span className="font-semibold">{r.total}</span></td>
                                        <td className={numCls}><span className="text-yellow-600">{r.excused}</span></td>
                                        <td className={numCls}><span className="text-red-500">{r.unexcused}</span></td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                )
            })}
        </div>
    )
}

export default function StatisticsPage() {
    const { data: schools = [] } = useGetSchools()
    const [selectedReport, setSelectedReport] = useState<ReportConfig>(REPORTS[0])
    const [selectedSchoolId, setSelectedSchoolId] = useState<string>('')
    const [reportData, setReportData] = useState<ReportData>(null)
    const [error, setError] = useState<string | null>(null)

    const mutation = useMutation({
        mutationFn: async (cfg: ReportConfig) => {
            const body: { report_type: string; school_ids?: number[] } = {
                report_type: cfg.type,
            }
            if (cfg.scope === 'school' && selectedSchoolId) {
                body.school_ids = [Number(selectedSchoolId)]
            }
            const { data } = await axiosInstance.post<{ data: ReportData; error: boolean; message: string }>(
                `${API_URL}/reporting`, body
            )
            if (data.error) throw new Error(data.message)
            return data.data
        },
        onSuccess: (data) => { setReportData(data); setError(null) },
        onError: (err) => { setError(parseApiError(err)); setReportData(null) },
    })

    const handleGenerate = () => {
        if (selectedReport.scope === 'school' && !selectedSchoolId) return
        mutation.mutate(selectedReport)
    }

    const inputCls = 'px-3 py-2 text-sm border border-slate-300 dark:border-slate-600 rounded-md bg-white dark:bg-slate-900 text-slate-800 dark:text-slate-100 focus:outline-none focus:ring-2 focus:ring-indigo-400'
    const btnPrimary = 'px-4 py-2 text-sm font-medium rounded-md bg-indigo-500 hover:bg-indigo-600 text-white border-none cursor-pointer transition-colors disabled:opacity-50'

    return (
        <main className="max-w-5xl mx-auto px-4 py-10 flex flex-col gap-6">
            <h1 className="text-2xl font-semibold text-slate-800 dark:text-slate-100">Statistics</h1>

            <div className="bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl p-5 flex flex-col gap-4">
                <h2 className="text-sm font-semibold text-slate-600 dark:text-slate-300 uppercase tracking-wide">Generate Report</h2>

                <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                    <div className="flex flex-col gap-1">
                        <label className="text-xs font-medium text-slate-500 dark:text-slate-400">Report type</label>
                        <select
                            className={inputCls}
                            value={selectedReport.id}
                            onChange={e => {
                                const cfg = REPORTS.find(r => r.id === e.target.value)!
                                setSelectedReport(cfg)
                                setReportData(null)
                                setError(null)
                                if (cfg.scope !== 'school') setSelectedSchoolId('')
                            }}
                        >
                            {REPORTS.map(r => (
                                <option key={r.id} value={r.id}>{r.label}</option>
                            ))}
                        </select>
                    </div>

                    {selectedReport.scope === 'school' && (
                        <div className="flex flex-col gap-1">
                            <label className="text-xs font-medium text-slate-500 dark:text-slate-400">School</label>
                            <select
                                className={inputCls}
                                value={selectedSchoolId}
                                onChange={e => { setSelectedSchoolId(e.target.value); setReportData(null) }}
                            >
                                <option value="">— Select school —</option>
                                {schools.map((s: School) => (
                                    <option key={s.school_id} value={s.school_id}>{s.school_name}</option>
                                ))}
                            </select>
                        </div>
                    )}
                </div>

                <div className="flex justify-end">
                    <button
                        className={btnPrimary}
                        disabled={mutation.isPending || (selectedReport.scope === 'school' && !selectedSchoolId)}
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
                    <h2 className="text-base font-semibold text-slate-700 dark:text-slate-200">{selectedReport.label}</h2>
                    {selectedReport.type === 'grades' ? (
                        <GradesTable data={reportData as SchoolGrades[]} perSubject={selectedReport.breakdown === 'per_subject'} />
                    ) : (
                        <AbsencesTable data={reportData as SchoolAbsences[]} perSubject={selectedReport.breakdown === 'per_subject'} />
                    )}
                </div>
            )}
        </main>
    )
}
