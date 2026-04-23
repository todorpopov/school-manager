import { createContext } from 'react'
import type { Role, AuthResponse } from '../types/auth'

export interface AuthUser {
    sessionId: string
    token?: string
    activeRole: Role
    roles: Role[]
    firstName: string
    lastName: string
    email: string
}

export interface AuthContextType {
    user: AuthUser | null
    pendingAuth: AuthResponse | null
    login: (response: AuthResponse) => void
    selectRole: (role: Role) => void
    logout: () => void
    loading: boolean
}

export const AuthContext = createContext<AuthContextType | undefined>(undefined)
