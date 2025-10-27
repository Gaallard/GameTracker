#!/bin/bash
# Script para ejecutar todas las pruebas (backend y frontend)

echo "=========================================="
echo "🧪 Running All Tests (Backend + Frontend)"
echo "=========================================="
echo ""

# Backend tests
echo "=========================================="
echo "📦 Running Backend Tests (Go)"
echo "=========================================="
cd backend
bash run-tests.sh
if [ $? -ne 0 ]; then
    echo ""
    echo "❌ Backend tests failed!"
    cd ..
    exit 1
fi
cd ..
echo ""

# Frontend tests
echo "=========================================="
echo "🎨 Running Frontend Tests (React/TypeScript)"
echo "=========================================="
cd frontend/gameTracker
npm run coverage
if [ $? -ne 0 ]; then
    echo ""
    echo "❌ Frontend tests failed!"
    cd ../..
    exit 1
fi
cd ../..
echo ""

echo "=========================================="
echo "✅ All tests completed successfully!"
echo "=========================================="
