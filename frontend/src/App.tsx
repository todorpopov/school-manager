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
                    </Route>
                    <Route path="*" element={<Navigate to="/login" replace />} />
                </Routes>
            </AuthProvider>
        </QueryClientProvider>
    )
}

export default App
