import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import { createMemoryRouter, RouterProvider } from 'react-router'
import MainPage from '../pages/MainPage'
import { CustomizationPage } from '../pages/CustomizationPage'

const routes = [
  { path: '/', element: <MainPage /> },
  { path: '/customize', element: <CustomizationPage /> },
]

describe('routing', () => {
  it('renders MainPage at /', () => {
    const router = createMemoryRouter(routes, { initialEntries: ['/'] })
    render(<RouterProvider router={router} />)
    expect(screen.getByText('Seasonal')).toBeInTheDocument()
  })

  it('renders CustomizationPage at /customize', () => {
    const router = createMemoryRouter(routes, { initialEntries: ['/customize'] })
    render(<RouterProvider router={router} />)
    expect(screen.getByText('Customize')).toBeInTheDocument()
  })

  it('Browse link on MainPage points to /', () => {
    const router = createMemoryRouter(routes, { initialEntries: ['/'] })
    render(<RouterProvider router={router} />)
    const browseLinks = screen.getAllByRole('link', { name: 'Browse' })
    browseLinks.forEach(link => expect(link).toHaveAttribute('href', '/'))
  })

  it('Custom link on MainPage points to /customize', () => {
    const router = createMemoryRouter(routes, { initialEntries: ['/'] })
    render(<RouterProvider router={router} />)
    const customizeLink = screen.getByRole('link', { name: 'Custom' })
    expect(customizeLink).toHaveAttribute('href', '/customize')
  })
})
