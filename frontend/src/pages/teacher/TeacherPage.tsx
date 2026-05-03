import React, { useState } from 'react'
import { useGetStudents } from '../../hooks/useStudents'
import { useGetCurricula } from '../../hooks/useCurricula'
import {
    useGetAllGrades, useGetAllAbsences,
    useCreateGrade, useDeleteGrade,
    useCreateAbsence, useDeleteAbsence,
    useBulkCreateGrades, useBulkCreateAbsences,
} from '../../hooks/useGradesAbsences'
import type { Grade, Absence } from '../../types/gradesAbsences'
import { Toast } from '../../components/Toast'
import { useToast } from '../../hooks/useToast'
import type { Student } from '../../types/students'
import type { Curriculum } from '../../types/curricula'

type Tab = 'grades' | 'absences'
type Mode = 'individual' | 'bulk'

const inputCls = 'w-full px-3 py-2 text-sm border border-slate-300 dark:border-slate-600 rounded-md bg-white dark:bg-slate-900 text-slate-800 dark:text-slate-100 focus:outline-none focus:ring-2 focus:ring-indigo-400'
const btnPrimary = 'px-3 py-1.5 text-sm font-medium rounded-md bg-indigo-500 hover:bg-indigo-600 text-white border-none cursor-pointer transition-colors disabled:opacity-50'
const btnDanger  = 'px-2 py-1 text-xs font-medium rounded border border-red-200 dark:border-red-800 text-red-500 hover:bg-red-50 dark:hover:bg-red-900/30 cursor-pointer bg-transparent transition-colors'
const btnGhost   = 'px-3 py-1.5 text-sm font-medium rounded-md border border-slate-300 dark:border-slate-600 text-slate-600 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-700 cursor-pointer bg-transparent transition-colors'
const thCls = 'px-4 py-3 text-left text-xs font-semibold text-slate-500 dark:text-slate-400 uppercase tracking-wide whitespace-nowrap'
const tdCls = 'px-4 py-3 text-slate-700 dark:text-slate-300'

interface GradeFormState { curriculum_id: string; grade_value: string; grade_date: string }
interface AbsenceFormState { curriculum_id: string; absence_date: string; is_excused: boolean }

interface BulkGradeRow { student_id: number; first_name: string; last_name: string; grade_value: string; include: boolean }
interface BulkAbsenceRow { student_id: number; first_name: string; last_name: string; is_excused: boolean; include: boolean }

