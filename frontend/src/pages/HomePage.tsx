import { useAuth } from '../hooks/useAuth'
import { Link } from 'react-router-dom'

const RoleContent = () => {
    const { user } = useAuth()
    const role = user?.activeRole

    if (role === 'ADMIN') {
        return (
            <div className="space-y-6">
                <h2 className="text-2xl font-semibold text-slate-800 dark:text-slate-100">
                    Administrator Dashboard
                </h2>
                <p className="text-slate-600 dark:text-slate-300">
                    As an administrator, you have full control over the school management system. You can manage schools, principals, teachers, students, parents, classes, terms, subjects, and curricula.
                </p>
                <div className="grid gap-4 sm:grid-cols-2">
                    <Link
                        to="/admin/statistics"
                        className="p-6 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl hover:border-indigo-400 dark:hover:border-indigo-500 transition-colors"
                    >
                        <h3 className="text-lg font-semibold text-slate-800 dark:text-slate-100 mb-2">
                            Statistics
                        </h3>
                        <p className="text-sm text-slate-600 dark:text-slate-400">
                            View system-wide statistics and reports
                        </p>
                    </Link>
                    <Link
                        to="/admin/management"
                        className="p-6 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl hover:border-indigo-400 dark:hover:border-indigo-500 transition-colors"
                    >
                        <h3 className="text-lg font-semibold text-slate-800 dark:text-slate-100 mb-2">
                            Management
                        </h3>
                        <p className="text-sm text-slate-600 dark:text-slate-400">
                            Manage all system entities and resources
                        </p>
                    </Link>
                </div>
            </div>
        )
    }

    if (role === 'DIRECTOR') {
        return (
            <div className="space-y-6">
                <h2 className="text-2xl font-semibold text-slate-800 dark:text-slate-100">
                    Principal Dashboard
                </h2>
                <p className="text-slate-600 dark:text-slate-300">
                    As a principal, you can oversee your school's operations, manage subjects, teachers, students, and parents. Access detailed statistics and reports for your school.
                </p>
                <Link
                    to="/principal"
                    className="block p-6 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl hover:border-indigo-400 dark:hover:border-indigo-500 transition-colors"
                >
                    <h3 className="text-lg font-semibold text-slate-800 dark:text-slate-100 mb-2">
                        Overview
                    </h3>
                    <p className="text-sm text-slate-600 dark:text-slate-400">
                        View your school's overview and quick stats
                    </p>
                </Link>
            </div>
        )
    }

    if (role === 'TEACHER') {
        return (
            <div className="space-y-6">
                <h2 className="text-2xl font-semibold text-slate-800 dark:text-slate-100">
                    Teacher Dashboard
                </h2>
                <p className="text-slate-600 dark:text-slate-300">
                    As a teacher, you can manage your students, record grades and absences, and track academic progress for your classes.
                </p>
                <Link
                    to="/teacher"
                    className="block p-6 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl hover:border-indigo-400 dark:hover:border-indigo-500 transition-colors"
                >
                    <h3 className="text-lg font-semibold text-slate-800 dark:text-slate-100 mb-2">
                        My Students
                    </h3>
                    <p className="text-sm text-slate-600 dark:text-slate-400">
                        View and manage your students' grades and absences
                    </p>
                </Link>
            </div>
        )
    }

    if (role === 'PARENT') {
        return (
            <div className="space-y-6">
                <h2 className="text-2xl font-semibold text-slate-800 dark:text-slate-100">
                    Parent Dashboard
                </h2>
                <p className="text-slate-600 dark:text-slate-300">
                    As a parent, you can monitor your children's academic performance, view their grades, track absences, and stay informed about their progress.
                </p>
                <Link
                    to="/parent"
                    className="block p-6 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl hover:border-indigo-400 dark:hover:border-indigo-500 transition-colors"
                >
                    <h3 className="text-lg font-semibold text-slate-800 dark:text-slate-100 mb-2">
                        My Children
                    </h3>
                    <p className="text-sm text-slate-600 dark:text-slate-400">
                        View your children's grades and absences
                    </p>
                </Link>
            </div>
        )
    }

    if (role === 'STUDENT') {
        return (
            <div className="space-y-6">
                <h2 className="text-2xl font-semibold text-slate-800 dark:text-slate-100">
                    Student Dashboard
                </h2>
                <p className="text-slate-600 dark:text-slate-300">
                    As a student, you can view your academic records, check your grades across all subjects, and track your attendance.
                </p>
                <Link
                    to="/student"
                    className="block p-6 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl hover:border-indigo-400 dark:hover:border-indigo-500 transition-colors"
                >
                    <h3 className="text-lg font-semibold text-slate-800 dark:text-slate-100 mb-2">
                        My Grades and Absences
                    </h3>
                    <p className="text-sm text-slate-600 dark:text-slate-400">
                        View your grades and absences
                    </p>
                </Link>
            </div>
        )
    }

    return (
        <div className="space-y-6">
            <h2 className="text-2xl font-semibold text-slate-800 dark:text-slate-100">
                Welcome to School Manager
            </h2>
            <p className="text-slate-600 dark:text-slate-300">
                Your account is set up. Please contact your administrator to assign you a role.
            </p>
        </div>
    )
}

export default function HomePage() {
    const { user } = useAuth()

    return (
        <main className="max-w-5xl mx-auto px-4 py-10">
            <div className="mb-8">
                <h1 className="text-3xl font-bold text-slate-800 dark:text-slate-100 mb-2">
                    Welcome back, {user?.firstName}!
                </h1>
            </div>
            <RoleContent />
        </main>
    )
}

