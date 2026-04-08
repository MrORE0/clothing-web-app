import { describe, it, expect, beforeEach } from 'vitest'
import { render, screen, fireEvent, act } from '@testing-library/react'
import { MemoryRouter } from 'react-router'
import { NavBar } from '../partials/NavBar'

const renderNavBar = () => render(<MemoryRouter><NavBar /></MemoryRouter>)

describe('NavBar', () => {
    describe('rendering', () => {
        it('renders the Browse home link', () => {
            renderNavBar()
            expect(screen.getByText('Browse')).toBeInTheDocument()
        })

        it('renders Him and Her nav items', () => {
            renderNavBar()
            expect(screen.getByText('Him')).toBeInTheDocument()
            expect(screen.getByText('Her')).toBeInTheDocument()
        })

        it('renders Hot nav item', () => {
            renderNavBar()
            expect(screen.getByText('Hot')).toBeInTheDocument()
        })

        it('renders Custom nav link pointing to /customize', () => {
            renderNavBar()
            const link = screen.getByRole('link', { name: 'Custom' })
            expect(link).toHaveAttribute('href', '/customize')
        })

        it('renders search input with placeholder Explore', () => {
            renderNavBar()
            expect(screen.getByPlaceholderText('Explore')).toBeInTheDocument()
        })

        it('renders Find button', () => {
            renderNavBar()
            expect(screen.getByRole('button', { name: 'Find' })).toBeInTheDocument()
        })

        it('renders the brand name AssK', () => {
            renderNavBar()
            expect(screen.getByText('AssK')).toBeInTheDocument()
        })

        it('renders For Him and For Her hot-section labels', () => {
            renderNavBar()
            expect(screen.getByText('For Him')).toBeInTheDocument()
            expect(screen.getByText('For Her')).toBeInTheDocument()
        })

        it('renders the navbar with initial height of 110px', () => {
            renderNavBar()
            const nav = document.querySelector('nav.navbar') as HTMLElement
            expect(nav).toHaveStyle({ height: '110px' })
        })
    })

    describe('scroll behaviour', () => {
        beforeEach(() => {
            Object.defineProperty(window, 'scrollY', { value: 0, configurable: true, writable: true })
        })

        it('hides the navbar when scrolling down', () => {
            renderNavBar()
            const nav = document.querySelector('nav.navbar') as HTMLElement

            act(() => {
                Object.defineProperty(window, 'scrollY', { value: 300, configurable: true, writable: true })
                window.dispatchEvent(new Event('scroll'))
            })

            expect(nav.style.opacity).toBe('0')
            expect(nav.style.height).toBe('0px')
        })

        it('shows the navbar again when scrolling up', () => {
            renderNavBar()
            const nav = document.querySelector('nav.navbar') as HTMLElement

            act(() => {
                Object.defineProperty(window, 'scrollY', { value: 300, configurable: true, writable: true })
                window.dispatchEvent(new Event('scroll'))
            })
            act(() => {
                Object.defineProperty(window, 'scrollY', { value: 100, configurable: true, writable: true })
                window.dispatchEvent(new Event('scroll'))
            })

            expect(nav.style.opacity).toBe('1')
            expect(nav.style.height).toBe('110px')
        })
    })

    describe('dropdown content', () => {
        it('updates clothes list to man types on Him hover', () => {
            renderNavBar()
            fireEvent.mouseEnter(screen.getByText('Him').closest('li')!)
            // Him dropdown renders Clothes/Shoes/Accessories headings
            const clothesHeadings = screen.getAllByText('Clothes')
            expect(clothesHeadings.length).toBeGreaterThan(0)
        })

        it('updates clothes list to woman types on Her hover', () => {
            renderNavBar()
            fireEvent.mouseEnter(screen.getByText('Her').closest('li')!)
            const clothesHeadings = screen.getAllByText('Clothes')
            expect(clothesHeadings.length).toBeGreaterThan(0)
        })
    })
})
