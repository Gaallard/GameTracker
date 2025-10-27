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
    baseURL: import.meta.env.VITE_API_URL || "", // Usar variable de entorno o rutas relativas
    headers: {
        "Content-Type": "application/json",
    },
})

// Interceptor para agregar token a las peticiones
API.interceptors.request.use(
    (config) => {
        const token = localStorage.getItem('token')
        if (token) {
            config.headers.Authorization = `Bearer ${token}`
        }
        return config
    },
    (error) => {
        return Promise.reject(error)
    }
)

// Interceptor para manejar respuestas de error
API.interceptors.response.use(
    (response) => response,
    (error) => {
        if (error.response?.status === 401) {
            // Token expirado o inválido
            localStorage.removeItem('token')
            localStorage.removeItem('user')
            window.location.reload()
        }
        return Promise.reject(error)
    }
)

// Game endpoints
export const getGames = () => API.get<Game[]>("/games/")
export const getGameById = (id: number) => API.get<Game>(`/games/${id}`)
export const searchGameByTitle = (title: string) => API.get<Game[]>(`/games/search?title=${title}`)
export const createGame = (data: Partial<Game>) => API.post("/games/", data)
export const getStats = () => API.get("/stats")
export const updateGame = (id: number, data: Partial<Game>) => API.put<Game>(`/games/${id}`, data)
export const deleteGame = (id: number) => API.delete(`/games/${id}`)

// Auth endpoints
export interface LoginRequest {
    username: string
    password: string
}

export interface RegisterRequest {
    username: string
    email: string
    password: string
    firstName?: string
    lastName?: string
}

export interface User {
    id: number
    username: string
    email: string
    firstName?: string
    lastName?: string
    createdAt: string
    updatedAt: string
}

export interface AuthResponse {
    token: string
    user: User
}

export const login = (data: LoginRequest) => API.post<AuthResponse>("/auth/login", data)
export const register = (data: RegisterRequest) => API.post<AuthResponse>("/auth/register", data)
export const getProfile = () => API.get("/api/profile")

export default API
