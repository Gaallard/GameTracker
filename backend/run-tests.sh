#!/bin/bash

# Script para ejecutar pruebas unitarias del backend con cobertura

echo "🔧 Instalando herramientas de testing..."
go install gotest.tools/gotestsum@latest
go install github.com/jstemmer/go-junit-report/v2@latest

echo "🧪 Ejecutando pruebas unitarias con cobertura..."
gotestsum --format=standard-verbose --junitfile test-results-go.xml -- \
  -covermode=atomic -coverprofile=coverage.out ./service ./controller

echo "📊 Generando reportes de cobertura..."
go tool cover -func=coverage.out > coverage.txt
go tool cover -html=coverage.out -o coverage.html

echo "✅ Verificando cobertura mínima del 70%..."
awk '/total:/ { gsub("%","",$3); if ($3+0 < 70) { print "❌ Coverage below 70%: " $3"%"; exit 1 } else { print "✅ Coverage: " $3"%"; exit 0 } }' coverage.txt

echo "📁 Archivos generados:"
echo "  - test-results-go.xml (reporte JUnit)"
echo "  - coverage.out (datos de cobertura)"
echo "  - coverage.txt (resumen de cobertura)"
echo "  - coverage.html (reporte HTML de cobertura)"
