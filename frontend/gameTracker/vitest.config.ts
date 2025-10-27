import { defineConfig } from 'vitest/config'
import react from '@vitejs/plugin-react'

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
      thresholds: {
        lines: 70,
        branches: 70,
        functions: 70,
        statements: 70
      }
    },
    globals: true,
    setupFiles: ['./src/test/setup.ts']
  }
})
