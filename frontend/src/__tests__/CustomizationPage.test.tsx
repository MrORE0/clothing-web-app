import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import { MemoryRouter } from 'react-router'
import { CustomizationPage } from '../pages/CustomizationPage'

const renderPage = () => render(<MemoryRouter><CustomizationPage /></MemoryRouter>)

describe('CustomizationPage', () => {
    describe('rendering', () => {
        it('renders without crashing', () => {
            expect(() => renderPage()).not.toThrow()
        })

        it('renders the Customize heading', () => {
            renderPage()
            expect(screen.getByText('Customize')).toBeInTheDocument()
        })

        it('renders the intro text', () => {
            renderPage()
            expect(screen.getByText(/Take picking your style to another level/)).toBeInTheDocument()
        })
    })

    describe('Upper Body section', () => {
        it('renders Upper Body heading', () => {
            renderPage()
            expect(screen.getByText('Upper Body')).toBeInTheDocument()
        })

        it('renders Neck Circumference input', () => {
            renderPage()
            expect(screen.getByLabelText('Neck Circumference')).toBeInTheDocument()
        })

        it('renders Chest Circumference input', () => {
            renderPage()
            expect(screen.getByLabelText('Chest Circumference')).toBeInTheDocument()
        })

        it('renders Tors Length input', () => {
            renderPage()
            expect(screen.getByLabelText('Tors Length')).toBeInTheDocument()
        })
    })

    describe('Lower Body section', () => {
        it('renders Lower Body heading', () => {
            renderPage()
            expect(screen.getByText('Lower Body')).toBeInTheDocument()
        })

        it('renders Waist Circumference input', () => {
            renderPage()
            expect(screen.getByLabelText('Waist Circumference')).toBeInTheDocument()
        })

        it('renders Leg Length input', () => {
            renderPage()
            expect(screen.getByLabelText('Leg Length')).toBeInTheDocument()
        })

        it('renders Leg Width input', () => {
            renderPage()
            expect(screen.getByLabelText('Leg Width')).toBeInTheDocument()
        })
    })

    describe('Additional section', () => {
        it('renders Additional heading', () => {
            renderPage()
            expect(screen.getByText('Additional')).toBeInTheDocument()
        })

        it('renders Clothing Color input', () => {
            renderPage()
            expect(screen.getByLabelText('Clothing Color')).toBeInTheDocument()
        })

        it('renders Clothing material input', () => {
            renderPage()
            expect(screen.getByLabelText('Clothing material')).toBeInTheDocument()
        })

        it('renders Clothing Style input', () => {
            renderPage()
            expect(screen.getByLabelText('Clothing Style')).toBeInTheDocument()
        })
    })

    describe('icons', () => {
        it('renders Shirt Icon alt text', () => {
            renderPage()
            expect(screen.getByAltText('Shirt Icon')).toBeInTheDocument()
        })

        it('renders Pants Icon alt text', () => {
            renderPage()
            expect(screen.getByAltText('Pants Icon')).toBeInTheDocument()
        })
    })
})
