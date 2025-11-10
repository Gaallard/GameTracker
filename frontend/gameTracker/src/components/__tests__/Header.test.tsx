import React from 'react'
import { vi, describe, it, expect, beforeEach } from 'vitest'
import { render, screen, fireEvent } from '@testing-library/react'
import type { AuthContextType } from '../../contexts/AuthContext'

// 1) Estado hoisted para el mock del contexto (no rompe con el hoisting de vi.mock)
const auth = vi.hoisted(() => ({ state: {} as Partial<AuthContextType> }))

// 2) Mock del contexto usando el MISMO path relativo que usa Header.tsx
vi.mock('../../contexts/AuthContext', () => ({
  useAuth: () => auth.state,
}))

// 3) Mock del botón por la ruta que resuelve Header ('./button' -> '../ui/button' desde este test)
vi.mock('../ui/button', () => {
  const Button = (props: React.ComponentProps<'button'>) => React.createElement('button', props, props.children)
  return { Button }
})

// 4) Import del Header y resolución default/named
import * as HeaderModule from '../ui/Header'
const Header = (HeaderModule.default ?? HeaderModule.Header) as React.ComponentType

// 5) Mock window.confirm
const mockConfirm = vi.fn()
Object.defineProperty(window, 'confirm', {
  value: mockConfirm,
  writable: true,
})

describe('Header Component', () => {
  const mockUser = {
    id: 1,
    username: 'testuser',
    email: 'test@example.com',
    firstName: 'Test',
    lastName: 'User',
    createdAt: '2024-01-01T00:00:00Z',
    updatedAt: '2024-01-01T00:00:00Z',
  }

  const mockLogout = vi.fn()

  beforeEach(() => {
    vi.clearAllMocks()
    mockConfirm.mockReturnValue(true)
  })

  it('should render header with user information', () => {
    auth.state = {
      user: mockUser,
      logout: mockLogout,
      login: vi.fn(),
      register: vi.fn(),
      isLoading: false,
      isAuthenticated: true,
    }

    render(<Header />)

    expect(screen.getByText('Game Tracker')).toBeInTheDocument()
    expect(screen.getByText(/¡Hola,\s*Test!/)).toBeInTheDocument()
    expect(screen.getByText('test@example.com')).toBeInTheDocument()
    expect(screen.getByText('Cerrar Sesión')).toBeInTheDocument()
  })

  it('should render header with username when firstName is not available', () => {
    auth.state = {
      user: { ...mockUser, firstName: undefined },
      logout: mockLogout,
      login: vi.fn(),
      register: vi.fn(),
      isLoading: false,
      isAuthenticated: true,
    }

    render(<Header />)
    expect(screen.getByText(/¡Hola,\s*testuser!/)).toBeInTheDocument()
  })

  it('should call logout when logout button is clicked and user confirms', () => {
    auth.state = {
      user: mockUser,
      logout: mockLogout,
      login: vi.fn(),
      register: vi.fn(),
      isLoading: false,
      isAuthenticated: true,
    }
    mockConfirm.mockReturnValue(true)

    render(<Header />)
    fireEvent.click(screen.getByText('Cerrar Sesión'))

    expect(mockConfirm).toHaveBeenCalledWith('¿Estás seguro de que quieres cerrar sesión?')
    expect(mockLogout).toHaveBeenCalledTimes(1)
  })

  it('should not call logout when user cancels confirmation', () => {
    auth.state = {
      user: mockUser,
      logout: mockLogout,
      login: vi.fn(),
      register: vi.fn(),
      isLoading: false,
      isAuthenticated: true,
    }
    mockConfirm.mockReturnValue(false)

    render(<Header />)
    fireEvent.click(screen.getByText('Cerrar Sesión'))

    expect(mockConfirm).toHaveBeenCalledWith('¿Estás seguro de que quieres cerrar sesión?')
    expect(mockLogout).not.toHaveBeenCalled()
  })

  it('should render game controller emoji', () => {
    auth.state = {
      user: mockUser,
      logout: mockLogout,
      login: vi.fn(),
      register: vi.fn(),
      isLoading: false,
      isAuthenticated: true,
    }

    render(<Header />)
    expect(screen.getByText('🎮')).toBeInTheDocument()
  })

  it('should have correct CSS classes', () => {
    auth.state = {
      user: mockUser,
      logout: mockLogout,
      login: vi.fn(),
      register: vi.fn(),
      isLoading: false,
      isAuthenticated: true,
    }

    render(<Header />)

    const header = screen.getByRole('banner')
    expect(header).toHaveClass('bg-white', 'shadow-sm', 'border-b')

    const title = screen.getByText('Game Tracker')
    expect(title).toHaveClass('text-2xl', 'font-bold', 'text-gray-900')
  })
})