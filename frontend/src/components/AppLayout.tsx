import { NavLink, Outlet } from 'react-router-dom'
import { useAuth } from '../hooks/useAuth'

const ROLE_LABELS: Record<string, string> = {
    ADMIN: 'Administrator',
    DIRECTOR: 'Principal',
    TEACHER: 'Teacher',
    PARENT: 'Parent',
    STUDENT: 'Student',
    USER: 'User',
}

export default function AppLayout() {
    const { user: auth, logout } = useAuth()
    const isAdmin = auth?.activeRole === 'ADMIN'
    const isPrincipal = auth?.activeRole === 'DIRECTOR'
    const isParent = auth?.activeRole === 'PARENT'
    const isTeacher = auth?.activeRole === 'TEACHER'
    const isStudent = auth?.activeRole === 'STUDENT'

    return (
        <div className="min-h-screen bg-slate-100 dark:bg-slate-950">
            <header className="bg-white dark:bg-slate-800 border-b border-slate-200 dark:border-slate-700 px-6 py-3 flex items-center justify-between">
                <div className="flex items-center gap-6">
                    <NavLink
                        to="/home"
                        className="font-semibold text-slate-800 dark:text-slate-100 hover:text-indigo-600 dark:hover:text-indigo-400 transition-colors"
                    >
                        School Manager
                    </NavLink>
                    {isAdmin && (
                        <nav className="flex items-center gap-4">
                            <NavLink
                                to="/admin/statistics"
                                className={({ isActive }) =>
                                    `text-sm font-medium transition-colors ${isActive ? 'text-indigo-600 dark:text-indigo-400' : 'text-slate-500 hover:text-slate-800 dark:hover:text-slate-100'}`
                                }
                            >
                                Statistics
                            </NavLink>
                            <NavLink
                                to="/admin/management"
                                className={({ isActive }) =>
                                    `text-sm font-medium transition-colors ${isActive ? 'text-indigo-600 dark:text-indigo-400' : 'text-slate-500 hover:text-slate-800 dark:hover:text-slate-100'}`
                                }
                            >
                                Management
                            </NavLink>
                        </nav>
                    )}
                    {isParent && (
                        <nav className="flex items-center gap-4">
                            <NavLink to="/parent" className={({ isActive }) => `text-sm font-medium transition-colors ${isActive ? 'text-indigo-600 dark:text-indigo-400' : 'text-slate-500 hover:text-slate-800 dark:hover:text-slate-100'}`}>
                                My Children
                            </NavLink>
                        </nav>
                    )}
                    {isTeacher && (
                        <nav className="flex items-center gap-4">
                            <NavLink to="/teacher" className={({ isActive }) => `text-sm font-medium transition-colors ${isActive ? 'text-indigo-600 dark:text-indigo-400' : 'text-slate-500 hover:text-slate-800 dark:hover:text-slate-100'}`}>
                                My Students
                            </NavLink>
                        </nav>
                    )}
                    {isStudent && (
                        <nav className="flex items-center gap-4">
                            <NavLink to="/student" className={({ isActive }) => `text-sm font-medium transition-colors ${isActive ? 'text-indigo-600 dark:text-indigo-400' : 'text-slate-500 hover:text-slate-800 dark:hover:text-slate-100'}`}>
                                My Grades
                            </NavLink>
                        </nav>
                    )}
                    {isPrincipal && (
                        <nav className="flex items-center gap-4">
                            <NavLink
                                to="/principal"
                                className={({ isActive }) =>
                                    `text-sm font-medium transition-colors ${isActive ? 'text-indigo-600 dark:text-indigo-400' : 'text-slate-500 hover:text-slate-800 dark:hover:text-slate-100'}`
                                }
                            >
                                Overview
                            </NavLink>
                        </nav>
                    )}
                </div>
                <div className="flex items-center gap-4">
                    <span className="text-sm text-slate-500 dark:text-slate-400">
                        {auth?.firstName} {auth?.lastName}
                        <span className="ml-2 px-2 py-0.5 bg-indigo-100 dark:bg-indigo-900 text-indigo-700 dark:text-indigo-300 rounded-full text-xs font-medium">
                            {ROLE_LABELS[auth?.activeRole ?? ''] ?? auth?.activeRole}
                        </span>
                    </span>
                    <button
                        onClick={logout}
                        className="text-sm text-slate-500 hover:text-slate-800 dark:hover:text-slate-100 transition-colors cursor-pointer bg-transparent border-none"
                    >
                        Sign out
                    </button>
                </div>
            </header>
            <Outlet />
        </div>
    )
}

