# Tests & Coverage - GameTracker

Este documento describe la configuración completa de pruebas unitarias, cobertura de código y CI/CD para el proyecto GameTracker.

## 📋 Resumen

El proyecto implementa un sistema completo de testing que incluye:
- **Backend (Go)**: Pruebas unitarias con cobertura mínima del 70%
- **Frontend (React/TypeScript)**: Pruebas unitarias con Vitest
- **CI/CD**: GitHub Actions con reportes automáticos
- **Reportes**: JUnit XML y HTML para análisis detallado

## 🔧 Herramientas Utilizadas

### Backend (Go)
- **gotestsum**: Ejecutor de pruebas con reportes JUnit
- **go-junit-report**: Generador de reportes JUnit
- **go test**: Framework nativo de Go para testing
- **sqlmock**: Mock de base de datos para pruebas aisladas
- **testify**: Librería de assertions y mocks

### Frontend (React/TypeScript)
- **Vitest**: Framework de testing moderno y rápido
- **@testing-library/react**: Utilidades para testing de componentes React
- **@testing-library/jest-dom**: Matchers adicionales para DOM
- **jsdom**: Entorno de testing para DOM
- **@vitest/coverage-v8**: Herramienta de cobertura de código

## 🚀 Ejecutar Pruebas Localmente

### Backend

```bash
cd backend

# Instalar herramientas de testing (solo la primera vez)
go install gotest.tools/gotestsum@latest
go install github.com/jstemmer/go-junit-report/v2@latest

# Ejecutar pruebas con cobertura
gotestsum --format=standard-verbose --junitfile test-results-go.xml -- \
  -covermode=atomic -coverprofile=coverage.out ./...

# Generar reportes de cobertura
go tool cover -func=coverage.out > coverage.txt
go tool cover -html=coverage.out -o coverage.html

# Verificar cobertura mínima del 70%
awk '/total:/ { gsub("%","",$3); if ($3+0 < 70) { print "❌ Coverage below 70%: " $3"%"; exit 1 } else { print "✅ Coverage: " $3"%" } }' coverage.txt
```

**Script automatizado (Linux/Mac):**
```bash
cd backend
chmod +x run-tests.sh
./run-tests.sh
```

**Script automatizado (Windows):**
```cmd
cd backend
run-tests.bat
```

### Frontend

```bash
cd frontend/gameTracker

# Instalar dependencias (solo la primera vez)
npm ci

# Ejecutar pruebas
npm run test

# Ejecutar pruebas con cobertura
npm run coverage
```

## 📊 Cobertura de Código

### Backend - Funciones Cubiertas

1. **Service Layer** (`service/service_test.go`):
   - `GetAllGames()` - Obtener todos los juegos
   - `GetGameByID()` - Obtener juego por ID (casos exitoso y error)
   - `CreateGame()` - Crear nuevo juego
   - `UpdateGame()` - Actualizar juego existente
   - `DeleteGame()` - Eliminar juego
   - `GetByTitle()` - Buscar juegos por título
   - `GetByStatus()` - Buscar juegos por estado
   - `GetByGenre()` - Buscar juegos por género
   - `GetStats()` - Obtener estadísticas de juegos

2. **Controller Layer** (`controller/controller_test.go`):
   - `GetAllGames()` - Handler GET /games
   - `GetGameByID()` - Handler GET /games/:id
   - `CreateGame()` - Handler POST /games
   - `UpdateGame()` - Handler PUT /games/:id
   - `DeleteGame()` - Handler DELETE /games/:id
   - `GetByTitle()` - Handler GET /games/search/title
   - `GetByStatus()` - Handler GET /games/search/status
   - `GetByGenre()` - Handler GET /games/search/genre
   - `GetStats()` - Handler GET /games/stats

### Frontend - Componentes y Utilidades Cubiertos

1. **Utilidades** (`lib/__tests__/utils.test.ts`):
   - `cn()` - Función de merge de clases CSS
   - Casos: clases simples, condicionales, arrays, objetos, valores nulos

2. **Servicios API** (`services/__tests__/api.test.ts`):
   - Endpoints de juegos: GET, POST, PUT, DELETE
   - Endpoints de autenticación: login, register, profile
   - Interceptores de axios para tokens y errores

3. **Componentes React** (`components/__tests__/`):
   - `Header.test.tsx` - Componente de navegación
   - `LoadingScreen.test.tsx` - Pantalla de carga
   - Casos: renderizado, interacciones, props, CSS classes

## 🏗️ Arquitectura de Testing

### Backend - Patrón AAA (Arrange, Act, Assert)

```go
func TestGetGameByID_Success(t *testing.T) {
    // Arrange
    db, mock := setupTestDB(t)
    defer db.Close()
    
    // Configurar mock data...
    
    // Act
    result, err := GetGameByID("1")
    
    // Assert
    require.NoError(t, err)
    assert.Equal(t, expectedGame.Title, result.Title)
}
```

