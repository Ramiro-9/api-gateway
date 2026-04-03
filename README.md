# API Gateway

Gateway centralizado construido en Go con Gin. Actúa como punto de entrada único para múltiples servicios, manejando autenticación, rate limiting, circuit breaker, cache, métricas y load balancing.

## Stack

- **Go** + **Gin** — servidor HTTP
- **JWT** — validación centralizada de tokens
- **APScheduler** — no aplica (Go nativo)

## Características

- Enrutamiento y proxy inverso hacia múltiples servicios
- Autenticación centralizada con JWT — los servicios no necesitan validar tokens
- Rate limiting por IP — 100 requests por minuto
- Circuit breaker — si un servicio falla 5 veces, se bloquea 30 segundos
- Timeout configurable por servicio
- Load balancing round-robin entre múltiples instancias
- Cache de respuestas GET con TTL configurable
- Headers de seguridad en todas las respuestas
- API keys internas — identifica que los requests vienen del gateway
- Métricas en tiempo real — requests, errores, latencia promedio
- Logging en consola y archivo

## Servicios registrados

| Servicio | URL | Timeout |
|----------|-----|---------|
| auth-api | http://localhost:8080 | 10s |
| crypto-etl | http://localhost:8000 | 30s |

## Endpoints del gateway

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/health` | Estado del gateway |
| GET | `/services` | Estado de los servicios y circuit breakers |
| GET | `/metrics` | Métricas en tiempo real |
| ANY | `/auth/*` | Proxy público → auth-api |
| ANY | `/pipeline/*` | Proxy protegido → crypto-etl |
| ANY | `/scheduler/*` | Proxy protegido con cache → crypto-etl |
| ANY | `/admin/*` | Proxy solo admin → crypto-etl |

## Instalación

### Requisitos

- Go 1.21+
- auth-api corriendo en puerto 8080
- crypto-etl corriendo en puerto 8000

### Pasos

1. Cloná el repositorio
```bash
git clone https://github.com/Ramiro-9/api-gateway.git
cd api-gateway
```

2. Instalá las dependencias
```bash
go mod tidy
```

3. Creá el archivo `.env` en la raíz
```env
GATEWAY_PORT=9000
JWT_SECRET=el_mismo_secreto_de_tu_auth_api

AUTH_API_URL=http://localhost:8080
CRYPTO_ETL_URL=http://localhost:8000

INTERNAL_API_KEY=un_secreto_interno_largo
```

4. Correlo
```bash
go run ./cmd/
```

## Arquitectura
```
Cliente
   │
   ▼
API Gateway :9000
   ├── Rate Limiting
   ├── Security Headers
   ├── JWT Validation
   ├── Cache (GET)
   ├── Metrics Tracking
   │
   ├── /auth/*     → auth-api :8080
   ├── /pipeline/* → crypto-etl :8000
   ├── /scheduler/*→ crypto-etl :8000
   └── /admin/*    → crypto-etl :8000
```
