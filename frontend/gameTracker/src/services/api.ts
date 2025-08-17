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
    totalGames: number
    averageHours: number
    mostPlayedGenre: string
    pendingGames: number
    byStatus: Record<string, number>
}

const API = axios.create({
    baseURL: "http://localhost:8080", // Cambiá si tu backend corre en otra URL/puerto
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