const TeacherPage: React.FC = () => {
    const { data: students = [], isLoading: loadingStudents } = useGetStudents()
    const { data: curricula = [] } = useGetCurricula()
    const { data: allGrades = [], isLoading: loadingGrades } = useGetAllGrades()
    const { data: allAbsences = [], isLoading: loadingAbsences } = useGetAllAbsences()

    const createGrade   = useCreateGrade()
    const deleteGrade   = useDeleteGrade()
    const createAbsence = useCreateAbsence()
    const deleteAbsence = useDeleteAbsence()
    const bulkCreateGrades   = useBulkCreateGrades()
    const bulkCreateAbsences = useBulkCreateAbsences()

    const { toast, show, dismiss } = useToast()

    const [mode, setMode] = useState<Mode>('individual')
    const [selectedStudent, setSelectedStudent] = useState<Student | null>(null)
    const [tab, setTab] = useState<Tab>('grades')
    const [showGradeForm, setShowGradeForm] = useState(false)
    const [showAbsenceForm, setShowAbsenceForm] = useState(false)
    const [gradeForm, setGradeForm] = useState<GradeFormState>({ curriculum_id: '', grade_value: '', grade_date: '' })
    const [absenceForm, setAbsenceForm] = useState<AbsenceFormState>({ curriculum_id: '', absence_date: '', is_excused: false })

    const [bulkTab, setBulkTab] = useState<Tab>('grades')
    const [bulkClassId, setBulkClassId] = useState('')
    const [bulkCurriculumId, setBulkCurriculumId] = useState('')
    const [bulkDate, setBulkDate] = useState('')
    const [bulkGradeRows, setBulkGradeRows] = useState<BulkGradeRow[]>([])
    const [bulkAbsenceRows, setBulkAbsenceRows] = useState<BulkAbsenceRow[]>([])

    const today = new Date().toISOString().split('T')[0]

    const classOptions = Array.from(
        new Map(curricula.map((c: Curriculum) => [c.class.class_id, c.class])).values()
    ).sort((a, b) => a.grade_level - b.grade_level || a.class_name.localeCompare(b.class_name))

    const curriculaForClass = bulkClassId
        ? curricula.filter((c: Curriculum) => String(c.class.class_id) === bulkClassId)
        : []

    const studentsInClass = bulkClassId
        ? students.filter(s => s.class && String((s.class as { class_id: number }).class_id) === bulkClassId)
        : []

    const curriculaOptions = curricula.map((c: Curriculum) => ({
        value: c.curriculum_id,
        label: `${c.subject.subject_name} — ${c.term.name}`,
    }))

    const studentGrades   = selectedStudent ? allGrades.filter(g => g.student.student_id === selectedStudent.student_id) : []
    const studentAbsences = selectedStudent ? allAbsences.filter(a => a.student.student_id === selectedStudent.student_id) : []

    const handleCreateGrade = async () => {
        if (!selectedStudent || !gradeForm.curriculum_id || !gradeForm.grade_value || !gradeForm.grade_date) {
            show('Please fill in all fields', 'error'); return
        }
        if (gradeForm.grade_date > today) { show('Grade date cannot be in the future', 'error'); return }
        const val = parseFloat(gradeForm.grade_value)
        if (val < 2 || val > 6) { show('Grade must be between 2.00 and 6.00', 'error'); return }
        try {
            await createGrade.mutateAsync({
                student_id: selectedStudent.student_id,
                curriculum_id: Number(gradeForm.curriculum_id),
                grade_value: val,
                grade_date: gradeForm.grade_date,
            })
            setGradeForm({ curriculum_id: '', grade_value: '', grade_date: '' })
            setShowGradeForm(false)
        } catch (e) { show((e as Error).message, 'error') }
    }

    const handleCreateAbsence = async () => {
        if (!selectedStudent || !absenceForm.curriculum_id || !absenceForm.absence_date) {
            show('Please fill in all fields', 'error'); return
        }
        if (absenceForm.absence_date > today) { show('Absence date cannot be in the future', 'error'); return }
        try {
            await createAbsence.mutateAsync({
                student_id: selectedStudent.student_id,
                curriculum_id: Number(absenceForm.curriculum_id),
                absence_date: absenceForm.absence_date,
                is_excused: absenceForm.is_excused,
            })
            setAbsenceForm({ curriculum_id: '', absence_date: '', is_excused: false })
            setShowAbsenceForm(false)
        } catch (e) { show((e as Error).message, 'error') }
    }

    const handleDeleteGrade = async (id: number) => {
        try { await deleteGrade.mutateAsync(id) } catch (e) { show((e as Error).message, 'error') }
    }

    const handleDeleteAbsence = async (id: number) => {
        try { await deleteAbsence.mutateAsync(id) } catch (e) { show((e as Error).message, 'error') }
    }

    const handleBulkClassChange = (classId: string) => {
        setBulkClassId(classId)
        setBulkCurriculumId('')
        const classStudents = classId
            ? students.filter(s => s.class && String((s.class as { class_id: number }).class_id) === classId)
            : []
        setBulkGradeRows(classStudents.map(s => ({
            student_id: s.student_id, first_name: s.first_name as string, last_name: s.last_name as string,
            grade_value: '', include: true,
        })))
        setBulkAbsenceRows(classStudents.map(s => ({
            student_id: s.student_id, first_name: s.first_name as string, last_name: s.last_name as string,
            is_excused: false, include: false,
        })))
    }

    const handleBulkSubmitGrades = async () => {
        if (!bulkCurriculumId || !bulkDate) { show('Please select a curriculum and date', 'error'); return }
        if (bulkDate > today) { show('Date cannot be in the future', 'error'); return }
        const included = bulkGradeRows.filter(r => r.include)
        if (!included.length) { show('No students selected', 'error'); return }
        for (const r of included) {
            const val = parseFloat(r.grade_value)
            if (isNaN(val) || val < 2 || val > 6) { show(`Invalid grade for ${r.first_name} ${r.last_name}`, 'error'); return }
        }
        try {
            await bulkCreateGrades.mutateAsync({
                entries: included.map(r => ({
                    student_id: r.student_id,
                    curriculum_id: Number(bulkCurriculumId),
                    grade_value: parseFloat(r.grade_value),
                    grade_date: bulkDate,
                })),
            })
            show(`${included.length} grade(s) created`, 'success')
            setBulkGradeRows(prev => prev.map(r => ({ ...r, grade_value: '', include: true })))
        } catch (e) { show((e as Error).message, 'error') }
    }

    const handleBulkSubmitAbsences = async () => {
        if (!bulkCurriculumId || !bulkDate) { show('Please select a curriculum and date', 'error'); return }
        if (bulkDate > today) { show('Date cannot be in the future', 'error'); return }
        const included = bulkAbsenceRows.filter(r => r.include)
        if (!included.length) { show('No students selected', 'error'); return }
        try {
            await bulkCreateAbsences.mutateAsync({
                entries: included.map(r => ({
                    student_id: r.student_id,
                    curriculum_id: Number(bulkCurriculumId),
                    absence_date: bulkDate,
                    is_excused: r.is_excused,
                })),
            })
            show(`${included.length} absence(s) created`, 'success')
            setBulkAbsenceRows(prev => prev.map(r => ({ ...r, include: false })))
        } catch (e) { show((e as Error).message, 'error') }
    }

    return (
        <main className="max-w-5xl mx-auto px-4 py-10 flex flex-col gap-6">
            {toast && <Toast message={toast.message} variant={toast.variant} onDismiss={dismiss} />}

            <div className="flex items-center justify-between">
                <h1 className="text-2xl font-semibold text-slate-800 dark:text-slate-100">Students</h1>
                <div className="flex gap-1 bg-slate-100 dark:bg-slate-800 p-1 rounded-lg">
                    {(['individual', 'bulk'] as Mode[]).map(m => (
                        <button key={m} onClick={() => setMode(m)}
                            className={`px-4 py-1.5 text-sm font-medium rounded-md capitalize transition-colors cursor-pointer border-none ${
                                mode === m ? 'bg-white dark:bg-slate-700 text-indigo-600 dark:text-indigo-400 shadow-sm' : 'text-slate-500 hover:text-slate-700 dark:hover:text-slate-300 bg-transparent'
                            }`}>
                            {m === 'bulk' ? 'Bulk Entry' : 'Individual'}
                        </button>
                    ))}
                </div>
            </div>
            {mode === 'individual' && (
                <div className="flex gap-6">
                    <div className="w-56 flex-shrink-0">
                        <div className="bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl overflow-hidden">
                            <div className="px-4 py-3 bg-slate-50 dark:bg-slate-900 border-b border-slate-200 dark:border-slate-700 text-xs font-semibold text-slate-500 dark:text-slate-400 uppercase tracking-wide">
                                Select Student
                            </div>
                            {loadingStudents && <p className="text-center text-slate-400 py-6 text-sm">Loading…</p>}
                            {!loadingStudents && students.length === 0 && <p className="text-center text-slate-400 py-6 text-sm">No students.</p>}
                            <ul>
                                {students.map(s => (
                                    <li key={s.student_id}>
                                        <button
                                            onClick={() => { setSelectedStudent(s as Student); setTab('grades'); setShowGradeForm(false); setShowAbsenceForm(false) }}
                                            className={`w-full text-left px-4 py-2.5 text-sm border-b border-slate-100 dark:border-slate-700 last:border-0 transition-colors cursor-pointer ${
                                                selectedStudent?.student_id === s.student_id
                                                    ? 'bg-indigo-50 dark:bg-indigo-900/20 text-indigo-700 dark:text-indigo-300 font-medium'
                                                    : 'hover:bg-slate-50 dark:hover:bg-slate-700 text-slate-700 dark:text-slate-300'
                                            }`}>
                                            {s.first_name as string} {s.last_name as string}
                                        </button>
                                    </li>
                                ))}
                            </ul>
                        </div>
                    </div>

                    <div className="flex-1 min-w-0">
                        {!selectedStudent ? (
                            <div className="flex items-center justify-center h-48 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-slate-400 text-sm">
                                Select a student to view their records
                            </div>
                        ) : (
                            <div className="bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl overflow-hidden">
                                <div className="px-5 py-4 bg-slate-50 dark:bg-slate-900 border-b border-slate-200 dark:border-slate-700">
                                    <h2 className="text-base font-semibold text-slate-800 dark:text-slate-100">
                                        {selectedStudent.first_name as string} {selectedStudent.last_name as string}
                                    </h2>
                                </div>
                                <div className="flex border-b border-slate-200 dark:border-slate-700">
                                    {(['grades', 'absences'] as Tab[]).map(t => (
                                        <button key={t} onClick={() => { setTab(t); setShowGradeForm(false); setShowAbsenceForm(false) }}
                                            className={`px-6 py-3 text-sm font-medium capitalize border-b-2 transition-colors cursor-pointer bg-transparent ${
                                                tab === t ? 'border-indigo-500 text-indigo-600 dark:text-indigo-400' : 'border-transparent text-slate-500 hover:text-slate-700 dark:hover:text-slate-300'
                                            }`}>
                                            {t}
                                        </button>
                                    ))}
                                </div>

                                {tab === 'grades' && (
                                    <div className="flex flex-col">
                                        <div className="flex justify-end px-4 py-3 border-b border-slate-100 dark:border-slate-700">
                                            <button className={btnPrimary} onClick={() => setShowGradeForm(v => !v)}>
                                                {showGradeForm ? 'Cancel' : '+ Add Grade'}
                                            </button>
                                        </div>
                                        {showGradeForm && (
                                            <div className="flex flex-col gap-3 px-4 py-4 border-b border-slate-100 dark:border-slate-700 bg-slate-50 dark:bg-slate-900/40">
                                                <div className="grid grid-cols-3 gap-3">
                                                    <div className="flex flex-col gap-1">
                                                        <label className="text-xs font-medium text-slate-500 dark:text-slate-400">Curriculum</label>
                                                        <select className={inputCls} value={gradeForm.curriculum_id} onChange={e => setGradeForm(f => ({ ...f, curriculum_id: e.target.value }))}>
                                                            <option value="">— Select —</option>
                                                            {curriculaOptions.map(o => <option key={o.value} value={o.value}>{o.label}</option>)}
                                                        </select>
                                                    </div>
                                                    <div className="flex flex-col gap-1">
                                                        <label className="text-xs font-medium text-slate-500 dark:text-slate-400">Grade (2.00 – 6.00)</label>
                                                        <input type="number" min={2} max={6} step={0.01} className={inputCls}
                                                            value={gradeForm.grade_value}
                                                            onChange={e => setGradeForm(f => ({ ...f, grade_value: e.target.value }))} />
                                                    </div>
                                                    <div className="flex flex-col gap-1">
                                                        <label className="text-xs font-medium text-slate-500 dark:text-slate-400">Date</label>
                                                        <input type="date" className={inputCls} max={today}
                                                            value={gradeForm.grade_date}
                                                            onChange={e => setGradeForm(f => ({ ...f, grade_date: e.target.value }))} />
                                                    </div>
                                                </div>
                                                <div className="flex justify-end gap-2">
                                                    <button className={btnGhost} onClick={() => setShowGradeForm(false)}>Cancel</button>
                                                    <button className={btnPrimary} disabled={createGrade.isPending} onClick={handleCreateGrade}>
                                                        {createGrade.isPending ? 'Saving…' : 'Save Grade'}
                                                    </button>
                                                </div>
                                            </div>
                                        )}
                                        {loadingGrades ? (
                                            <div className="flex justify-center py-8"><div className="w-6 h-6 border-4 border-slate-200 border-t-indigo-500 rounded-full animate-spin" /></div>
                                        ) : studentGrades.length === 0 ? (
                                            <p className="text-center text-slate-400 py-8 text-sm">No grades recorded.</p>
                                        ) : (
                                            <table className="w-full border-collapse text-sm">
                                                <thead className="bg-slate-100 dark:bg-slate-900">
                                                    <tr>{['Subject', 'Term', 'Grade', 'Date', ''].map(h => <th key={h} className={thCls}>{h}</th>)}</tr>
                                                </thead>
                                                <tbody>
                                                    {studentGrades.map((g: Grade) => (
                                                        <tr key={g.grade_id} className="border-t border-slate-100 dark:border-slate-700 hover:bg-slate-50 dark:hover:bg-slate-900/50">
                                                            <td className={tdCls}>{g.curriculum.subject.subject_name}</td>
                                                            <td className={tdCls}>{g.curriculum.term.name}</td>
                                                            <td className={tdCls}><span className={`font-semibold ${g.grade_value >= 5 ? 'text-green-600' : g.grade_value >= 3.5 ? 'text-yellow-600' : 'text-red-500'}`}>{g.grade_value.toFixed(2)}</span></td>
                                                            <td className={tdCls}>{new Date(g.grade_date).toLocaleDateString()}</td>
                                                            <td className="px-4 py-3 text-right"><button className={btnDanger} disabled={deleteGrade.isPending} onClick={() => handleDeleteGrade(g.grade_id)}>Delete</button></td>
                                                        </tr>
                                                    ))}
                                                </tbody>
                                            </table>
                                        )}
                                    </div>
                                )}

                                {tab === 'absences' && (
                                    <div className="flex flex-col">
                                        <div className="flex justify-end px-4 py-3 border-b border-slate-100 dark:border-slate-700">
                                            <button className={btnPrimary} onClick={() => setShowAbsenceForm(v => !v)}>
                                                {showAbsenceForm ? 'Cancel' : '+ Add Absence'}
                                            </button>
                                        </div>
                                        {showAbsenceForm && (
                                            <div className="flex flex-col gap-3 px-4 py-4 border-b border-slate-100 dark:border-slate-700 bg-slate-50 dark:bg-slate-900/40">
                                                <div className="grid grid-cols-3 gap-3">
                                                    <div className="flex flex-col gap-1">
                                                        <label className="text-xs font-medium text-slate-500 dark:text-slate-400">Curriculum</label>
                                                        <select className={inputCls} value={absenceForm.curriculum_id} onChange={e => setAbsenceForm(f => ({ ...f, curriculum_id: e.target.value }))}>
                                                            <option value="">— Select —</option>
                                                            {curriculaOptions.map(o => <option key={o.value} value={o.value}>{o.label}</option>)}
                                                        </select>
                                                    </div>
                                                    <div className="flex flex-col gap-1">
                                                        <label className="text-xs font-medium text-slate-500 dark:text-slate-400">Date</label>
                                                        <input type="date" className={inputCls} max={today}
                                                            value={absenceForm.absence_date}
                                                            onChange={e => setAbsenceForm(f => ({ ...f, absence_date: e.target.value }))} />
                                                    </div>
                                                    <div className="flex flex-col gap-1 justify-end">
                                                        <label className="flex items-center gap-2 text-sm text-slate-600 dark:text-slate-300 cursor-pointer">
                                                            <input type="checkbox" className="accent-indigo-500" checked={absenceForm.is_excused} onChange={e => setAbsenceForm(f => ({ ...f, is_excused: e.target.checked }))} />
                                                            Excused
                                                        </label>
                                                    </div>
                                                </div>
                                                <div className="flex justify-end gap-2">
                                                    <button className={btnGhost} onClick={() => setShowAbsenceForm(false)}>Cancel</button>
                                                    <button className={btnPrimary} disabled={createAbsence.isPending} onClick={handleCreateAbsence}>
                                                        {createAbsence.isPending ? 'Saving…' : 'Save Absence'}
                                                    </button>
                                                </div>
                                            </div>
                                        )}
                                        {loadingAbsences ? (
                                            <div className="flex justify-center py-8"><div className="w-6 h-6 border-4 border-slate-200 border-t-indigo-500 rounded-full animate-spin" /></div>
                                        ) : studentAbsences.length === 0 ? (
                                            <p className="text-center text-slate-400 py-8 text-sm">No absences recorded.</p>
                                        ) : (
                                            <table className="w-full border-collapse text-sm">
                                                <thead className="bg-slate-100 dark:bg-slate-900">
                                                    <tr>{['Subject', 'Term', 'Date', 'Status', ''].map(h => <th key={h} className={thCls}>{h}</th>)}</tr>
                                                </thead>
                                                <tbody>
                                                    {studentAbsences.map((a: Absence) => (
                                                        <tr key={a.absence_id} className="border-t border-slate-100 dark:border-slate-700 hover:bg-slate-50 dark:hover:bg-slate-900/50">
                                                            <td className={tdCls}>{a.curriculum.subject.subject_name}</td>
                                                            <td className={tdCls}>{a.curriculum.term.name}</td>
                                                            <td className={tdCls}>{new Date(a.absence_date).toLocaleDateString()}</td>
                                                            <td className={tdCls}>
                                                                <span className={`text-xs font-medium px-2 py-0.5 rounded-full ${a.is_excused ? 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900/30 dark:text-yellow-400' : 'bg-red-100 text-red-600 dark:bg-red-900/30 dark:text-red-400'}`}>
                                                                    {a.is_excused ? 'Excused' : 'Unexcused'}
                                                                </span>
                                                            </td>
                                                            <td className="px-4 py-3 text-right"><button className={btnDanger} disabled={deleteAbsence.isPending} onClick={() => handleDeleteAbsence(a.absence_id)}>Delete</button></td>
                                                        </tr>
                                                    ))}
                                                </tbody>
                                            </table>
                                        )}
                                    </div>
                                )}
                            </div>
                        )}
                    </div>
                </div>
            )}
            {mode === 'bulk' && (
                <div className="bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl overflow-hidden">
                    <div className="flex border-b border-slate-200 dark:border-slate-700">
                        {(['grades', 'absences'] as Tab[]).map(t => (
                            <button key={t} onClick={() => setBulkTab(t)}
                                className={`px-6 py-3 text-sm font-medium capitalize border-b-2 transition-colors cursor-pointer bg-transparent ${
                                    bulkTab === t ? 'border-indigo-500 text-indigo-600 dark:text-indigo-400' : 'border-transparent text-slate-500 hover:text-slate-700 dark:hover:text-slate-300'
                                }`}>
                                {t}
                            </button>
                        ))}
                    </div>
                    <div className="grid grid-cols-3 gap-4 px-5 py-4 bg-slate-50 dark:bg-slate-900/40 border-b border-slate-100 dark:border-slate-700">
                        <div className="flex flex-col gap-1">
                            <label className="text-xs font-medium text-slate-500 dark:text-slate-400">Class</label>
                            <select className={inputCls} value={bulkClassId} onChange={e => handleBulkClassChange(e.target.value)}>
                                <option value="">— Select class —</option>
                                {classOptions.map(cl => (
                                    <option key={cl.class_id} value={cl.class_id}>{cl.grade_level}{cl.class_name}</option>
                                ))}
                            </select>
                        </div>
                        <div className="flex flex-col gap-1">
                            <label className="text-xs font-medium text-slate-500 dark:text-slate-400">Curriculum</label>
                            <select className={inputCls} value={bulkCurriculumId} onChange={e => setBulkCurriculumId(e.target.value)} disabled={!bulkClassId}>
                                <option value="">— Select curriculum —</option>
                                {curriculaForClass.map((c: Curriculum) => (
                                    <option key={c.curriculum_id} value={c.curriculum_id}>{c.subject.subject_name} — {c.term.name}</option>
                                ))}
                            </select>
                        </div>
                        <div className="flex flex-col gap-1">
                            <label className="text-xs font-medium text-slate-500 dark:text-slate-400">Date</label>
                            <input type="date" className={inputCls} max={today} value={bulkDate} onChange={e => setBulkDate(e.target.value)} />
                        </div>
                    </div>
                    {bulkTab === 'grades' && (
                        <>
                            {!bulkClassId ? (
                                <p className="text-center text-slate-400 py-10 text-sm">Select a class to start bulk grade entry.</p>
                            ) : studentsInClass.length === 0 ? (
                                <p className="text-center text-slate-400 py-10 text-sm">No students in this class.</p>
                            ) : (
                                <>
                                    <table className="w-full border-collapse text-sm">
                                        <thead className="bg-slate-100 dark:bg-slate-900">
                                            <tr>
                                                <th className={thCls}>Include</th>
                                                <th className={thCls}>Student</th>
                                                <th className={thCls}>Grade (2.00 – 6.00)</th>
                                            </tr>
                                        </thead>
                                        <tbody>
                                            {bulkGradeRows.map((row, i) => (
                                                <tr key={row.student_id} className="border-t border-slate-100 dark:border-slate-700">
                                                    <td className="px-4 py-2.5 text-center">
                                                        <input type="checkbox" className="accent-indigo-500" checked={row.include}
                                                            onChange={e => setBulkGradeRows(prev => prev.map((r, idx) => idx === i ? { ...r, include: e.target.checked } : r))} />
                                                    </td>
                                                    <td className={tdCls}>{row.first_name} {row.last_name}</td>
                                                    <td className="px-4 py-2">
                                                        <input type="number" min={2} max={6} step={0.01} placeholder="e.g. 5.50"
                                                            className={`${inputCls} max-w-[120px]`}
                                                            value={row.grade_value}
                                                            disabled={!row.include}
                                                            onChange={e => setBulkGradeRows(prev => prev.map((r, idx) => idx === i ? { ...r, grade_value: e.target.value } : r))} />
                                                    </td>
                                                </tr>
                                            ))}
                                        </tbody>
                                    </table>
                                    <div className="flex justify-between items-center px-5 py-3 border-t border-slate-100 dark:border-slate-700">
                                        <span className="text-xs text-slate-400">{bulkGradeRows.filter(r => r.include).length} student(s) selected</span>
                                        <button className={btnPrimary} disabled={bulkCreateGrades.isPending} onClick={handleBulkSubmitGrades}>
                                            {bulkCreateGrades.isPending ? 'Saving…' : 'Save All Grades'}
                                        </button>
                                    </div>
                                </>
                            )}
                        </>
                    )}
                    {bulkTab === 'absences' && (
                        <>
                            {!bulkClassId ? (
                                <p className="text-center text-slate-400 py-10 text-sm">Select a class to start bulk absence entry.</p>
                            ) : studentsInClass.length === 0 ? (
                                <p className="text-center text-slate-400 py-10 text-sm">No students in this class.</p>
                            ) : (
                                <>
                                    <table className="w-full border-collapse text-sm">
                                        <thead className="bg-slate-100 dark:bg-slate-900">
                                            <tr>
                                                <th className={thCls}>Absent</th>
                                                <th className={thCls}>Student</th>
                                                <th className={thCls}>Excused</th>
                                            </tr>
                                        </thead>
                                        <tbody>
                                            {bulkAbsenceRows.map((row, i) => (
                                                <tr key={row.student_id} className="border-t border-slate-100 dark:border-slate-700">
                                                    <td className="px-4 py-2.5 text-center">
                                                        <input type="checkbox" className="accent-indigo-500" checked={row.include}
                                                            onChange={e => setBulkAbsenceRows(prev => prev.map((r, idx) => idx === i ? { ...r, include: e.target.checked } : r))} />
                                                    </td>
                                                    <td className={tdCls}>{row.first_name} {row.last_name}</td>
                                                    <td className="px-4 py-2.5">
                                                        <input type="checkbox" className="accent-yellow-500" checked={row.is_excused}
                                                            disabled={!row.include}
                                                            onChange={e => setBulkAbsenceRows(prev => prev.map((r, idx) => idx === i ? { ...r, is_excused: e.target.checked } : r))} />
                                                    </td>
                                                </tr>
                                            ))}
                                        </tbody>
                                    </table>
                                    <div className="flex justify-between items-center px-5 py-3 border-t border-slate-100 dark:border-slate-700">
                                        <span className="text-xs text-slate-400">{bulkAbsenceRows.filter(r => r.include).length} student(s) marked absent</span>
                                        <button className={btnPrimary} disabled={bulkCreateAbsences.isPending} onClick={handleBulkSubmitAbsences}>
                                            {bulkCreateAbsences.isPending ? 'Saving…' : 'Save All Absences'}
                                        </button>
                                    </div>
                                </>
                            )}
                        </>
                    )}
                </div>
            )}
        </main>
    )
}

export default TeacherPage
