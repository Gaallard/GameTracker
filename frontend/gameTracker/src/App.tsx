import { useEffect, useState } from "react"
import { getGames, deleteGame, getStats, getGameById, type Game } from "@/services/api"
import {
    Card,
    CardHeader,
    CardTitle,
    CardDescription,
    CardContent,
} from "@/components/ui/card"
import GameForm from "@/components/ui/gameForm"
import { useNavigate } from "react-router-dom"



type Stats = {
    total_games: number
    by_status: Record<string, number>
    average_hours_played: number
    most_played_genre: string
    pending_games: number
}

function App() {
    const [games, setGames] = useState<Game[]>([])
    const [editingGame, setEditingGame] = useState<Game | null>(null)
    const [detailedGameId, setDetailedGameId] = useState<number | null>(null)
    const [detailedGameData, setDetailedGameData] = useState<Game | null>(null)
    const [stats, setStats] = useState<Stats | null>(null)
    const navigate = useNavigate()


    useEffect(() => {
        const fetchGames = async () => {
            try {
                const res = await getGames()
                setGames(res.data)
            } catch (err) {
                console.error("Error cargando juegos:", err)
            }
        }

        const fetchStats = async () => {
            try {
                const res = await getStats()
                setStats(res.data)
            } catch (err) {
                console.error("Error cargando estadísticas:", err)
            }
        }

        fetchGames()
        fetchStats()
    }, [])

    const refreshStats = async () => {
        try {
            const res = await getStats()
            setStats(res.data)
        } catch (err) {
            console.error("Error actualizando estadísticas:", err)
        }
    }

    const handleGameAddedOrUpdated = (game: Game) => {
        setGames((prevGames) => {
            const exists = prevGames.find((g) => g.id === game.id)
            if (exists) {
                return prevGames.map((g) => (g.id === game.id ? game : g))
            } else {
                return [...prevGames, game]
            }
        })
        setEditingGame(null)
        refreshStats()
    }

    const handleDeleteGame = async (id: number) => {
        if (!confirm("¿Querés eliminar este juego?")) return

        try {
            await deleteGame(id)
            setGames(prev => prev.filter(game => game.id !== id))
            alert("Juego eliminado correctamente")
            if (detailedGameId === id) {
                setDetailedGameId(null)
                setDetailedGameData(null)
            }
            refreshStats()
        } catch (err) {
            console.error("Error eliminando juego:", err)
            alert("No se pudo eliminar el juego")
        }
    }

    const handleViewDetails = async (id: number) => {
        if (detailedGameId === id) {
            setDetailedGameId(null)
            setDetailedGameData(null)
            return
        }

        try {
            const res = await getGameById(id)
            setDetailedGameId(id)
            setDetailedGameData(res.data)
        } catch (err) {
            console.error("Error cargando detalle del juego:", err)
            alert("No se pudo cargar el detalle del juego")
        }
    }

    return (
        <div className="min-h-svh flex flex-col items-center p-6 gap-6 bg-muted">
            <h1 className="text-3xl font-bold">🎮 Game Tracker</h1>

            {stats && (
                <div className="w-full max-w-3xl p-4 mb-4 bg-white rounded-lg shadow">
                    <h2 className="text-2xl font-semibold mb-2">📊 Estadísticas</h2>
                    <p><strong>Total de juegos:</strong> {stats.total_games}</p>
                    <p><strong>Juegos pendientes:</strong> {stats.pending_games}</p>
                    <p><strong>Horas promedio jugadas:</strong> {stats.average_hours_played.toFixed(2)}</p>
                    <p><strong>Género más jugado:</strong> {stats.most_played_genre}</p>

                    <div className="mt-2">
                        <strong>Juegos por estado:</strong>
                        <ul className="list-disc list-inside ml-4">
                            {Object.entries(stats.by_status).map(([status, count]) => (
                                <li key={status}>
                                    {status}: {count}
                                </li>
                            ))}
                        </ul>
                    </div>
                </div>
            )}


            <GameForm
                gameToEdit={editingGame}
                onGameAdded={handleGameAddedOrUpdated}
                onEditComplete={() => setEditingGame(null)}
            />
            <button
                onClick={() => navigate("/stats")}
                className="px-4 py-2 rounded bg-purple-600 text-white text-sm hover:bg-purple-700"
            >
                Ver estadísticas
            </button>
            <div className="w-full max-w-3xl grid gap-4">
                {games.length === 0 ? (
                    <Card>
                        <CardHeader>
                            <CardTitle>No hay juegos cargados</CardTitle>
                            <CardDescription>Agregá un juego para empezar</CardDescription>
                        </CardHeader>
                    </Card>
                ) : (
                    games.map((game) => (
                        <div key={game.id}>
                            <Card>
                                <CardHeader>
                                    <CardTitle>{game.title}</CardTitle>
                                    <CardDescription>
                                        {game.platform} - {game.genre}
                                    </CardDescription>
                                </CardHeader>
                                <CardContent>
                                    <p>Status: {game.status}</p>
                                    <p>Horas jugadas: {game.hoursPlayed}</p>

                                    <button
                                        className="mt-2 px-3 py-1 mr-2 rounded bg-blue-600 text-white text-sm hover:bg-blue-700"
                                        onClick={() => setEditingGame(game)}
                                    >
                                        Editar
                                    </button>

                                    <button
                                        className="mt-2 px-3 py-1 mr-2 rounded bg-green-600 text-white text-sm hover:bg-green-700"
                                        onClick={() => handleViewDetails(game.id)}
                                    >
                                        {detailedGameId === game.id ? "Cerrar detalles" : "Ver detalles"}
                                    </button>

                                    <button
                                        className="mt-2 px-3 py-1 rounded bg-red-600 text-white text-sm hover:bg-red-700"
                                        onClick={() => handleDeleteGame(game.id)}
                                    >
                                        Eliminar
                                    </button>
                                </CardContent>
                            </Card>

                            {detailedGameId === game.id && detailedGameData && (
                                <div className="p-4 mt-2 border rounded-xl shadow bg-white">
                                    <h3 className="text-xl font-semibold mb-2">{detailedGameData.title}</h3>
                                    <p><strong>Plataforma:</strong> {detailedGameData.platform}</p>
                                    <p><strong>Género:</strong> {detailedGameData.genre}</p>
                                    <p><strong>Status:</strong> {detailedGameData.status}</p>
                                    <p><strong>Progreso:</strong> {detailedGameData.progress}%</p>
                                    <p><strong>Horas jugadas:</strong> {detailedGameData.hoursPlayed}</p>
                                    <p><strong>Nota personal:</strong> {detailedGameData.personalNote}</p>
                                    <p><strong>Score:</strong> {detailedGameData.score}</p>
                                    <p>
                                        <strong>Comenzado:</strong> {new Date(detailedGameData.startedAt).toLocaleDateString()}
                                    </p>
                                    <p>
                                        <strong>Terminado:</strong> {new Date(detailedGameData.finishedAt).toLocaleDateString()}
                                    </p>
                                </div>
                            )}
                        </div>
                    ))
                )}
            </div>
        </div>
    )
}

export default App
