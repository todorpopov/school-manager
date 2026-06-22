import { render } from '@testing-library/react'
import { describe, it, expect, vi } from 'vitest'
import { ConfirmDelete } from '../ResourceManager/ConfirmDelete'

describe('ConfirmDelete', () => {
    it('renders the confirmation dialog', () => {
        const { container } = render(
            <ConfirmDelete onConfirm={vi.fn()} onCancel={vi.fn()} />
        )
        expect(container).toMatchSnapshot()
    })
})

