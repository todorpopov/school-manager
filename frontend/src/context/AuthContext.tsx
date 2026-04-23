import { useState, useCallback } from 'react'
import { useNavigate } from 'react-router-dom'
import type { Role, AuthResponse } from '../types/auth'
import { AMBIGUOUS_ROLE_PAIRS } from '../types/auth'
import type { AuthUser } from './authContext.types'
import { AuthContext } from './authContext.types'

const USER_KEY = 'sm_user'

function needsRolePicker(roles: Role[]): boolean {
    return AMBIGUOUS_ROLE_PAIRS.some((pair) => pair.every((r) => roles.includes(r)))
}

function buildUser(authResponse: AuthResponse, activeRole: Role): AuthUser {
    return {
        sessionId: authResponse.sessionId,
        token: authResponse.token,
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
    const [loading] = useState(false)
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
        const newUser = buildUser(authResponse, authResponse.roles[0])
        saveUser(newUser)
        navigate('/dashboard')
    }, [navigate])

    const selectRole = useCallback((role: Role) => {
        if (!pendingAuth) return
        const newUser = buildUser(pendingAuth, role)
        saveUser(newUser)
        setPendingAuth(null)
        navigate('/dashboard')
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
