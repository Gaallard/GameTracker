import { createContext, useContext } from 'react'

type User = {
  id?: number
  username?: string
  email?: string
  firstName?: string
  lastName?: string
  createdAt?: string
  updatedAt?: string
}

type AuthContextValue = {
  user: User | null
  login?: (...args: any[]) => any
  register?: (...args: any[]) => any
  logout: () => void
  isLoading?: boolean
  isAuthenticated?: boolean
}

// valor por defecto inofensivo
const defaultValue: AuthContextValue = {
  user: null,
  logout: () => {},
  isLoading: false,
  isAuthenticated: false,
}

const AuthContext = createContext<AuthContextValue>(defaultValue)

export const useAuth = () => useContext(AuthContext)

// opcional: provider mínimo para que compile si lo llegás a usar
export const AuthProvider: React.FC<React.PropsWithChildren> = ({ children }) => {
  return <AuthContext.Provider value={defaultValue}>{children}</AuthContext.Provider>
}
