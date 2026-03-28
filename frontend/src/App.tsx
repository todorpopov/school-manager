import { useState } from 'react'
import { ResourceManager } from './components/ResourceManager'
import type { FieldConfig } from './components/ResourceManager'

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
  { key: 'first_name', label: 'First name', type: 'text', required: true, placeholder: 'Ivan' },
  { key: 'last_name', label: 'Last name', type: 'text', required: true, placeholder: 'Petrov' },
  { key: 'email', label: 'Email', type: 'email', required: true, placeholder: 'ivan@school.bg' },
  { key: 'roles', label: 'Roles', type: 'multiselect', required: true, options: ROLE_OPTIONS },
]

const MOCK_USERS: User[] = [
  { user_id: 1, first_name: 'Maria', last_name: 'Georgieva', email: 'maria@school.bg', roles: ['TEACHER'] },
  { user_id: 2, first_name: 'Georgi', last_name: 'Stoyanov', email: 'georgi@school.bg', roles: ['DIRECTOR'] },
]

let nextId = 3

function App() {
  const [users, setUsers] = useState<User[]>(MOCK_USERS)

  const handleCreate = async (values: Partial<User>) => {
    const newUser: User = { ...(values as User), user_id: nextId++ }
    setUsers((prev) => [...prev, newUser])
  }

  const handleUpdate = async (id: User[keyof User], values: Partial<User>) => {
    setUsers((prev) =>
      prev.map((u) => (u.user_id === id ? { ...u, ...values } : u))
    )
  }

  const handleDelete = async (id: User[keyof User]) => {
    setUsers((prev) => prev.filter((u) => u.user_id !== id))
  }

  return (
    <div className="min-h-screen bg-slate-100 dark:bg-slate-950 py-10">
      <div className="max-w-4xl mx-auto px-4">
        <ResourceManager<User>
        title="Users"
        data={users}
        fields={USER_FIELDS}
        idKey="user_id"
        onCreate={handleCreate}
        onUpdate={handleUpdate}
        onDelete={handleDelete}
      />
      </div>
    </div>
  )
}

export default App
