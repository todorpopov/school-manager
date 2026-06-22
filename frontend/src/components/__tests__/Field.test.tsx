import { render } from '@testing-library/react'
import { describe, it, expect } from 'vitest'
import { Field } from '../Field'

describe('Field', () => {
    it('renders label and children', () => {
        const { container } = render(
            <Field label="Email">
                <input type="email" />
            </Field>
        )
        expect(container).toMatchSnapshot()
    })

    it('renders required indicator', () => {
        const { container } = render(
            <Field label="Password" required>
                <input type="password" />
            </Field>
        )
        expect(container).toMatchSnapshot()
    })

    it('renders error message', () => {
        const { container } = render(
            <Field label="Email" error="Invalid email address">
                <input type="email" />
            </Field>
        )
        expect(container).toMatchSnapshot()
    })
})

