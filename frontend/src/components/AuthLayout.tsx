import type React from 'react'

interface AuthLayoutProps {
    children: React.ReactNode
}

export const AuthLayout: React.FC<AuthLayoutProps> = ({ children }) => {
    return (
        <div className="min-h-screen bg-slate-100 dark:bg-slate-950 flex items-center justify-center px-4">
            <div className="w-full max-w-sm bg-white dark:bg-slate-800 rounded-xl shadow-sm border border-slate-200 dark:border-slate-700 p-8">
                {children}
            </div>
        </div>
    )
}

