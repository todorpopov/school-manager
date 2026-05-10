import type React from 'react';
import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { useAuth } from '../hooks/useAuth'
import { apiRegisterAdmin } from '../api/auth'
import { useToast } from '../hooks/useToast'
import { Toast } from '../components/Toast'
import { Field } from '../components/Field'
import { AuthLayout } from '../components/AuthLayout'
import { inputCls, primaryBtnCls } from '../styles/formStyles'

export default function AdminRegisterPage() {
    const { login } = useAuth()
    const navigate = useNavigate()

    const { toast, show: showToast, dismiss } = useToast()

    const [firstName, setFirstName] = useState('')
    const [lastName, setLastName] = useState('')
    const [email, setEmail] = useState('')
    const [password, setPassword] = useState('')
    const [systemAuthToken, setSystemAuthToken] = useState('')
    const [errors, setErrors] = useState<Record<string, string>>({})
    const [loading, setLoading] = useState(false)

    const clearError = (key: string) => {
        setErrors((prev) => {
            const next = { ...prev }
            delete next[key]
            return next
        })
    }

    const validate = () => {
        const e: Record<string, string> = {}

        if (!firstName) {
            e.firstName = 'First name is required'
        }

        if (!lastName) {
            e.lastName = 'Last name is required'
        }

        if (!email) {
            e.email = 'Email is required'
        } else if (!email.includes('@')) {
            e.email = 'Email must contain @'
        }

        if (!password) {
            e.password = 'Password is required'
        } else if (password.length < 8) {
            e.password = 'Password must be at least 8 characters'
        }

        if (!systemAuthToken) {
            e.systemAuthToken = 'System authorization token is required'
        }

        setErrors(e)
        return Object.keys(e).length === 0
    }

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault()

        if (!validate()) {
            return
        }

        setLoading(true)

        try {
            const response = await apiRegisterAdmin(firstName, lastName, email, password, systemAuthToken)
            login(response)
            navigate('/home')
        } catch (err) {
            showToast(err instanceof Error ? err.message : 'Admin registration failed')
        } finally {
            setLoading(false)
        }
    }

    return (
        <AuthLayout>
            <div className="mb-6 text-center">
                <div className="inline-flex items-center justify-center w-12 h-12 rounded-full bg-indigo-100 dark:bg-indigo-900 mb-3">
                    <svg className="w-6 h-6 text-indigo-600 dark:text-indigo-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
                    </svg>
                </div>
                <h1 className="text-2xl font-semibold text-slate-800 dark:text-slate-100">
                    Create Admin Account
                </h1>
                <p className="text-sm text-slate-500 dark:text-slate-400 mt-2">
                    Register as a system administrator
                </p>
            </div>

            <form onSubmit={handleSubmit} noValidate className="flex flex-col gap-4">
                <div className="grid grid-cols-2 gap-3">
                    <Field label="First name" error={errors.firstName}>
                        <input
                            type="text"
                            value={firstName}
                            onChange={(e) => { setFirstName(e.target.value); clearError('firstName') }}
                            placeholder="First name"
                            className={inputCls(!!errors.firstName)}
                        />
                    </Field>
                    <Field label="Last name" error={errors.lastName}>
                        <input
                            type="text"
                            value={lastName}
                            onChange={(e) => { setLastName(e.target.value); clearError('lastName') }}
                            placeholder="Last name"
                            className={inputCls(!!errors.lastName)}
                        />
                    </Field>
                </div>

                <Field label="Admin Email" error={errors.email}>
                    <input
                        type="email"
                        value={email}
                        onChange={(e) => { setEmail(e.target.value); clearError('email') }}
                        placeholder="admin@school.com"
                        className={inputCls(!!errors.email)}
                    />
                </Field>

                <Field label="Password" error={errors.password}>
                    <input
                        type="password"
                        value={password}
                        onChange={(e) => { setPassword(e.target.value); clearError('password') }}
                        placeholder="••••••••"
                        className={inputCls(!!errors.password)}
                    />
                </Field>

                <Field label="System Authorization Token" error={errors.systemAuthToken}>
                    <input
                        type="password"
                        value={systemAuthToken}
                        onChange={(e) => { setSystemAuthToken(e.target.value); clearError('systemAuthToken') }}
                        placeholder="Enter system auth token"
                        className={inputCls(!!errors.systemAuthToken)}
                    />
                </Field>

                <button type="submit" disabled={loading} className={primaryBtnCls}>
                    {loading ? 'Creating admin account…' : 'Create admin account'}
                </button>
            </form>

            <p className="mt-5 text-center text-sm text-slate-500 dark:text-slate-400">
                Already have an account?{' '}
                <Link to="/login" className="text-indigo-500 hover:underline font-medium">
                    Sign in
                </Link>
            </p>
            {toast && <Toast message={toast.message} variant={toast.variant} onDismiss={dismiss} />}
        </AuthLayout>
    )
}

