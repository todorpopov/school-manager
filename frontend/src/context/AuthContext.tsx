import { useState, useCallback } from 'react'
import { useNavigate } from 'react-router-dom'
import type { Role, AuthResponse } from '../types/auth'
import type { AuthUser } from './authContext.types'
import { AuthContext } from './authContext.types'
import { apiSelectRole } from '../api/auth'

const USER_KEY = 'sm_user'

function needsRolePicker(roles: Role[]): boolean {
    // Filter out the generic USER role as it's not a specific role
    const specificRoles = roles.filter(r => r !== 'USER')
    return specificRoles.length > 1
}

function buildUser(authResponse: AuthResponse, activeRole: Role): AuthUser {
    return {
        sessionId: authResponse.sessionId,
        token: authResponse.token,
        userId: authResponse.userId,
        activeRole,
        roles: authResponse.roles,
        firstName: authResponse.firstName,
        lastName: authResponse.lastName,
        email: authResponse.email,
    }
}

export function AuthProvider({ children }: { children: React.ReactNode }) {
    const [user, setUser] = useState<AuthUser | null>(() => {
        const stored = localStorage.getItem(USER_KEY)
        return stored ? (JSON.parse(stored) as AuthUser) : null
    })
    const [pendingAuth, setPendingAuth] = useState<AuthResponse | null>(null)
    const [loading, setLoading] = useState(false)
    const navigate = useNavigate()

    const saveUser = (newUser: AuthUser) => {
        setUser(newUser)
        localStorage.setItem(USER_KEY, JSON.stringify(newUser))
        if (newUser.token) {
            localStorage.setItem('token', newUser.token)
        }
    }

    const login = useCallback((authResponse: AuthResponse) => {
        if (needsRolePicker(authResponse.roles)) {
            setPendingAuth(authResponse)
            return
        }
        const specificRoles = authResponse.roles.filter(r => r !== 'USER')
        const newUser = buildUser(authResponse, specificRoles[0])
        saveUser(newUser)
        navigate('/home')
    }, [navigate])

    const selectRole = useCallback(async (role: Role) => {
        if (!pendingAuth) return

        setLoading(true)
        try {
            await apiSelectRole(pendingAuth.sessionId, role)

            const newUser = buildUser(pendingAuth, role)
            saveUser(newUser)
            setPendingAuth(null)
            navigate('/home')
        } catch (error) {
            console.error('Failed to set session role:', error)
            const newUser = buildUser(pendingAuth, role)
            saveUser(newUser)
            setPendingAuth(null)
            navigate('/home')
        } finally {
            setLoading(false)
        }
    }, [pendingAuth, navigate])

    const logout = useCallback(() => {
        setUser(null)
        setPendingAuth(null)
        localStorage.removeItem(USER_KEY)
        localStorage.removeItem('token')
        navigate('/login')
    }, [navigate])

    return (
        <AuthContext.Provider value={{ user, pendingAuth, login, selectRole, logout, loading }}>
            {children}
        </AuthContext.Provider>
    )
}
