import { useNavigate } from 'react-router-dom'
import type { IconType } from 'react-icons'
import { FaSchool, FaUserTie, FaUserGraduate, FaChalkboardTeacher, FaUserFriends, FaBookOpen } from 'react-icons/fa'

const RESOURCES: { label: string; description: string; path: string; Icon: IconType }[] = [
    { label: 'School', description: 'Manage school information and settings', path: '/admin/management/school', Icon: FaSchool },
    { label: 'Principal', description: 'Manage school principals', path: '/admin/management/principal', Icon: FaUserTie },
    { label: 'Students', description: 'Manage student records', path: '/admin/management/students', Icon: FaUserGraduate },
    { label: 'Teachers', description: 'Manage teaching staff', path: '/admin/management/teachers', Icon: FaChalkboardTeacher },
    { label: 'Parents', description: 'Manage parent contacts', path: '/admin/management/parents', Icon: FaUserFriends },
    { label: 'Curriculum', description: 'Manage subjects and curriculum', path: '/admin/management/curriculum', Icon: FaBookOpen },
]

export default function ManagementPage() {
    const navigate = useNavigate()

    return (
        <main className="max-w-5xl mx-auto px-4 py-10">
            <h1 className="text-2xl font-semibold text-slate-800 dark:text-slate-100 mb-8 text-center">Management</h1>
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
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
