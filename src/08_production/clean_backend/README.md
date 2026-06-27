# Clean Architecture Backend in Go

> From first principles: build a production-ready REST API with JWT auth, MongoDB, and clean layering.

## Architecture

```
┌──────────────────────────────────────────────────────┐
│                    HTTP Request                        │
└──────────────┬───────────────────────────────────────┘
               ▼
┌──────────────────────────────────────────────────────┐
│              Router (Gin Engine)                       │
│   Public: /signup, /login, /refresh                   │
│   Protected: /profile, /task (JWT middleware)          │
└──────────────┬───────────────────────────────────────┘
               ▼
┌──────────────────────────────────────────────────────┐
│              Controller Layer                          │
│   Parses HTTP request → calls Usecase → HTTP response │
└──────────────┬───────────────────────────────────────┘
               ▼
┌──────────────────────────────────────────────────────┐
│              Usecase Layer (Business Logic)            │
│   Orchestrates domain logic with timeouts             │
└──────────────┬───────────────────────────────────────┘
               ▼
┌──────────────────────────────────────────────────────┐
│              Repository Layer (Data Access)            │
│   Talks to MongoDB through abstraction interfaces     │
└──────────────┬───────────────────────────────────────┘
               ▼
┌──────────────────────────────────────────────────────┐
│              Domain Layer                              │
│   Models, Interfaces, DTOs ZERO dependencies        │
└──────────────────────────────────────────────────────┘
```

## Key Principles

1. **Dependency Rule**: Dependencies point inward. Domain knows nothing about HTTP or DB.
2. **Interface Segregation**: Each usecase defines only the methods it needs.
3. **Testability**: MongoDB is abstracted behind interfaces → mock everything.
4. **Context Timeouts**: Every DB call respects `context.WithTimeout`.
5. **JWT Auth**: Access + Refresh token flow with middleware protection.

## Folder Structure

```
clean_backend/
├── cmd/main.go              # Entry point wires everything together
├── bootstrap/
│   ├── app.go               # Application struct (Env + DB client)
│   ├── env.go               # Viper-based config from .env
│   └── database.go          # MongoDB connection lifecycle
├── domain/                  # Pure domain: models + interfaces
│   ├── user.go
│   ├── task.go
│   ├── jwt_custom.go
│   ├── login.go
│   ├── signup.go
│   ├── profile.go
│   ├── refresh_token.go
│   ├── error_response.go
│   └── success_response.go
├── repository/              # Data access implementations
│   ├── user_repository.go
│   └── task_repository.go
├── usecase/                 # Business logic
│   ├── signup_usecase.go
│   ├── login_usecase.go
│   ├── profile_usecase.go
│   ├── refresh_token_usecase.go
│   └── task_usecase.go
├── api/
│   ├── controller/          # HTTP handlers
│   │   ├── signup_controller.go
│   │   ├── login_controller.go
│   │   ├── profile_controller.go
│   │   ├── refresh_token_controller.go
│   │   └── task_controller.go
│   ├── middleware/
│   │   └── jwt_auth_middleware.go
│   └── route/
│       ├── route.go         # Master router (public + protected)
│       ├── signup_route.go
│       ├── login_route.go
│       ├── profile_route.go
│       ├── refresh_token_route.go
│       └── task_route.go
├── mongo/                   # MongoDB abstraction (interfaces + wrappers)
│   └── mongo.go
├── internal/
│   └── tokenutil/
│       └── tokenutil.go     # JWT creation, validation, extraction
├── Dockerfile
├── docker-compose.yaml
├── .env.example
└── README.md
```

## How to Run

```bash
# Without Docker (needs local MongoDB)
cp .env.example .env
# Edit .env: set DB_HOST=localhost
go run cmd/main.go

# With Docker
cp .env.example .env
docker-compose up -d
```

## API Endpoints

| Method | Path      | Auth     | Description            |
|--------|-----------|----------|------------------------|
| POST   | /signup   | Public   | Register new user      |
| POST   | /login    | Public   | Login, get tokens      |
| POST   | /refresh  | Public   | Refresh access token   |
| GET    | /profile  | Bearer   | Get user profile       |
| POST   | /task     | Bearer   | Create a task          |
| GET    | /task     | Bearer   | List user's tasks      |

## What You'll Learn

- Clean Architecture in Go (Dependency Inversion, Interface Segregation)
- JWT access/refresh token flow
- MongoDB with interface abstraction for testability
- Gin framework routing and middleware
- Viper for configuration management
- Docker multi-stage builds
- Context timeout patterns
- Password hashing with bcrypt
