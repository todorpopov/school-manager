import type React from 'react';
import { useState } from 'react'
import { Link } from 'react-router-dom'
import { useAuth } from '../hooks/useAuth'
import { apiLogin } from '../api/auth'
import { useToast } from '../hooks/useToast'
import { Toast } from '../components/Toast'
import { Field } from '../components/Field'
import { AuthLayout } from '../components/AuthLayout'
import { inputCls, primaryBtnCls } from '../styles/formStyles'

const ROLE_LABELS: Record<string, string> = {
    ADMIN: 'Administrator',
    DIRECTOR: 'Director',
    TEACHER: 'Teacher',
    PARENT: 'Parent',
    STUDENT: 'Student',
    USER: 'User',
}

export default function LoginPage() {
    const { login, pendingAuth, selectRole } = useAuth()

    const { toast, show: showToast, dismiss } = useToast()

    const [email, setEmail] = useState('')
    const [password, setPassword] = useState('')
    const [errors, setErrors] = useState<{ email?: string; password?: string }>({})
    const [loading, setLoading] = useState(false)

    const validate = () => {
        const e: typeof errors = {}

        if (!email) {
            e.email = 'Email is required'
        } else if (!email.includes('@')) {
            e.email = 'Email must contain @'
        }

        if (!password) {
            e.password = 'Password is required'
        }

        setErrors(e)
        return Object.keys(e).length === 0
    }

    const handleEmailChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setEmail(e.target.value)
        setErrors((prev) => ({ ...prev, email: undefined }))
    }

    const handlePasswordChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setPassword(e.target.value)
        setErrors((prev) => ({ ...prev, password: undefined }))
    }

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault()

        if (!validate()) {
            return
        }

        setLoading(true)

        try {
            const response = await apiLogin(email, password)
            login(response)
        } catch (err) {
            showToast(err instanceof Error ? err.message : 'Login failed')
        } finally {
            setLoading(false)
        }
    }

    const handleRoleSelect = (role: string) => {
        selectRole(role as Parameters<typeof selectRole>[0])
    }

    if (pendingAuth) {
        return (
            <AuthLayout>
                <h1 className="text-2xl font-semibold text-slate-800 dark:text-slate-100 mb-4 text-center">
                    Log in as…
                </h1>
                <p className="text-sm text-slate-500 dark:text-slate-400 mb-6">
                    Your account has multiple roles. Choose how you want to continue.
                </p>
                <div className="flex flex-col gap-3">
                    {pendingAuth.roles.map((role) => (
                        <button
                            key={role}
                            onClick={() => handleRoleSelect(role)}
                            className="w-full py-3 px-4 rounded-lg border-2 border-slate-200 dark:border-slate-700 hover:border-indigo-400 dark:hover:border-indigo-500 hover:bg-indigo-50 dark:hover:bg-indigo-950 text-slate-700 dark:text-slate-200 font-medium transition-colors text-left"
                        >
                            {ROLE_LABELS[role] ?? role}
                        </button>
                    ))}
                </div>
                {toast && <Toast message={toast.message} variant={toast.variant} onDismiss={dismiss} />}
            </AuthLayout>
        )
    }

    return (
        <AuthLayout>
            <h1 className="text-2xl font-semibold text-slate-800 dark:text-slate-100 mb-4 text-center">
                Sign in to your account
            </h1>

            <form onSubmit={handleSubmit} noValidate className="flex flex-col gap-4">
                <Field label="Email" error={errors.email}>
                    <input
                        type="email"
                        value={email}
                        onChange={handleEmailChange}
                        placeholder="Email"
                        className={inputCls(!!errors.email)}
                    />
                </Field>

                <Field label="Password" error={errors.password}>
                    <input
                        type="password"
                        value={password}
                        onChange={handlePasswordChange}
                        placeholder="••••••••"
                        className={inputCls(!!errors.password)}
                    />
                </Field>

                <button type="submit" disabled={loading} className={primaryBtnCls}>
                    {loading ? 'Signing in…' : 'Sign in'}
                </button>
            </form>

            <p className="mt-5 text-center text-sm text-slate-500 dark:text-slate-400">
                Don&apos;t have an account?{' '}
                <Link to="/signup" className="text-indigo-500 hover:underline font-medium">
                    Sign up
                </Link>
            </p>
            {toast && <Toast message={toast.message} variant={toast.variant} onDismiss={dismiss} />}
        </AuthLayout>
    )
}
