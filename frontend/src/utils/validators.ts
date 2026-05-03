const LETTERS_ONLY = /^[A-Za-zÀ-ÖØ-öø-ÿ\s'-]+$/

export const validateName = (value: unknown): string | undefined => {
    const str = String(value ?? '').trim()
    if (!str) return 'This field is required'
    if (!LETTERS_ONLY.test(str)) return 'Only letters are allowed'
    if (str.length < 2) return 'Must be at least 2 characters'
    if (str.length > 255) return 'Must be at most 255 characters'
}

export const validateEmail = (value: unknown): string | undefined => {
    const str = String(value ?? '').trim()
    if (!str) return 'Email is required'
    if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(str)) return 'Enter a valid email address'
}

export const validatePassword = (value: unknown): string | undefined => {
    const str = String(value ?? '')
    if (!str) return 'Password is required'
    if (str.length < 8) return 'Password must be at least 8 characters'
}

export const validateText = (value: unknown): string | undefined => {
    const str = String(value ?? '').trim()
    if (!str) return 'This field is required'
    if (str.length > 500) return 'Must be at most 500 characters'
}

