import { describe, it, expect } from 'vitest'
import { cn } from '../utils'

describe('Utils', () => {
  describe('cn function', () => {
    it('should merge class names correctly', () => {
      const result = cn('text-red-500', 'bg-blue-500')

      const toSet = (s: string) => new Set(s.split(/\s+/).filter(Boolean))
      
      expect(toSet(result)).toEqual(toSet('bg-blue-500 text-red-500'))
    })
    

    it('should handle conditional classes', () => {
      // Arrange
      const isActive = true
      const isDisabled = false
      
      // Act
      const result = cn(
        'base-class',
        isActive && 'active-class',
        isDisabled && 'disabled-class'
      )
      
      // Assert
      expect(result).toBe('base-class active-class')
    })

    it('should handle empty inputs', () => {
      // Arrange & Act
      const result = cn()
      
      // Assert
      expect(result).toBe('')
    })

    it('should handle undefined and null values', () => {
      // Arrange & Act
      const result = cn('base-class', undefined, null, 'another-class')
      
      // Assert
      expect(result).toBe('base-class another-class')
    })

    it('should merge conflicting Tailwind classes correctly', () => {
      // Arrange & Act
      const result = cn('text-red-500', 'text-blue-500')
      
      // Assert
      expect(result).toBe('text-blue-500')
    })

    it('should handle arrays of classes', () => {
      // Arrange & Act
      const result = cn(['class1', 'class2'], 'class3')
      
      // Assert
      expect(result).toBe('class1 class2 class3')
    })

    it('should handle objects with boolean values', () => {
      // Arrange & Act
      const result = cn({
        'active': true,
        'disabled': false,
        'visible': true
      })
      
      // Assert
      expect(result).toBe('active visible')
    })
  })
})
