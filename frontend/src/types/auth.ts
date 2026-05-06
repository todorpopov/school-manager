export type Role = 'ADMIN' | 'DIRECTOR' | 'TEACHER' | 'PARENT' | 'STUDENT' | 'USER'

export const AMBIGUOUS_ROLE_PAIRS: Role[][] = [
    ['TEACHER', 'PARENT'],
]

export interface AuthResponse {
    sessionId: string
    token?: string
    userId: number
    roles: Role[]
    firstName: string
    lastName: string
    email: string
}

export interface AuthState {
    sessionId: string
    activeRole: Role
    userId: number
    roles: Role[]
    firstName: string
    lastName: string
    email: string
}
