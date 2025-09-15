import axios from "axios"

// Definimos el tipo de juego
export interface Game {
    id: number
    title: string
    platform: string
    genre: string
    status: string
    progress: number
    hoursPlayed: number
    personalNote: string
    score: number
    startedAt: string
    finishedAt: string
    coverURL: string
    createdAt: string
    updatedAt: string
}
export interface GameStats {
    total_games: number
    average_hours_played: number
    most_played_genre: string
    pending_games: number
    by_status: Record<string, number>
}

const API = axios.create({
    baseURL: "", // Usar rutas relativas para que Nginx haga el proxy
    headers: {
        "Content-Type": "application/json",
    },
})

export const getGames = () => API.get<Game[]>("/games/")
export const getGameById = (id: number) => API.get<Game>(`/games/${id}`)
export const searchGameByTitle = (title: string) => API.get<Game[]>(`/games/search?title=${title}`)
export const createGame = (data: Partial<Game>) => API.post("/games/", data)
export const getStats = () => API.get("/stats")
export const updateGame = (id: number, data: Partial<Game>) => API.put<Game>(`/games/${id}`, data)
export const deleteGame = (id: number) => API.delete(`/games/${id}`)


export default API
