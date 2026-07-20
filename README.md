# finance-go

Base de una API financiera en Go con Gin, GORM, PostgreSQL, golang-migrate, Viper, godotenv, validator y JWT.

## Requisitos

- Go 1.26+
- Docker y Docker Compose
- `just`

## Desarrollo

1. Copiá `.env.example` como `.env` y cambiá `JWT_SECRET` por un valor aleatorio.
2. Levantá PostgreSQL con `docker compose up -d postgres`.
3. Ejecutá `just migrate-up`.
4. Generá la documentación OpenAPI con `just docs`.
5. Inicializá la API con `just run`.

La documentación interactiva queda disponible en `http://localhost:8080/docs` y el esquema JSON en `http://localhost:8080/docs/openapi.json`. La ruta `GET /health` no requiere autenticación. La ruta `GET /api/v1/me` requiere `Authorization: Bearer <token>` y muestra la integración del middleware JWT. La emisión de tokens queda encapsulada en `internal/auth.TokenService` para conectarla al flujo real de login.

## Verificación

```sh
just test
just vet
```
