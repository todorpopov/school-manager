import type React from 'react'
import { Navigate } from 'react-router-dom'
import { useAuth } from '../hooks/useAuth'
import type { Role } from '../types/auth'

interface ProtectedRouteProps {
    children: React.ReactNode
    allowedRoles?: Role[]
}

const ProtectedRoute: React.FC<ProtectedRouteProps> = ({ children, allowedRoles }) => {
    const { user, loading } = useAuth()

    if (loading) return null

    if (!user) {
        return <Navigate to="/login" replace />
    }

    if (allowedRoles && !allowedRoles.includes(user.activeRole)) {
        return <Navigate to="/access-denied" replace />
    }

    return <>{children}</>
}

export default ProtectedRoute
