import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import LoadingScreen from '../ui/LoadingScreen'

describe('LoadingScreen Component', () => {
  it('should render loading screen with correct content', () => {
    // Act
    render(<LoadingScreen />)

    // Assert
    expect(screen.getByText('Game Tracker')).toBeInTheDocument()
    expect(screen.getByText('Cargando tu experiencia de juego...')).toBeInTheDocument()
  })

  it('should have correct CSS classes for styling', () => {
    // Act
    render(<LoadingScreen />)

    // Assert
    const mainDiv = screen.getByText('Game Tracker').closest('div')?.parentElement
    expect(mainDiv).toHaveClass('min-h-screen', 'flex', 'items-center', 'justify-center', 'bg-gradient-to-br', 'from-purple-600', 'via-blue-600', 'to-indigo-700')

    const title = screen.getByText('Game Tracker')
    expect(title).toHaveClass('text-2xl', 'font-bold', 'text-white', 'mb-2')

    const subtitle = screen.getByText('Cargando tu experiencia de juego...')
    expect(subtitle).toHaveClass('text-white/80')
  })

  it('should render loading spinner', () => {
    // Act
    render(<LoadingScreen />)

    // Assert
    const spinner = document.querySelector('.animate-spin')
    expect(spinner).toBeInTheDocument()
    expect(spinner).toHaveClass('animate-spin', 'rounded-full', 'h-12', 'w-12', 'border-b-2', 'border-purple-600')
  })

  it('should have proper structure with container and content', () => {
    // Act
    render(<LoadingScreen />)

    // Assert
    const container = screen.getByText('Game Tracker').closest('.text-center')
    expect(container).toBeInTheDocument()
    expect(container).toHaveClass('text-center')

    const whiteCircle = document.querySelector('.bg-white.rounded-full')
    expect(whiteCircle).toBeInTheDocument()
    expect(whiteCircle).toHaveClass('mx-auto', 'w-20', 'h-20', 'bg-white', 'rounded-full', 'flex', 'items-center', 'justify-center', 'mb-6', 'shadow-2xl')
  })
})
