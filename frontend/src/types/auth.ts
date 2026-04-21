export type Role = 'ADMIN' | 'DIRECTOR' | 'TEACHER' | 'PARENT' | 'STUDENT' | 'USER'

export const AMBIGUOUS_ROLE_PAIRS: Role[][] = [
  ['TEACHER', 'PARENT'],
]

export interface AuthResponse {
  sessionId: string
  roles: Role[]
  firstName: string
  lastName: string
  email: string
}

export interface AuthState {
  sessionId: string
  activeRole: Role
  roles: Role[]
  firstName: string
  lastName: string
  email: string
}
