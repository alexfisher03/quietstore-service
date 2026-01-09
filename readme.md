# QuietStore

QuietStore is a self-hosted file storage service designed for learning modern backend architecture. It combines a Go (Fiber) API, PostgreSQL for metadata, MinIO for object storage, and Docker Compose for local orchestration. The project is intentionally structured to evolve toward Kubernetes-based deployment in the future, while remaining simple enough to run locally today.

This README focuses on **local development and usage** on macOS and Windows.

---

## High-Level Architecture

* **API Service**: Go + Fiber
* **Database**: PostgreSQL (file metadata, users, refresh tokens)
* **Object Storage**: MinIO (S3-compatible)
* **Reverse Proxy / TLS (optional, dev-only)**: Caddy
* **Orchestration**: Docker Compose

The API never serves files directly from disk. All file contents live in MinIO; the API only streams data and enforces access control.

---

## Prerequisites

### Required

* Docker Desktop

  * macOS: Docker Desktop for Mac
  * Windows: Docker Desktop with WSL2 enabled
* Git

### Optional (for development)

* Go 1.22+
* Insomnia / Postman (API testing)

---

## Repository Layout

```
quietstore-service/
├── docker-compose.yaml
├── Caddyfile                # optional (local HTTPS)
├── certs/
│   ├── minio/               # MinIO TLS certs
│   └── caddy/               # Caddy internal CA data
└── QuietStore/
    ├── cmd/server            # main entrypoint
    ├── api/v1                # HTTP routes
    ├── internal/
    │   ├── handlers          # request handlers
    │   ├── repo              # DB access layer
    │   ├── service           # storage + business logic
    │   ├── models            # domain models
    │   ├── db                # migrations + DB init
    │   └── config            # env config loading
```

---

## Environment Configuration

Create a file at:

```
QuietStore/.env
```

Example:

```
# Server
SERVER_PORT=8080
SERVER_HOST=0.0.0.0

# Database
DB_DSN=postgres://quietstore:quietstore@db:5432/quietstore?sslmode=disable
POSTGRES_USER=quietstore
POSTGRES_PASSWORD=quietstore
POSTGRES_DB=quietstore

# Auth
AUTH_JWT_SECRET=dev-secret-change-me
AUTH_ACCESS_TTL_MIN=10
AUTH_REFRESH_TTL_DAYS=7

# MinIO
MINIO_ENDPOINT=http://minio:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_BUCKET=quietstore
MINIO_USE_SSL=false
```

---

## Running the System (macOS / Windows)

From the project root:

```
docker compose up --build
```

Services started:

* API: [http://localhost:8080](http://localhost:8080)
* PostgreSQL: localhost:5432
* MinIO API: [http://localhost:9000](http://localhost:9000)
* MinIO Console: [http://localhost:9001](http://localhost:9001)

On first boot:

* DB migrations run automatically
* MinIO bucket is created if missing

---

## Health Check

```
GET /api/v1/health
```

Example:

```
{
  "service": "QuietStore",
  "status": "ok",
  "timestamp": "2025-08-26T05:53:20Z"
}
```

---

## Authentication Flow

### Login

```
POST /api/v1/auth/login
```

Body:

```
{
  "username": "alice",
  "password": "password123"
}
```

Response:

* `access_token` (JWT, short-lived)
* `refresh_token` (opaque, long-lived)

### Access Token Usage

For authenticated requests:

```
Authorization: Bearer <access_token>
```

### Refresh Token

```
POST /api/v1/auth/refresh
```

Used when the access token expires. Refresh tokens are:

* Stored **hashed** in the database
* Single-use (rotated on refresh)
* Revoked on logout

### Logout

```
POST /api/v1/auth/logout
```

Requires:

* Authorization header
* Refresh token in body

---

## File Operations

All file routes are scoped to the authenticated user.

### Upload

```
POST /api/v1/me/files/upload
```

Multipart form field:

* `file`

### List Files

```
GET /api/v1/me/files?limit=50&offset=0
```

### Download

```
GET /api/v1/me/files/:fileID
```

### Delete

```
DELETE /api/v1/me/files/:fileID
```

The backend enforces ownership at every step; users cannot access files they do not own.

---

## Rate Limiting

Rate limiting is applied per-route group using Fiber middleware:

* Auth endpoints (login / refresh / logout)
* User management endpoints
* File upload and delete endpoints

Limits are intentionally conservative to protect local hardware.

---

## MinIO Usage

MinIO is used as an S3-compatible object store.

* Buckets are created automatically
* Objects are keyed by internal file IDs
* Metadata lives in PostgreSQL

You can inspect objects via the MinIO Console:

```
http://localhost:9001
```

---

## HTTPS (Development Only)

Caddy is included **only for local experimentation** with TLS.

* Uses an internal CA
* Requires trusting a local root certificate
* Not required for normal development

In production (or Kubernetes), TLS will be terminated by:

* Ingress controller
* Cloud load balancer
* Reverse proxy

You can safely ignore HTTPS locally if you prefer.

---

## Swagger / API Docs (Planned)

Swagger support can be added using:

[https://github.com/gofiber/swagger](https://github.com/gofiber/swagger)

This will expose:

```
GET /swagger/index.html
```
