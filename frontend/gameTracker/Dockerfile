# ---------- STAGE 1: Build de Vite/React ----------
    FROM node:20-alpine AS builder

    WORKDIR /app
    
    # Copiar archivos de dependencias
    COPY package*.json ./
    
    # Instalar dependencias
    RUN npm ci
    
    # Copiar el código fuente
    COPY . .
    
    # Inyectar variable de entorno para la API (Render)
    ENV VITE_API_URL="https://gametracker-backend-rf2e.onrender.com"
    
    # Build de producción
    RUN npm run build
    
    # ---------- STAGE 2: nginx para servir los archivos ----------
    FROM nginx:alpine
    
    # Limpiar archivos por defecto
    RUN rm -rf /usr/share/nginx/html/*
    
    # Copiar el build del frontend
    COPY --from=builder /app/dist /usr/share/nginx/html
    
    # Exponer el puerto estándar
    EXPOSE 80
    
    # Iniciar nginx (config por defecto, sin upstreams)
    CMD ["nginx", "-g", "daemon off;"]
    