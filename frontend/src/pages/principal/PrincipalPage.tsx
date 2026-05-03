import { useNavigate } from 'react-router-dom'
import type { IconType } from 'react-icons'
import { FaChalkboardTeacher, FaUserGraduate, FaUserFriends, FaBookOpen } from 'react-icons/fa'

const RESOURCES: { label: string; description: string; path: string; Icon: IconType }[] = [
    { label: 'Subjects', description: 'View all subjects', path: '/director/subjects', Icon: FaBookOpen },
    { label: 'Teachers', description: 'View teaching staff', path: '/director/teachers', Icon: FaChalkboardTeacher },
    { label: 'Students', description: 'View student records', path: '/director/students', Icon: FaUserGraduate },
    { label: 'Parents', description: 'View parent contacts', path: '/director/parents', Icon: FaUserFriends },
]

export default function PrincipalPage() {
    const navigate = useNavigate()

    return (
        <main className="max-w-5xl mx-auto px-4 py-10">
            <h1 className="text-2xl font-semibold text-slate-800 dark:text-slate-100 mb-8 text-center">Director Overview</h1>
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-2 gap-6">
                {RESOURCES.map(({ path, label, description, Icon }) => (
                    <button
                        key={path}
                        onClick={() => navigate(path)}
                        className="flex flex-col items-center gap-3 p-6 bg-white dark:bg-slate-800 rounded-2xl border border-slate-200 dark:border-slate-700 shadow-sm hover:shadow-md hover:border-indigo-300 dark:hover:border-indigo-600 transition-all text-center cursor-pointer"
                    >
                        <Icon className="text-3xl text-indigo-500 dark:text-indigo-400" />
                        <div>
                            <p className="font-semibold text-slate-800 dark:text-slate-100">{label}</p>
                            <p className="text-sm text-slate-500 dark:text-slate-400 mt-1">{description}</p>
                        </div>
                    </button>
                ))}
            </div>
        </main>
    )
}

