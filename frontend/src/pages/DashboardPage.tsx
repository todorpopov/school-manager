import { useNavigate } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'
import { ResourceManager } from '../components/ResourceManager'
import type { FieldConfig } from '../components/ResourceManager'
import { useState } from 'react'

interface User {
  [key: string]: unknown
  user_id: number
  first_name: string
  last_name: string
  email: string
  roles: string[]
}

const ROLE_OPTIONS = [
  { label: 'Admin', value: 'ADMIN' },
  { label: 'Director', value: 'DIRECTOR' },
  { label: 'Teacher', value: 'TEACHER' },
  { label: 'Parent', value: 'PARENT' },
  { label: 'Student', value: 'STUDENT' },
]

const USER_FIELDS: FieldConfig<User>[] = [
  { key: 'user_id', label: 'ID', type: 'number', hideInForm: true },
  { key: 'first_name', label: 'First name', type: 'text', required: true, placeholder: 'First name' },
  { key: 'last_name', label: 'Last name', type: 'text', required: true, placeholder: 'Last name' },
  { key: 'email', label: 'Email', type: 'email', required: true, placeholder: 'Email' },
  { key: 'roles', label: 'Roles', type: 'multiselect', required: true, options: ROLE_OPTIONS },
]

const MOCK_USERS: User[] = [
  { user_id: 1, first_name: 'Maria', last_name: 'Georgieva', email: 'maria@school.bg', roles: ['TEACHER'] },
  { user_id: 2, first_name: 'Georgi', last_name: 'Stoyanov', email: 'georgi@school.bg', roles: ['DIRECTOR'] },
]

let nextId = 3

const ROLE_LABELS: Record<string, string> = {
  ADMIN: 'Administrator',
  DIRECTOR: 'Director',
  TEACHER: 'Teacher',
  PARENT: 'Parent',
  STUDENT: 'Student',
  USER: 'User',
}

export default function DashboardPage() {
  const { auth, logout } = useAuth()
  const navigate = useNavigate()
  const [users, setUsers] = useState<User[]>(MOCK_USERS)

  const handleLogout = () => { logout(); navigate('/login') }

  const handleCreate = async (values: Partial<User>) => {
    setUsers((prev) => [...prev, { ...(values as User), user_id: nextId++ }])
  }

  const handleUpdate = async (id: User[keyof User], values: Partial<User>) => {
    setUsers((prev) => prev.map((u) => (u.user_id === id ? { ...u, ...values } : u)))
  }

  const handleDelete = async (id: User[keyof User]) => {
    setUsers((prev) => prev.filter((u) => u.user_id !== id))
  }

  return (
    <div className="min-h-screen bg-slate-100 dark:bg-slate-950">
      <header className="bg-white dark:bg-slate-800 border-b border-slate-200 dark:border-slate-700 px-6 py-3 flex items-center justify-between">
        <span className="font-semibold text-slate-800 dark:text-slate-100">School Manager</span>
        <div className="flex items-center gap-4">
          <span className="text-sm text-slate-500 dark:text-slate-400">
            {auth?.firstName} {auth?.lastName}
            <span className="ml-2 px-2 py-0.5 bg-indigo-100 dark:bg-indigo-900 text-indigo-700 dark:text-indigo-300 rounded-full text-xs font-medium">
              {ROLE_LABELS[auth?.activeRole ?? ''] ?? auth?.activeRole}
            </span>
          </span>
          <button
            onClick={handleLogout}
            className="text-sm text-slate-500 hover:text-slate-800 dark:hover:text-slate-100 transition-colors cursor-pointer bg-transparent border-none"
          >
            Sign out
          </button>
        </div>
      </header>

      <main className="max-w-4xl mx-auto px-4 py-10">
        <ResourceManager<User>
          title="Users"
          data={users}
          fields={USER_FIELDS}
          idKey="user_id"
          onCreate={handleCreate}
          onUpdate={handleUpdate}
          onDelete={handleDelete}
        />
      </main>
    </div>
  )
}
