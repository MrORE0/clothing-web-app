import { describe, it, expect } from 'vitest'
import { render, screen, fireEvent } from '@testing-library/react'
import { MemoryRouter } from 'react-router'
import { FourCardsCarousel } from '../components/FourCardsCarousel'

const renderCarousel = () => render(<MemoryRouter><FourCardsCarousel /></MemoryRouter>)

describe('FourCardsCarousel', () => {
    describe('rendering', () => {
        it('renders without crashing', () => {
            expect(() => renderCarousel()).not.toThrow()
        })

        it('renders Woman and Man toggle buttons', () => {
            renderCarousel()
            expect(screen.getByText('Woman')).toBeInTheDocument()
            expect(screen.getByText('Man')).toBeInTheDocument()
        })

        it('renders left and right arrow buttons', () => {
            renderCarousel()
            expect(screen.getByAltText('arrow left icon')).toBeInTheDocument()
            expect(screen.getByAltText('arrow right icon')).toBeInTheDocument()
        })

        it('renders See More link', () => {
            renderCarousel()
            expect(screen.getByRole('link', { name: 'See More' })).toBeInTheDocument()
        })

        it('renders product cards with names and prices', () => {
            renderCarousel()
            const bikiniItems = screen.getAllByText('Bikini')
            expect(bikiniItems.length).toBeGreaterThan(0)
        })
    })

    describe('gender toggle', () => {
        it('Woman button is initially bold (active)', () => {
            renderCarousel()
            expect(screen.getByText('Woman')).toHaveClass('fw-bold')
        })

        it('Man button is initially not bold (inactive)', () => {
            renderCarousel()
            expect(screen.getByText('Man')).not.toHaveClass('fw-bold')
        })

        it('clicking Man makes it bold and Woman not bold', () => {
            renderCarousel()
            fireEvent.click(screen.getByText('Man'))
            expect(screen.getByText('Man')).toHaveClass('fw-bold')
            expect(screen.getByText('Woman')).not.toHaveClass('fw-bold')
        })

        it('clicking Man then Woman restores Woman as active', () => {
            renderCarousel()
            fireEvent.click(screen.getByText('Man'))
            fireEvent.click(screen.getByText('Woman'))
            expect(screen.getByText('Woman')).toHaveClass('fw-bold')
            expect(screen.getByText('Man')).not.toHaveClass('fw-bold')
        })

        it('switching to Man loads male products (Swim Trunks)', () => {
            renderCarousel()
            fireEvent.click(screen.getByText('Man'))
            expect(screen.getByText('Swim Trunks')).toBeInTheDocument()
        })

        it('switching back to Woman removes Swim Trunks', () => {
            renderCarousel()
            fireEvent.click(screen.getByText('Man'))
            fireEvent.click(screen.getByText('Woman'))
            expect(screen.queryByText('Swim Trunks')).not.toBeInTheDocument()
        })
    })

    describe('carousel navigation', () => {
        it('clicking right arrow does not throw', () => {
            renderCarousel()
            const rightContainer = screen.getByAltText('arrow right icon').closest('div.right-arrow-button-container')!
            expect(() => fireEvent.click(rightContainer.parentElement!)).not.toThrow()
        })

        it('clicking left arrow does not throw', () => {
            renderCarousel()
            const leftContainer = screen.getByAltText('arrow left icon').closest('div.left-arrow-button-container')!
            expect(() => fireEvent.click(leftContainer.parentElement!)).not.toThrow()
        })
    })
})
