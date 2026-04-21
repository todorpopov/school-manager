import { createContext, useContext, useState, useCallback } from 'react'
import type { AuthState, AuthResponse, Role } from '../types/auth'
import { AMBIGUOUS_ROLE_PAIRS } from '../types/auth'

interface AuthContextValue {
  auth: AuthState | null
  pendingAuth: AuthResponse | null
  login: (response: AuthResponse) => void
  selectRole: (role: Role) => void
  logout: () => void
}

const AuthContext = createContext<AuthContextValue | null>(null)

const SESSION_KEY = 'sm_auth'

function loadAuth(): AuthState | null {
  try {
    const raw = sessionStorage.getItem(SESSION_KEY)
    return raw ? (JSON.parse(raw) as AuthState) : null
  } catch {
    return null
  }
}

function saveAuth(state: AuthState) {
  sessionStorage.setItem(SESSION_KEY, JSON.stringify(state))
}

function clearAuth() {
  sessionStorage.removeItem(SESSION_KEY)
}

function needsRolePicker(roles: Role[]): boolean {
  return AMBIGUOUS_ROLE_PAIRS.some(
    (pair) => pair.every((r) => roles.includes(r))
  )
}

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [auth, setAuth] = useState<AuthState | null>(loadAuth)
  const [pendingAuth, setPendingAuth] = useState<AuthResponse | null>(null)

  const login = useCallback((response: AuthResponse) => {
    if (needsRolePicker(response.roles)) {
      setPendingAuth(response)
      return
    }

    const state: AuthState = {
      sessionId: response.sessionId,
      activeRole: response.roles[0],
      roles: response.roles,
      firstName: response.firstName,
      lastName: response.lastName,
      email: response.email,
    }
    saveAuth(state)
    setAuth(state)
  }, [])

  const selectRole = useCallback((role: Role) => {
    if (!pendingAuth) {
      return
    }

    const state: AuthState = {
      sessionId: pendingAuth.sessionId,
      activeRole: role,
      roles: pendingAuth.roles,
      firstName: pendingAuth.firstName,
      lastName: pendingAuth.lastName,
      email: pendingAuth.email,
    }
    saveAuth(state)
    setAuth(state)
    setPendingAuth(null)
  }, [pendingAuth])

  const logout = useCallback(() => {
    clearAuth()
    setAuth(null)
    setPendingAuth(null)
  }, [])

  return (
    <AuthContext.Provider value={{ auth, pendingAuth, login, selectRole, logout }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error('useAuth must be used inside AuthProvider')
  return ctx
}

