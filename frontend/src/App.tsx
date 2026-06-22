import { Routes, Route, Navigate } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { AuthProvider } from './context/AuthContext'
import ProtectedRoute from './components/ProtectedRoute'
import AppLayout from './components/AppLayout'
import LoginPage from './pages/LoginPage'
import SignupPage from './pages/SignupPage'
import AdminRegisterPage from './pages/AdminRegisterPage'
import HomePage from './pages/HomePage'
import StatisticsPage from './pages/admin/StatisticsPage'
import ManagementPage from './pages/admin/ManagementPage'
import SchoolManagementPage from './pages/admin/management/SchoolManagementPage'
import PrincipalManagementPage from './pages/admin/management/PrincipalManagementPage'
import ClassManagementPage from './pages/admin/management/ClassManagementPage'
import TermManagementPage from './pages/admin/management/TermManagementPage'
import SubjectManagementPage from './pages/admin/management/SubjectManagementPage'
import StudentsManagementPage from './pages/admin/management/StudentsManagementPage'
import TeachersManagementPage from './pages/admin/management/TeachersManagementPage'
import ParentsManagementPage from './pages/admin/management/ParentsManagementPage'
import CurriculumManagementPage from './pages/admin/management/CurriculumManagementPage'
import PrincipalPage from './pages/principal/PrincipalPage'
import PrincipalStatisticsPage from './pages/principal/PrincipalStatisticsPage'
import PrincipalSubjectsPage from './pages/principal/PrincipalSubjectsPage'
import PrincipalClassesPage from './pages/principal/PrincipalClassesPage'
import PrincipalTeachersPage from './pages/principal/PrincipalTeachersPage'
import PrincipalStudentsPage from './pages/principal/PrincipalStudentsPage'
import PrincipalParentsPage from './pages/principal/PrincipalParentsPage'
import ParentPage from './pages/parent/ParentPage'
import StudentPage from './pages/student/StudentPage'
import TeacherPage from './pages/teacher/TeacherPage.tsx';
import AccessDeniedPage from './pages/AccessDeniedPage'

const queryClient = new QueryClient()

function App() {
    return (
        <QueryClientProvider client={queryClient}>
            <AuthProvider>
                <Routes>
                    <Route path="/login" element={<LoginPage />} />
                    <Route path="/signup" element={<SignupPage />} />
                    <Route path="/admin/register" element={<AdminRegisterPage />} />
                    <Route path="/access-denied" element={<AccessDeniedPage />} />
                    <Route
                        element={
                            <ProtectedRoute>
                                <AppLayout />
                            </ProtectedRoute>
                        }
                    >
                        <Route path="/home" element={<HomePage />} />
                        <Route path="/admin/statistics" element={
                            <ProtectedRoute allowedRoles={['ADMIN']}>
                                <StatisticsPage />
                            </ProtectedRoute>
                        } />
                        <Route path="/admin/management" element={
                            <ProtectedRoute allowedRoles={['ADMIN']}>
                                <ManagementPage />
                            </ProtectedRoute>
                        } />
                        <Route path="/admin/management/school" element={
                            <ProtectedRoute allowedRoles={['ADMIN']}>
                                <SchoolManagementPage />
                            </ProtectedRoute>
                        } />
                        <Route path="/admin/management/principal" element={
                            <ProtectedRoute allowedRoles={['ADMIN']}>
                                <PrincipalManagementPage />
                            </ProtectedRoute>
                        } />
                        <Route path="/admin/management/classes" element={
                            <ProtectedRoute allowedRoles={['ADMIN']}>
                                <ClassManagementPage />
                            </ProtectedRoute>
                        } />
                        <Route path="/admin/management/terms" element={
                            <ProtectedRoute allowedRoles={['ADMIN']}>
                                <TermManagementPage />
                            </ProtectedRoute>
                        } />
                        <Route path="/admin/management/subjects" element={
                            <ProtectedRoute allowedRoles={['ADMIN']}>
                                <SubjectManagementPage />
                            </ProtectedRoute>
                        } />
                        <Route path="/admin/management/students" element={
                            <ProtectedRoute allowedRoles={['ADMIN']}>
                                <StudentsManagementPage />
                            </ProtectedRoute>
                        } />
                        <Route path="/admin/management/teachers" element={
                            <ProtectedRoute allowedRoles={['ADMIN']}>
                                <TeachersManagementPage />
                            </ProtectedRoute>
                        } />
                        <Route path="/admin/management/parents" element={
                            <ProtectedRoute allowedRoles={['ADMIN']}>
                                <ParentsManagementPage />
                            </ProtectedRoute>
                        } />
                        <Route path="/admin/management/curriculum" element={
                            <ProtectedRoute allowedRoles={['ADMIN']}>
                                <CurriculumManagementPage />
                            </ProtectedRoute>
                        } />
                        <Route path="/principal" element={
                            <ProtectedRoute allowedRoles={['DIRECTOR']}>
                                <PrincipalPage />
                            </ProtectedRoute>
                        } />
                        <Route path="/principal/statistics" element={
                            <ProtectedRoute allowedRoles={['DIRECTOR']}>
                                <PrincipalStatisticsPage />
                            </ProtectedRoute>
                        } />
                        <Route path="/principal/subjects" element={
                            <ProtectedRoute allowedRoles={['DIRECTOR']}>
                                <PrincipalSubjectsPage />
                            </ProtectedRoute>
                        } />
                        <Route path="/principal/classes" element={
                            <ProtectedRoute allowedRoles={['DIRECTOR']}>
                                <PrincipalClassesPage />
                            </ProtectedRoute>
                        } />
                        <Route path="/principal/teachers" element={
                            <ProtectedRoute allowedRoles={['DIRECTOR']}>
                                <PrincipalTeachersPage />
                            </ProtectedRoute>
                        } />
                        <Route path="/principal/students" element={
                            <ProtectedRoute allowedRoles={['DIRECTOR']}>
                                <PrincipalStudentsPage />
                            </ProtectedRoute>
                        } />
                        <Route path="/principal/parents" element={
                            <ProtectedRoute allowedRoles={['DIRECTOR']}>
                                <PrincipalParentsPage />
                            </ProtectedRoute>
                        } />
                        <Route path="/parent" element={
                            <ProtectedRoute allowedRoles={['PARENT']}>
                                <ParentPage />
                            </ProtectedRoute>
                        } />
                        <Route path="/student" element={
                            <ProtectedRoute allowedRoles={['STUDENT']}>
                                <StudentPage />
                            </ProtectedRoute>
                        } />
                        <Route path="/teacher" element={
                            <ProtectedRoute allowedRoles={['TEACHER']}>
                                <TeacherPage />
                            </ProtectedRoute>
                        } />
                    </Route>
                    <Route path="*" element={<Navigate to="/login" replace />} />
                </Routes>
            </AuthProvider>
        </QueryClientProvider>
    )
}

export default App
