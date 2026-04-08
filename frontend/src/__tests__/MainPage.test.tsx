import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import { MemoryRouter } from 'react-router'
import MainPage from '../pages/MainPage'

const renderPage = () => render(<MemoryRouter><MainPage /></MemoryRouter>)

describe('MainPage', () => {
    it('renders without crashing', () => {
        expect(() => renderPage()).not.toThrow()
    })

    it('renders the Seasonal section heading', () => {
        renderPage()
        expect(screen.getByText('Seasonal')).toBeInTheDocument()
    })

    it('renders the navbar Browse link', () => {
        renderPage()
        expect(screen.getByText('Browse')).toBeInTheDocument()
    })

    it('renders the FourCardsCarousel gender toggle', () => {
        renderPage()
        expect(screen.getByText('Woman')).toBeInTheDocument()
        expect(screen.getByText('Man')).toBeInTheDocument()
    })

    it('renders the hero banner text', () => {
        renderPage()
        expect(screen.getByText(/A way for you to dress like/)).toBeInTheDocument()
    })
})
