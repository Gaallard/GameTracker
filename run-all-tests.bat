@echo off
REM Script para ejecutar todas las pruebas (backend y frontend)

echo ==========================================
echo 🧪 Running All Tests (Backend + Frontend)
echo ==========================================
echo.

REM Backend tests
echo ==========================================
echo 📦 Running Backend Tests (Go)
echo ==========================================
cd backend
call run-tests.bat
if errorlevel 1 (
    echo.
    echo ❌ Backend tests failed!
    cd ..
    exit /b 1
)
cd ..
echo.

REM Frontend tests
echo ==========================================
echo 🎨 Running Frontend Tests (React/TypeScript)
echo ==========================================
cd frontend\gameTracker
call npm run coverage
if errorlevel 1 (
    echo.
    echo ❌ Frontend tests failed!
    cd ..\..
    exit /b 1
)
cd ..\..
echo.

echo ==========================================
echo ✅ All tests completed successfully!
echo ==========================================
