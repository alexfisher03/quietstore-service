# QuietStore Development Plan

## Project Overview

Building a custom cloud storage and hosting platform to learn system design fundamentals including Docker, Kubernetes, logging/monitoring, and Go development.

**Hardware**: Ubuntu Desktop, 64GB RAM, NVIDIA P5000 Quadro, Ryzen 5 5600

## Technology Stack

### Frontend
- **React (TypeScript)**
  - Nginx serves statically
  - CRUD operations to backend
  - File uploads via multipart/form-data
  - Downloads via blob links
  - JSON responses from API

### Backend
- **Fiber (Go Web Framework)**
  - HTTP routing and middleware
  - JWT authentication
  - File storage management
  - Streaming files to/from disk
  - Goroutines for concurrency
  - Key endpoints:
    - `POST /login` → JWT authentication
    - `GET /files/` → list user files
    - `POST /upload` → file upload
    - `DELETE /files/:name` → file deletion

### Storage & Database
- **MinIO** (Self-hosted S3-compatible object storage)
- **PostgreSQL** (File metadata and user data)
- **Redis** (JWT token management)

### Infrastructure
- **Docker & Docker Compose**
- **Kubernetes (K3s)** for orchestration
- **Nginx** as reverse proxy

## Project Structure

```
QuietStore/
├── api/
│   └── v1/
│       └── routes.go          # API route definitions & grouping
├── cmd/
│   └── server/
│       └── main.go            # Application entry point
├── internal/                  # Private application code
│   ├── config/
│   │   └── config.go          # Configuration management
│   ├── handlers/              # HTTP request handlers
│   │   ├── auth.go           # Login, register endpoints
│   │   ├── files.go          # Upload, download, list endpoints
│   │   └── health.go         # Health check endpoint
│   ├── middleware/           # Cross-cutting concerns
│   │   ├── auth.go           # JWT validation & user context
│   │   └── cors.go           # CORS headers
│   ├── models/               # Data structures
│   │   ├── file.go           # File metadata struct
│   │   └── user.go           # User struct & validation
│   └── service/              # Business logic layer
│       ├── auth.go           # Authentication service
│       └── storage.go        # File operations service
├── pkg/
│   └── utils/                # Reusable utilities
├── go.mod
├── go.sum
└── main.go                   # Legacy - will be moved to cmd/
```

## Architecture Principles

### Dependency Flow
```
main.go → config → services → handlers → models
   ↓         ↓         ↓         ↓
   └─────────┴─────────┴─────────┴──→ middleware
```

### Layer Responsibilities

**1. `cmd/server/main.go` - The Orchestrator**
- Creates configuration
- Initializes services (auth, storage)
- Dependency injection to handlers
- Middleware setup
- Server startup
- **Keep to ~30-50 lines maximum**

**2. `internal/config/` - Configuration Management**
- Environment variable reading
- Default value provision
- Configuration validation
- Returns Config struct for other packages

**3. `internal/service/` - Business Logic Layer**
- Domain-specific operations
- Business rule enforcement
- Returns domain errors (not HTTP errors)
- **No HTTP knowledge** - pure Go types
- `auth.go`: User authentication, JWT operations
- `storage.go`: File operations, metadata management

**4. `internal/handlers/` - HTTP Layer**
- HTTP request/response translation
- Service layer integration
- HTTP status code management
- Request parsing & response formatting
- **Thin layer** - minimal business logic

**5. `internal/models/` - Data Structures**
- Pure data definitions
- JSON serialization tags
- Validation methods
- No business logic

**6. `internal/middleware/` - Cross-Cutting Concerns**
- Request preprocessing
- Authentication checks
- CORS handling
- Can short-circuit requests

**7. `api/v1/routes.go` - Route Organization**
- URL to handler mapping
- Route grouping
- Middleware application
- API versioning

## Development Phases

### Phase 1: Skeleton Structure ✅
- [x] Basic Go setup with Fiber
- [x] Understanding pointers, contexts, error handling
- [x] Middleware concepts and implementation
- [ ] **NEXT**: Create proper file structure with stub functions
- [ ] Wire dependencies together
- [ ] Ensure compilation and basic route testing

### Phase 2: Core Application Logic
- [ ] Implement `internal/config/` (environment variables)
- [ ] Build `service/auth.go` (JWT generation/validation)
- [ ] Create corresponding `handlers/auth.go`
- [ ] Add authentication middleware
- [ ] File upload/download logic in `service/storage.go`
- [ ] File handlers implementation

### Phase 3: Storage Integration
- [ ] PostgreSQL setup and connection pooling
- [ ] User model with password hashing
- [ ] File metadata models
- [ ] MinIO integration for object storage
- [ ] Redis for token management

### Phase 4: Infrastructure
- [ ] Containerization with Docker
- [ ] Docker Compose for local development
- [ ] K3s setup and deployment
- [ ] Nginx reverse proxy configuration

### Phase 5: Frontend & Production
- [ ] React TypeScript application
- [ ] File upload/download UI
- [ ] Authentication flows
- [ ] Production monitoring and logging

## Key Go Concepts in Use

### Dependency Injection Pattern
```go
// Services depend on config
authService := service.NewAuthService(config)

// Handlers depend on services  
authHandler := handlers.NewAuthHandler(authService)
```

### Interface-Based Design
```go
type StorageService interface {
    SaveFile(filename string, data []byte) error
    GetFile(filename string) ([]byte, error)
}
```

### Error Handling Patterns
```go
// Services return domain errors
if err := authService.ValidateToken(token); err != nil {
    return handlers.BadRequest("Invalid token")
}
```

## Learning Objectives

**System Design Concepts:**
- Microservice architecture patterns
- Service communication and discovery
- Data consistency and transaction management
- Scalability and load balancing

**DevOps & Infrastructure:**
- Container orchestration with Kubernetes
- Service mesh concepts
- Monitoring and observability
- CI/CD pipeline design

**Go Programming:**
- Concurrent programming with goroutines
- Interface design and dependency injection
- Error handling best practices
- Performance optimization techniques

## Implementation Notes

**Current Status**: Completed basic Fiber setup with middleware understanding. Ready to implement proper file structure.

**Next Steps**: Create directory structure with stub functions, focus on making everything compile and run before adding real implementation.

**Key Principle**: Build incrementally, test each layer before adding the next. Understand the "why" behind each architectural decision, not just the "how".