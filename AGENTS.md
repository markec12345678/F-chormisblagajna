# Agent Guidance for NutrixPOS

## Project Overview
- Go 1.25 monorepo with MongoDB backend (mongo-driver v2.2.0)
- Point-of-sale system for restaurants/shops: inventory, sales, products
- Active development - no backward compatibility guarantee
- License: GPL v2

## Build Commands
```bash
go build ./...        # build all packages
go run ./cmd/pos      # run the CLI
go test ./...         # run all tests
go vet ./...          # static analysis
```

## Architecture
- `/cmd/` - CLI entrypoints (Cobra CLI)
- `/modules/` - business logic (core, hubsync, auth modules)
- `/common/` - shared utilities (database, config, logger, middlewares)
- `/frontend/` - Vue 3 SPA (separate build with Vite)

## Database
- Use `common.GetDatabaseClient()` singleton - never create new `mongo.Connect()` connections
- Singleton pattern in `common/database.go` ensures single connection
- **mongo-driver v2**: use `bson` package (NOT `bson/primitive`)
- **v2 Connect**: `mongo.Connect(opts)` - no context parameter

## Import Paths (mongo-driver v2)
```go
import (
    "go.mongodb.org/mongo-driver/v2/bson"              // ObjectID, NewObjectID, NilObjectID, etc.
    "go.mongodb.org/mongo-driver/v2/mongo"
    "go.mongodb.org/mongo-driver/v2/mongo/options"
)
```

## Common Pitfalls to Avoid

### 1. Database Connection Pattern (CRITICAL)
❌ Wrong:
```go
clientOptions := options.Client().ApplyURI(...)
ctx, cancel := context.WithTimeout(...)
client, err := mongo.Connect(ctx, clientOptions)
```
✅ Correct:
```go
client, err := common.GetDatabaseClient(logger, &config)
if err != nil {
    return err
}
ctx := context.Background()
```

### 2. bson Package (v2)
❌ Wrong:
```go
import "go.mongodb.org/mongo-driver/v2/bson/primitive"
id := primitive.NewObjectID()
```
✅ Correct:
```go
import "go.mongodb.org/mongo-driver/v2/bson"
id := bson.NewObjectID()
```

### 3. Update Options (v2)
❌ Wrong: `options.Update()`
✅ Correct: `options.UpdateOne()`

❌ Wrong: `&options.FindOptions{Sort: sort}`
✅ Correct: `options.Find().SetSort(sort)`

### 4. Error Handling in Services
❌ Wrong (CRASHES SERVER):
```go
client, err := common.GetDatabaseClient(...)
if err != nil {
    log.Fatal(err)  // or panic(err)
}
```
✅ Correct:
```go
client, err := common.GetDatabaseClient(...)
if err != nil {
    return fmt.Errorf("FunctionName: %w", err)
}
```

### 5. Regex Injection (ReDoS)
❌ Wrong:
```go
filter["name"] = bson.M{"$regex": fmt.Sprintf("(?i).*%s.*", userInput)}
```
✅ Correct:
```go
import "regexp"
filter["name"] = bson.M{"$regex": fmt.Sprintf("(?i).*%s.*", regexp.QuoteMeta(userInput))}
```

### 6. Error Response to Client
❌ Wrong (leaks internal errors):
```go
w.Write([]byte(err.Error()))
w.WriteHeader(http.StatusInternalServerError)
```
✅ Correct:
```go
http.Error(w, "failed to get data", http.StatusInternalServerError)
```

## Security Practices
- JWT secret auto-generated if empty (crypto/rand)
- Rate limiting: sliding window on auth endpoints
- User input escaped with `regexp.QuoteMeta` in MongoDB queries
- Error messages to clients are generic (no internal details)
- Docker runs as non-root user
- Passwords: bcrypt hashed
- Tokens: crypto/rand (not math/rand)

## Dependencies
- `go.mongodb.org/mongo-driver/v2` - MongoDB driver v2
- `github.com/gorilla/mux` - HTTP router
- `github.com/spf13/cobra` + `viper` - CLI framework
- `github.com/golang-jwt/jwt/v5` - JWT authentication
- `github.com/nutrixpos/crypt` - password hashing
- `github.com/nutrixpos/melody` - WebSocket

## Testing
- Tests in: `common/config`, `common/middlewares`, `modules/auth/middlewares`
- Rate limiter: `common/middlewares/ratelimit.go` with sliding window
- Frontend tests: `frontend/src/__tests__/`
- Run: `go test -race ./...`

## CI/CD
- GitHub Actions: `.github/workflows/ci.yml`
- Jobs: Go (vet+golangci-lint+test+build+vulncheck) → Vue (typecheck+lint+test+build) → Docker
- Docker: multi-stage build, non-root user, healthcheck, stripped binary

## Entities
- `Material`, `Component` and `Inventory Item` are the same entity
- `Product` and `Recipe` are the same entity

## API http schema
When calling the backend api from the frontend vue app, make sure to include the VITE_APP_BACKEND_HOST and VITE_APP_MODULE_CORE_API_PREFIX env vars in the request path
