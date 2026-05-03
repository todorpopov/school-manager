import { Routes, Route, Navigate } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { AuthProvider } from './context/AuthContext'
import ProtectedRoute from './components/ProtectedRoute'
import AppLayout from './components/AppLayout'
import LoginPage from './pages/LoginPage'
import SignupPage from './pages/SignupPage'
import DashboardPage from './pages/DashboardPage'
import StatisticsPage from './pages/admin/StatisticsPage'
import ManagementPage from './pages/admin/ManagementPage'
import SchoolManagementPage from './pages/admin/management/SchoolManagementPage'
import PrincipalManagementPage from './pages/admin/management/PrincipalManagementPage'
import StudentsManagementPage from './pages/admin/management/StudentsManagementPage'
import TeachersManagementPage from './pages/admin/management/TeachersManagementPage'
import ParentsManagementPage from './pages/admin/management/ParentsManagementPage'
import CurriculumManagementPage from './pages/admin/management/CurriculumManagementPage'
import PrincipalPage from './pages/principal/PrincipalPage'
import PrincipalSubjectsPage from './pages/principal/PrincipalSubjectsPage'
import PrincipalTeachersPage from './pages/principal/PrincipalTeachersPage'
import PrincipalStudentsPage from './pages/principal/PrincipalStudentsPage'
import PrincipalParentsPage from './pages/principal/PrincipalParentsPage'
import ParentPage from './pages/parent/ParentPage'
import TeacherPage from './pages/teacher/TeacherPage.tsx';

const queryClient = new QueryClient()

function App() {
    return (
        <QueryClientProvider client={queryClient}>
            <AuthProvider>
                <Routes>
                    <Route path="/login" element={<LoginPage />} />
                    <Route path="/signup" element={<SignupPage />} />
                    <Route
                        element={
                            <ProtectedRoute>
                                <AppLayout />
                            </ProtectedRoute>
                        }
                    >
                        <Route path="/dashboard" element={<DashboardPage />} />
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
                        <Route path="/director" element={
                            <ProtectedRoute allowedRoles={['DIRECTOR']}>
                                <PrincipalPage />
                            </ProtectedRoute>
                        } />
                        <Route path="/director/subjects" element={
                            <ProtectedRoute allowedRoles={['DIRECTOR']}>
                                <PrincipalSubjectsPage />
                            </ProtectedRoute>
                        } />
                        <Route path="/director/teachers" element={
                            <ProtectedRoute allowedRoles={['DIRECTOR']}>
                                <PrincipalTeachersPage />
                            </ProtectedRoute>
                        } />
                        <Route path="/director/students" element={
                            <ProtectedRoute allowedRoles={['DIRECTOR']}>
                                <PrincipalStudentsPage />
                            </ProtectedRoute>
                        } />
                        <Route path="/director/parents" element={
                            <ProtectedRoute allowedRoles={['DIRECTOR']}>
                                <PrincipalParentsPage />
                            </ProtectedRoute>
                        } />
                        <Route path="/parent" element={
                            <ProtectedRoute allowedRoles={['PARENT']}>
                                <ParentPage />
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
