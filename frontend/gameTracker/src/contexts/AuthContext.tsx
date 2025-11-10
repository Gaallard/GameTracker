import React, { createContext, useContext, useState, useEffect, type ReactNode } from 'react'
import { login as loginAPI, register as registerAPI, type User, type LoginRequest, type RegisterRequest } from '@/services/api'

interface AuthContextType {
    user: User | null
    token: string | null
    isAuthenticated: boolean
    isLoading: boolean
    login: (data: LoginRequest) => Promise<void>
    register: (data: RegisterRequest) => Promise<void>
    logout: () => void
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export const useAuth = () => {
    const context = useContext(AuthContext)
    if (context === undefined) {
        throw new Error('useAuth must be used within an AuthProvider')
    }
    return context
}

interface AuthProviderProps {
    children: ReactNode
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
    const [user, setUser] = useState<User | null>(null)
    const [token, setToken] = useState<string | null>(null)
    const [isLoading, setIsLoading] = useState(true)

    const isAuthenticated = !!user && !!token

    useEffect(() => {
        // Verificar si hay token guardado al cargar la app
        const savedToken = localStorage.getItem('token')
        const savedUser = localStorage.getItem('user')
        
        if (savedToken && savedUser) {
            try {
                setToken(savedToken)
                setUser(JSON.parse(savedUser))
            } catch (error) {
                console.error('Error parsing saved user:', error)
                localStorage.removeItem('token')
                localStorage.removeItem('user')
            }
        }
        setIsLoading(false)
    }, [])

    const login = async (data: LoginRequest) => {
        try {
            setIsLoading(true)
            const response = await loginAPI(data)
            const { token: newToken, user: newUser } = response.data
            
            setToken(newToken)
            setUser(newUser)
            
            // Guardar en localStorage
            localStorage.setItem('token', newToken)
            localStorage.setItem('user', JSON.stringify(newUser))
        } catch (error: any) {
            console.error('Login error:', error)
            throw new Error(error.response?.data?.error || 'Error al iniciar sesión')
        } finally {
            setIsLoading(false)
        }
    }

    const register = async (data: RegisterRequest) => {
        try {
            setIsLoading(true)
            const response = await registerAPI(data)
            const { token: newToken, user: newUser } = response.data
            
            setToken(newToken)
            setUser(newUser)
            
            // Guardar en localStorage
            localStorage.setItem('token', newToken)
            localStorage.setItem('user', JSON.stringify(newUser))
        } catch (error: any) {
            console.error('Register error:', error)
            throw new Error(error.response?.data?.error || 'Error al registrarse')
        } finally {
            setIsLoading(false)
        }
    }

    const logout = () => {
        setUser(null)
        setToken(null)
        localStorage.removeItem('token')
        localStorage.removeItem('user')
    }

    const value: AuthContextType = {
        user,
        token,
        isAuthenticated,
        isLoading,
        login,
        register,
        logout,
    }

    return (
        <AuthContext.Provider value={value}>
            {children}
        </AuthContext.Provider>
    )
}
