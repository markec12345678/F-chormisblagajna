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
- `/modules/` - business logic (core, hubsync, auth, fiscal modules)
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
- `golang.org/x/crypto/pkcs12` - PKCS#12 certificate parsing (fiscal module)

## Testing
- Backend tests (13 packages, 196 tests): `common/config`, `common/customerrors`, `common/helpers` (+ParsePagination, ParseAcceptLanguage), `common/logger`, `common/middlewares` (ratelimit +X-Forwarded-For), `common/userio` (+PropagateCounterIndexToTree), `modules/auth`, `modules/auth/middlewares` (jwt, bcrypt, auth), `modules/auth/models`, `modules/core/middlewares` (CORS), `modules/core/models`, `modules/core/dto`, `modules/fiscal/services` (ZOI, QR, format, tax grouping, receipt), `modules/fiscal_hr/services` (ZKI, OIB, XML, SOAP, XML-DSig, PDV grouping, receipt)
- Frontend tests (38 files, 230 tests): `frontend/src/__tests__/`
  - Components: ErrorBoundary, InventoryItem, Notification, Order, OrderItem, OrderItemView, AddCustomer, MealCard, PickMaterial, OrderItemRefund
  - Complex: QueueOrder, OrderView, StashedOrder
  - Pages: Settings, Categories, Customers, Orders, Kitchen, Languages, Profile, Home, Login, Sales, Products, Inventory, NoAccessView, AdminSetup, Setup, Hubsync
  - Not testable (uses `getCurrentInstance()`/`proxy.$auth`): Admin, Users
- Run backend: `$env:GOTOOLCHAIN = "auto"; & "C:\go\bin\go.exe" test -count=1 -timeout 120s ./...`
- Run frontend: `cd frontend && npx vitest run`
- Frontend lint: `cd frontend && npx eslint . --fix`
- Frontend format: `cd frontend && npx prettier --write src/`

## CI/CD
- GitHub Actions: `.github/workflows/ci.yml`
- Jobs: Go (vet+golangci-lint+test+build+vulncheck) → Vue (typecheck+lint+test+build) → Docker
- Docker: multi-stage build, non-root user, healthcheck, stripped binary

## Entities
- `Material`, `Component` and `Inventory Item` are the same entity
- `Product` and `Recipe` are the same entity

## Frontend i18n
- All UI strings use `$t()` in templates and `t()` in script via `useI18n()`
- Fallback keys defined in `frontend/src/main.ts` (130+ keys)
- Languages: SLO (primary), EN (fallback), AR (RTL)
- When adding new UI strings, add key to main.ts fallback first

## Frontend Components
- `ErrorBoundary.vue` — wraps RouterView in App.vue, catches child component errors with retry
- `InventoryItem.vue` — props interface (no self-recursion)
- `QueueOrder.vue`, `OrderView.vue`, `MainSearchResultView.vue` — order workflow
- `MealCard.vue` — product display with inventory consumption toggle

## API http schema
When calling the backend api from the frontend vue app, make sure to include the VITE_APP_BACKEND_HOST and VITE_APP_MODULE_CORE_API_PREFIX env vars in the request path

## Fiscal Module (FURS ZAPOS)
- `modules/fiscal/` — FURS ZAPOS fiscal receipt integration
- ZOI: RSA-SHA256 PKCS#1v1.5 signing → MD5 hash → 32-char hex
- QR data: 60 digits (ZOI decimal 39d + tax number 8d + datetime 12d + check digit 1d)
- API: mTLS (PKCS#12 certificate) + JWS RS256 signed requests
- Test env: `https://blagajne-test.fu.gov.si:9002`, prod: `https://blagajne.fu.gov.si:9003`
- Endpoints: `/fiscal/api/fiscal/echo`, `/fiscal/api/fiscal/invoice`, `/fiscal/api/fiscal/storno`, `/fiscal/api/fiscal/settings`
- Background worker: offline retry queue (every 5 min, exponential backoff up to 1 hour)
- Fiscal settings stored in core Settings model (`fiscal` field)
- i18n keys: fiscal, fiscal_description, fiscal_enabled, fiscal_environment, etc.

## Croatian Fiscal Module (Fina eRačun)
- `modules/fiscal_hr/` — Croatian CIS fiscalization (TR-2006 XML format)
- ZKI: RSA-SHA1 signing → MD5 hash → 32-char hex
- XML: TR-2006 `RacunZahtjev` with PDV (VAT) grouping, SOAP 1.1 envelope
- XML-DSig: Enveloped signature, Exclusive C14N, RSA-SHA1, SHA1 digest
- Protocol: SOAP XML web service (NOT REST), XML-DSig with PKCS#12 certificate
- Test env: `https://cistest.apis-it.hr:8449/FiskalizacijaServiceTest`, prod: `https://cis.porezna-uprava.hr:8449/FiskalizacijaService`
- Endpoints: `/fiscal_hr/api/fiscal_hr/echo`, `/fiscal_hr/api/fiscal_hr/invoice`, `/fiscal_hr/api/fiscal_hr/storno`, `/fiscal_hr/api/fiscal_hr/settings`
- Background worker: offline retry queue (every 5 min, max 10 retries)
- Fiscal settings stored in core Settings model (`fiscal_hr` field)
- i18n keys: fiscal_hr, fiscal_hr_description, fiscal_hr_enabled, fiscal_hr_oib, etc.
