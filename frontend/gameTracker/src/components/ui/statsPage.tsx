// src/pages/StatsPage.tsx
import { useEffect, useState } from "react"
import { getStats } from "@/services/api"

type Stats = {
    total_games: number
    by_status: Record<string, number>
    average_hours_played: number
    most_played_genre: string
    pending_games: number
}

export default function StatsPage() {
    const [stats, setStats] = useState<Stats | null>(null)
    const [loading, setLoading] = useState(true)

    useEffect(() => {
        const fetchStats = async () => {
            try {
                const res = await getStats()
                setStats(res.data)
            } catch (err) {
                console.error("Error cargando estadísticas:", err)
            } finally {
                setLoading(false)
            }
        }

        fetchStats()
    }, [])

    if (loading) return <p className="text-center">Cargando estadísticas...</p>
    if (!stats) return <p className="text-center text-red-500">No se pudieron cargar las estadísticas.</p>

    return (
        <div className="max-w-md mx-auto bg-white shadow rounded-xl p-6 space-y-4">
            <h2 className="text-2xl font-bold text-center">📊 Estadísticas</h2>
            <p>Total de juegos: <strong>{stats.total_games}</strong></p>
            <p>Promedio de horas jugadas: <strong>{stats.average_hours_played.toFixed(2)}</strong></p>
            <p>Género más jugado: <strong>{stats.most_played_genre || "N/A"}</strong></p>
            <p>Juegos pendientes: <strong>{stats.pending_games}</strong></p>

            <div>
                <h3 className="font-semibold mt-4 mb-2">Por estado:</h3>
                <ul className="list-disc list-inside">
                    {Object.entries(stats.by_status).map(([status, count]) => (
                        <li key={status}>{status}: {count}</li>
                    ))}
                </ul>
            </div>
        </div>
    )
}