### Frontend - Testing de Componentes

```typescript
describe('Header Component', () => {
  it('should render header with user information', () => {
    // Arrange
    vi.mocked(useAuth).mockReturnValue({
      user: mockUser,
      logout: mockLogout,
      // ...
    })

    // Act
    render(<Header />)

    // Assert
    expect(screen.getByText('Game Tracker')).toBeInTheDocument()
  })
})
```

## 🔄 CI/CD Pipeline

### GitHub Actions Workflow

El pipeline se ejecuta en cada push y pull request:

1. **Backend Job**:
   - Instala Go 1.22.x
   - Instala herramientas de testing
   - Ejecuta pruebas con cobertura
   - Verifica cobertura mínima del 70%
   - Sube reportes como artefactos

2. **Frontend Job**:
   - Instala Node.js 20
   - Instala dependencias
   - Ejecuta pruebas con cobertura
   - Sube reportes como artefactos

3. **Lint Job**:
   - Verifica calidad de código Go
   - Ejecuta ESLint para TypeScript

### Artefactos Generados

- `test-results-go.xml` - Reporte JUnit del backend
- `coverage.out` - Datos de cobertura del backend
- `coverage.txt` - Resumen de cobertura del backend
- `coverage.html` - Reporte HTML de cobertura del backend
- `test-results-vue.xml` - Reporte JUnit del frontend
- `coverage/` - Reportes de cobertura del frontend

## 📈 Métricas de Cobertura

### Criterios de Aceptación

- **Cobertura mínima**: 70% para statements, branches, functions y lines
- **Pipeline falla** si la cobertura está por debajo del umbral
- **Reportes automáticos** en cada ejecución de CI

### Tipos de Cobertura Medidos

1. **Statements**: Porcentaje de declaraciones ejecutadas
2. **Branches**: Porcentaje de ramas condicionales probadas
3. **Functions**: Porcentaje de funciones llamadas
4. **Lines**: Porcentaje de líneas ejecutadas

## 🛠️ Configuración Técnica

### Backend - Mock de Base de Datos

```go
func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
    db, mock, err := sqlmock.New()
    require.NoError(t, err)
    
    gormDB, err := gorm.Open(mysql.New(mysql.Config{
        Conn:                      db,
        SkipInitializeWithVersion: true,
    }), &gorm.Config{})
    
    return gormDB, mock
}
```

### Frontend - Configuración de Vitest

```typescript
export default defineConfig({
  plugins: [react()],
  test: {
    environment: 'jsdom',
    reporters: ['default', 'junit'],
    outputFile: { junit: 'test-results-vue.xml' },
    coverage: {
      provider: 'v8',
      reporter: ['text', 'lcov', 'cobertura', 'html'],
      reportsDirectory: 'coverage',
      statements: 70,
      branches: 70,
      functions: 70,
      lines: 70
    }
  }
})
```

## 🎯 Justificación de Decisiones

### Backend

1. **sqlmock**: Permite testing aislado sin dependencias de base de datos real
2. **gotestsum**: Proporciona reportes detallados y compatibilidad con CI
3. **testify**: Simplifica assertions y mejora legibilidad de tests
4. **Patrón AAA**: Estructura clara y mantenible para tests

### Frontend

1. **Vitest**: Más rápido que Jest, mejor integración con Vite
2. **@testing-library**: Enfoque en testing de comportamiento del usuario
3. **jsdom**: Simula DOM real para testing de componentes
4. **Coverage v8**: Mejor rendimiento que Istanbul

### CI/CD

1. **GitHub Actions**: Integración nativa con GitHub
2. **Artefactos**: Preservación de reportes para análisis posterior
3. **Umbrales de cobertura**: Garantiza calidad mínima del código
4. **Paralelización**: Jobs independientes para backend y frontend

## 🔍 Troubleshooting

### Backend

**Error: "no test files"**
```bash
# Verificar que los archivos terminen en _test.go
ls -la *_test.go
```

**Error de cobertura baja**
```bash
# Ejecutar con verbose para ver qué no está cubierto
go test -v -cover ./...
```

### Frontend

**Error: "Cannot resolve module"**
```bash
# Reinstalar dependencias
rm -rf node_modules package-lock.json
npm install
```

**Error de cobertura**
```bash
# Verificar configuración de Vitest
npx vitest --coverage --reporter=verbose
```

## 📚 Referencias

- [Go Testing Documentation](https://golang.org/pkg/testing/)
- [Vitest Documentation](https://vitest.dev/)
- [Testing Library Documentation](https://testing-library.com/)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)

---

**Última actualización**: Enero 2024  
**Versión**: 1.0.0  
**Mantenedor**: Equipo GameTracker
