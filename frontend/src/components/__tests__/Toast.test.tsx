import { render } from '@testing-library/react'
import { describe, it, expect, vi } from 'vitest'
import { Toast } from '../Toast'

describe('Toast', () => {
    it('renders error variant', () => {
        const { container } = render(
            <Toast message="Something went wrong" variant="error" onDismiss={vi.fn()} />
        )
        expect(container).toMatchSnapshot()
    })

    it('renders success variant', () => {
        const { container } = render(
            <Toast message="Saved successfully" variant="success" onDismiss={vi.fn()} />
        )
        expect(container).toMatchSnapshot()
    })

    it('defaults to error variant', () => {
        const { container } = render(
            <Toast message="Oops" onDismiss={vi.fn()} />
        )
        expect(container).toMatchSnapshot()
    })
})

