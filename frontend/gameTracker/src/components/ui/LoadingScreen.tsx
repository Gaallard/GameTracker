import React from 'react'

const LoadingScreen: React.FC = () => {
    return (
        <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-purple-600 via-blue-600 to-indigo-700">
            <div className="text-center">
                <div className="mx-auto w-20 h-20 bg-white rounded-full flex items-center justify-center mb-6 shadow-2xl">
                    <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-purple-600"></div>
                </div>
                <h2 className="text-2xl font-bold text-white mb-2">Game Tracker</h2>
                <p className="text-white/80">Cargando tu experiencia de juego...</p>
            </div>
        </div>
    )
}

export default LoadingScreen
