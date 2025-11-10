import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import axios from 'axios'
import type { AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios'

// ⬇️ Importá SOLO TIPOS del módulo (no ejecuta código en runtime)
import type {
  Game,
  GameStats,
  LoginRequest,
  RegisterRequest,
  AuthResponse,
} from '../api'

// Helper: instancia "fake" de axios con interceptores válidos
type MockAxiosInstance = {
  get: ReturnType<typeof vi.fn>
  post: ReturnType<typeof vi.fn>
  put: ReturnType<typeof vi.fn>
  delete: ReturnType<typeof vi.fn>
  interceptors: {
    request: { use: ReturnType<typeof vi.fn>; eject: ReturnType<typeof vi.fn> }
    response: { use: ReturnType<typeof vi.fn>; eject: ReturnType<typeof vi.fn> }
  }
}

const makeAxiosInstance = (): MockAxiosInstance => ({
  get: vi.fn(),
  post: vi.fn(),
  put: vi.fn(),
  delete: vi.fn(),
  interceptors: {
    request: { use: vi.fn(), eject: vi.fn() },
    response: { use: vi.fn(), eject: vi.fn() },
  },
})

type ApiModule = typeof import('../api')

let instance: MockAxiosInstance
let api: ApiModule

describe('API Service', () => {
  beforeEach(async () => {
    vi.resetModules() // fuerza re-import del módulo para recrear la instancia con nuestro spy
    vi.clearAllMocks()

    // mock de axios.create
    instance = makeAxiosInstance()
    vi.spyOn(axios, 'create').mockReturnValue(instance as unknown as AxiosInstance)

    // mock de localStorage
    Object.defineProperty(window, 'localStorage', {
      value: {
        getItem: vi.fn(),
        setItem: vi.fn(),
        removeItem: vi.fn(),
        clear: vi.fn(),
      },
      writable: true,
    })

    // importar recién ahora tu API (usa la instancia espía)
    api = await import('../api')
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  describe('Game endpoints', () => {
    const mockGame: Game = {
      id: 1,
      title: 'Test Game',
      platform: 'PC',
      genre: 'RPG',
      status: 'Completed',
      progress: 100,
      hoursPlayed: 25.5,
      personalNote: 'Great game',
      score: 8,
      startedAt: '2024-01-01T00:00:00Z',
      finishedAt: '2024-01-15T00:00:00Z',
      coverURL: 'http://example.com/cover.jpg',
      createdAt: '2024-01-01T00:00:00Z',
      updatedAt: '2024-01-15T00:00:00Z',
    }

    it('should get all games', async () => {
      const mockResponse: AxiosResponse<Game[]> = { data: [mockGame] } as AxiosResponse<Game[]>
      instance.get.mockResolvedValue(mockResponse)

      const result = await api.getGames()
      expect(result.data).toEqual([mockGame])
      expect(instance.get).toHaveBeenCalled()
    })

    it('should get game by ID', async () => {
      const mockResponse: AxiosResponse<Game> = { data: mockGame } as AxiosResponse<Game>
      instance.get.mockResolvedValue(mockResponse)

      const result = await api.getGameById(1)
      expect(result.data).toEqual(mockGame)
      expect(instance.get).toHaveBeenCalled()
    })

    it('should create a new game', async () => {
      const newGameData: Partial<Game> = {
        title: 'New Game',
        platform: 'PC',
        genre: 'RPG',
        status: 'Not Started',
        progress: 0,
        hoursPlayed: 0,
        personalNote: 'New game to play',
        score: 0,
        coverURL: 'http://example.com/cover.jpg',
      }
      const mockResponse: AxiosResponse<Game> = { data: { ...newGameData, id: 1 } as Game } as AxiosResponse<Game>
      instance.post.mockResolvedValue(mockResponse)

      const result = await api.createGame(newGameData)
      expect(result.data).toEqual({ ...newGameData, id: 1 })
      expect(instance.post).toHaveBeenCalled()
    })

    it('should update a game', async () => {
      const updateData: Partial<Game> = { title: 'Updated Game', score: 9 }
      const mockResponse: AxiosResponse<Game> = { data: { ...mockGame, ...updateData } } as AxiosResponse<Game>
      instance.put.mockResolvedValue(mockResponse)

      const result = await api.updateGame(1, updateData)
      expect(result.data).toEqual({ ...mockGame, ...updateData })
      expect(instance.put).toHaveBeenCalled()
    })

    it('should delete a game', async () => {
      const mockResponse: AxiosResponse<{ message: string }> = { data: { message: 'Game deleted successfully' } } as AxiosResponse<{ message: string }>
      instance.delete.mockResolvedValue(mockResponse)

      const result = await api.deleteGame(1)
      expect(result.data).toEqual({ message: 'Game deleted successfully' })
      expect(instance.delete).toHaveBeenCalled()
    })

    it('should get game statistics', async () => {
      const mockStats: GameStats = {
        total_games: 10,
        average_hours_played: 15.5,
        most_played_genre: 'RPG',
        pending_games: 3,
        by_status: {
          Completed: 5,
          Playing: 2,
          'Not Started': 3,
        },
      }
      const mockResponse: AxiosResponse<GameStats> = { data: mockStats } as AxiosResponse<GameStats>
      instance.get.mockResolvedValue(mockResponse)

      const result = await api.getStats()
      expect(result.data).toEqual(mockStats)
      expect(instance.get).toHaveBeenCalled()
    })
  })

  describe('Auth endpoints', () => {
    const mockAuthResponse: AuthResponse = {
      token: 'mock-jwt-token',
      user: {
        id: 1,
        username: 'testuser',
        email: 'test@example.com',
        firstName: 'Test',
        lastName: 'User',
        createdAt: '2024-01-01T00:00:00Z',
        updatedAt: '2024-01-01T00:00:00Z',
      },
    }

    it('should login user', async () => {
      const loginData: LoginRequest = { username: 'testuser', password: 'password123' }
      const mockResponse: AxiosResponse<AuthResponse> = { data: mockAuthResponse } as AxiosResponse<AuthResponse>
      instance.post.mockResolvedValue(mockResponse)

      const result = await api.login(loginData)
      expect(result.data).toEqual(mockAuthResponse)
      expect(instance.post).toHaveBeenCalled()
    })

    it('should register user', async () => {
      const registerData: RegisterRequest = {
        username: 'newuser',
        email: 'newuser@example.com',
        password: 'password123',
        firstName: 'New',
        lastName: 'User',
      }
      const mockResponse: AxiosResponse<AuthResponse> = { data: mockAuthResponse } as AxiosResponse<AuthResponse>
      instance.post.mockResolvedValue(mockResponse)

      const result = await api.register(registerData)
      expect(result.data).toEqual(mockAuthResponse)
      expect(instance.post).toHaveBeenCalled()
    })

    it('should get user profile', async () => {
      const mockResponse: AxiosResponse<typeof mockAuthResponse.user> = { data: mockAuthResponse.user } as AxiosResponse<typeof mockAuthResponse.user>
      instance.get.mockResolvedValue(mockResponse)

      const result = await api.getProfile()
      expect(result.data).toEqual(mockAuthResponse.user)
      expect(instance.get).toHaveBeenCalled()
    })
  })

  describe('API interceptors', () => {
    it('adds Authorization header when token exists', async () => {
      const mockGetItem = vi.fn().mockReturnValue('mock-token')
      Object.defineProperty(window, 'localStorage', {
        value: { getItem: mockGetItem },
        writable: true,
      })

      // El interceptor se cargó al importar; lo invocamos manualmente
      const reqHandler = instance.interceptors.request.use.mock.calls[0][0] as (config: AxiosRequestConfig) => AxiosRequestConfig | Promise<AxiosRequestConfig>
      const cfg = await reqHandler({ headers: {} } as AxiosRequestConfig)

      expect(cfg.headers?.Authorization).toBe('Bearer mock-token')
    })

    it('does not add Authorization header when token is missing', async () => {
      const mockGetItem = vi.fn().mockReturnValue(null)
      Object.defineProperty(window, 'localStorage', {
        value: { getItem: mockGetItem },
        writable: true,
      })

      const reqHandler = instance.interceptors.request.use.mock.calls[0][0] as (config: AxiosRequestConfig) => AxiosRequestConfig | Promise<AxiosRequestConfig>
      const cfg = await reqHandler({ headers: {} } as AxiosRequestConfig)

      expect(cfg.headers?.Authorization).toBeUndefined()
    })
  })
})
