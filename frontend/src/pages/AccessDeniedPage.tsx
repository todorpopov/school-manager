import { useNavigate } from 'react-router-dom'
import { useAuth } from '../hooks/useAuth'

export default function AccessDeniedPage() {
    const navigate = useNavigate()
    const { logout } = useAuth()

    return (
        <main className="min-h-screen flex items-center justify-center bg-slate-50 dark:bg-slate-900 px-4">
            <div className="text-center flex flex-col items-center gap-6 max-w-sm">
                <div className="text-6xl font-bold text-red-400">403</div>
                <div className="flex flex-col gap-2">
                    <h1 className="text-2xl font-semibold text-slate-800 dark:text-slate-100">Access Denied</h1>
                    <p className="text-sm text-slate-500 dark:text-slate-400">
                        You don't have permission to view this page.
                    </p>
                </div>
                <div className="flex gap-3">
                    <button
                        onClick={() => navigate('/home')}
                        className="px-4 py-2 text-sm font-medium rounded-md bg-indigo-500 hover:bg-indigo-600 text-white border-none cursor-pointer transition-colors"
                    >
                        Go to Home
                    </button>
                    <button
                        onClick={logout}
                        className="px-4 py-2 text-sm font-medium rounded-md border border-slate-300 dark:border-slate-600 text-slate-600 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-700 bg-transparent cursor-pointer transition-colors"
                    >
                        Log out
                    </button>
                </div>
            </div>
        </main>
    )
}

