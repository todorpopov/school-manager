import type { AuthResponse } from '../types/auth'

// const BASE = '/api'

export async function apiLogin(email: string, password: string): Promise<AuthResponse> {
    console.log(email, password);
    return {
        sessionId: 'mock-session-id',
        roles: ['ADMIN'],
        firstName: 'Ivan',
        lastName: 'Petrov',
        email,
    }

    // const res = await fetch(`${BASE}/auth/login`, {
    //   method: 'POST',
    //   headers: { 'Content-Type': 'application/json' },
    //   body: JSON.stringify({ email, password }),
    // })
    //
    // const json = await res.json()
    //
    // if (json.error) {
    //   throw new Error(json.message)
    // }
    //
    // return json.data as AuthResponse
}

export async function apiRegister(
    firstName: string,
    lastName: string,
    email: string,
    password: string,
): Promise<AuthResponse> {
    console.log(firstName, lastName, email, password);
    return {
        sessionId: 'mock-session-id',
        roles: ['TEACHER', 'PARENT'],
        firstName: 'Ivan',
        lastName: 'Petrov',
        email,
    }

    // const res = await fetch(`${BASE}/auth/register`, {
    //   method: 'POST',
    //   headers: { 'Content-Type': 'application/json' },
    //   body: JSON.stringify({ firstName, lastName, email, password }),
    // })
    //
    // const json = await res.json()
    //
    // if (json.error) {
    //   throw new Error(json.message)
    // }
    //
    // return json.data as AuthResponse
}
