# Docker Commands for GameTracker Multi-Environment Setup

## Environment Management

### Start QA Environment Only
```bash
docker-compose up db-qa backend-qa frontend-qa -d
```

### Start PROD Environment Only
```bash
docker-compose up db-prod backend-prod frontend-prod -d
```

### Start Both Environments
```bash
docker-compose up -d
```

### Stop QA Environment
```bash
docker-compose stop db-qa backend-qa frontend-qa
```

### Stop PROD Environment
```bash
docker-compose stop db-prod backend-prod frontend-prod
```

### Stop All Environments
```bash
docker-compose down
```

## Access Points

### QA Environment
- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8081
- **Database**: localhost:3307

### PROD Environment
- **Frontend**: http://localhost:8080
- **Backend API**: http://localhost:8082
- **Database**: localhost:3308

## Database Management

### Connect to QA Database
```bash
docker-compose exec db-qa mysql -u root -proot -e "USE gametracker_qa; SHOW TABLES;"
```

### Connect to PROD Database
```bash
docker-compose exec db-prod mysql -u root -proot -e "USE gametracker_prod; SHOW TABLES;"
```

### View QA Database Data
```bash
docker-compose exec db-qa mysql -u root -proot -e "USE gametracker_qa; SELECT COUNT(*) FROM games;"
```

### View PROD Database Data
```bash
docker-compose exec db-prod mysql -u root -proot -e "USE gametracker_prod; SELECT COUNT(*) FROM games;"
```

## Logs

### View QA Backend Logs
```bash
docker-compose logs backend-qa -f
```

### View PROD Backend Logs
```bash
docker-compose logs backend-prod -f
```

### View QA Frontend Logs
```bash
docker-compose logs frontend-qa -f
```

### View PROD Frontend Logs
```bash
docker-compose logs frontend-prod -f
```

## Rebuild Services

### Rebuild QA Backend
```bash
docker-compose up --build backend-qa -d
```

### Rebuild PROD Backend
```bash
docker-compose up --build backend-prod -d
```

### Rebuild QA Frontend
```bash
docker-compose up --build frontend-qa -d
```

### Rebuild PROD Frontend
```bash
docker-compose up --build frontend-prod -d
```

## Volume Management

### Remove QA Data (WARNING: This will delete all QA data)
```bash
docker-compose down -v
docker volume rm gametracker_db_qa_volume
```

### Remove PROD Data (WARNING: This will delete all PROD data)
```bash
docker-compose down -v
docker volume rm gametracker_db_prod_volume
```

## Environment Variables

### QA Environment Variables (env.qa)
- ENVIRONMENT=qa
- LOG_LEVEL=debug
- GIN_MODE=debug
- DB_NAME=gametracker_qa
- API_PORT=8080
- FRONTEND_PORT=3000

### PROD Environment Variables (env.prod)
- ENVIRONMENT=prod
- LOG_LEVEL=info
- GIN_MODE=release
- DB_NAME=gametracker_prod
- API_PORT=8080
- FRONTEND_PORT=80
