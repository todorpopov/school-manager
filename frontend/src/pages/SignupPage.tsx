import React, { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'
import { apiRegister } from '../api/auth'
import { useToast } from '../hooks/useToast'
import { Toast } from '../components/Toast'
import { Field } from '../components/Field'
import { AuthLayout } from '../components/AuthLayout'
import { inputCls, primaryBtnCls } from '../styles/formStyles'

export default function SignupPage() {
  const { login } = useAuth()
  const navigate = useNavigate()

  const { toast, show: showToast, dismiss } = useToast()

  const [firstName, setFirstName] = useState('')
  const [lastName, setLastName] = useState('')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [confirm, setConfirm] = useState('')
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
    }

    if (!confirm) {
      e.confirm = 'Please confirm your password'
    } else if (confirm !== password) {
      e.confirm = 'Passwords do not match'
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
      const response = await apiRegister(firstName, lastName, email, password)
      login(response)
      navigate('/dashboard')
    } catch (err) {
      showToast(err instanceof Error ? err.message : 'Registration failed')
    } finally {
      setLoading(false)
    }
  }

  return (
    <AuthLayout>
      <h1 className="text-2xl font-semibold text-slate-800 dark:text-slate-100 mb-4 text-center">
        Create an account
      </h1>

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

        <Field label="Email" error={errors.email}>
          <input
            type="email"
            value={email}
            onChange={(e) => { setEmail(e.target.value); clearError('email') }}
            placeholder="Email"
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

        <Field label="Confirm password" error={errors.confirm}>
          <input
            type="password"
            value={confirm}
            onChange={(e) => { setConfirm(e.target.value); clearError('confirm') }}
            placeholder="••••••••"
            className={inputCls(!!errors.confirm)}
          />
        </Field>

        <button type="submit" disabled={loading} className={primaryBtnCls}>
          {loading ? 'Creating account…' : 'Create account'}
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
