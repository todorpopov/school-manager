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

    return (
        <div className="min-h-screen bg-slate-100 dark:bg-slate-950">
            <header className="bg-white dark:bg-slate-800 border-b border-slate-200 dark:border-slate-700 px-6 py-3 flex items-center justify-between">
                <div className="flex items-center gap-6">
                    <span className="font-semibold text-slate-800 dark:text-slate-100">School Manager</span>
                    {isAdmin && (
                        <nav className="flex items-center gap-4">
                            <NavLink
                                to="/admin/statistics"
                                className={({ isActive }) =>
                                    `text-sm font-medium transition-colors ${
                                        isActive
                                            ? 'text-indigo-600 dark:text-indigo-400'
                                            : 'text-slate-500 hover:text-slate-800 dark:hover:text-slate-100'
                                    }`
                                }
                            >
                                Statistics
                            </NavLink>
                            <NavLink
                                to="/admin/management"
                                className={({ isActive }) =>
                                    `text-sm font-medium transition-colors ${
                                        isActive
                                            ? 'text-indigo-600 dark:text-indigo-400'
                                            : 'text-slate-500 hover:text-slate-800 dark:hover:text-slate-100'
                                    }`
                                }
                            >
                                Management
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

