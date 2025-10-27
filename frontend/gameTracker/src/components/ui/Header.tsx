import React from 'react'
// ⬇️ usar ruta relativa en vez de '@/contexts/AuthContext'
import { useAuth } from '../../contexts/AuthContext'
// ⬇️ y también relativa para el botón (está en la misma carpeta)
import { Button } from './button'

const Header: React.FC = () => {
  const { user, logout } = useAuth()

  const handleLogout = () => {
    if (window.confirm('¿Estás seguro de que quieres cerrar sesión?')) {
      logout()
    }
  }

  return (
    <header role="banner" className="bg-white shadow-sm border-b">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center py-4">
          <div className="flex items-center space-x-4">
            <div className="w-8 h-8 bg-gradient-to-r from-purple-600 to-blue-600 rounded-full flex items-center justify-center">
              <span className="text-white text-sm font-bold">🎮</span>
            </div>
            <h1 className="text-2xl font-bold text-gray-900">Game Tracker</h1>
          </div>

          <div className="flex items-center space-x-4">
            <div className="text-right">
              <p className="text-sm font-medium text-gray-900">
                ¡Hola, {user?.firstName || user?.username}!
              </p>
              <p className="text-xs text-gray-500">{user?.email}</p>
            </div>

            <Button
              onClick={handleLogout}
              variant="outline"
              className="text-gray-600 hover:text-gray-900 hover:bg-gray-50"
            >
              Cerrar Sesión
            </Button>
          </div>
        </div>
      </div>
    </header>
  )
}

export default Header
export { Header }
