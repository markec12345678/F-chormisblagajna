<p align="center">
  <img src="https://elmawardy.sirv.com/nutrixpos-docs/nutrixdocs0.png" alt="NutrixPOS" height="400" />
</p>

<p align="center">
  <a href="https://pkg.go.dev/github.com/nutrixpos/pos"><img src="https://pkg.go.dev/badge/github.com/nutrixpos/pos.svg" alt="Go Reference"></a>
  <a href="https://github.com/markec12345678/F-chormisblagajna/actions"><img src="https://github.com/markec12345678/F-chormisblagajna/actions/workflows/ci.yml/badge.svg" alt="CI"></a>
  <a href="https://nutrixpos.com/userguide/installation.html"><img src="https://img.shields.io/badge/docs-nutrixpos.com-teal" alt="Docs"></a>
  <a href="LICENSE"><img src="https://img.shields.io/badge/license-GPL--v2-blue" alt="License"></a>
</p>

# NutrixPOS

Celovit blagajniški sistem za restavracije in trgovine. **Go backend (MongoDB) + Vue 3 SPA frontend** z modularno arhitekturo — 48 modulov, ki jih poljubno vklapljaš.

## Pregled modulov

| Kategorija | Moduli |
|------------|--------|
| **POS & Naročila** | Core Orders, Multi-Payment, Split Payment, Online Orders, Delivery Management, Tableside Ordering, Self-Service Kiosk |
| **Kuhinja** | Kitchen Display (WebSocket), Menu Engineering, Receipt Customization |
| **Stranke** | Customer CRUD, Customer Feedback, Loyalty, Gift Cards |
| **Inventar & Nabava** | Inventory (materials), Inventory Alerts, Inventory Transfers, Purchase Orders, Suppliers, Waste Tracking |
| **Finance** | Accounting, Expenses, Employee Tips, FURS ZAPOS (SI fiscal), CIS (HR fiscal) |
| **Osebje** | Auth & Roles (JWT/OIDC/NoAuth), Employee Performance, Scheduling, Time Clock, Staff Chat, Staff Training |
| **Prostor** | Floor Plan (visual table editor), Reservations, Queue/Waitlist |
| **Marketing** | Marketing Campaigns, Promotions (auto-discount rules) |
| **Operativno** | Multi-Location Dashboard, Audit Log, Reports (CSV export), Notifications, HubSync (remote sync), Languages (i18n) |
| **Nastavitve** | App Settings, Admin Setup Wizard |

## Tech Stack

| Komponenta | Tehnologija |
|------------|-------------|
| Backend | Go 1.25 + gorilla/mux + Cobra CLI |
| Baza | MongoDB (mongo-driver v2.2.0) |
| Frontend | Vue 3 + TypeScript + PrimeVue 4 |
| Build | Vite 6 + viteSingleFile |
| State | Pinia |
| Auth | JWT (HS256) / Zitadel OIDC / NoAuth |
| i18n | vue-i18n (460+ ključev, SLO/EN) |
| CI/CD | GitHub Actions (golangci-lint) |
| Container | Docker (multi-stage, non-root) |
| Linting | ESLint + Prettier + oxlint |
| Testi | vitest (frontend) + go test (backend) |

## Hitri začetek

### Docker (priporočeno)

```bash
docker compose up
```

Ob prvem zagonu bo Setup Wizard na `http://localhost:8000` za nastavitev MongoDB povezave.

### Ročno

**Backend:**
```bash
# Potrebujete Go 1.25+ in MongoDB
cp config.example.yaml config.yaml
go run ./cmd/pos
```

**Frontend:**
```bash
cd frontend
npm install
npm run dev
```

Frontend bo dostopen na `http://localhost:3000`.

## Konfiguracija

Kopiraj `config.example.yaml` v `config.yaml`:

```yaml
databases:
  - name: core
    type: mongo
    host: 127.0.0.1
    port: 27017
    database: nutrix

auth:
  jwt_secret: ""          # prazen = naključno ob zagonu
  jwt_expire_hrs: 24
  enabled: true

zitadel:
  enabled: false          # vključi za OIDC

serve_frontend: true      # streži Vue SPA iz Go backend-a
```

## Razvoj

### Build

```bash
# Backend
go build -o pos ./cmd/pos

# Frontend
cd frontend && npm run build
```

### Testi

```bash
# Backend (13+ paketov, 196+ testov)
go test -race ./...

# Frontend (38 testnih datotek, 230+ testov)
cd frontend && npx vitest run
```

### CI/CD

GitHub Actions pipeline (`.github/workflows/ci.yml`):
- **Backend**: vet, golangci-lint, test (race + coverage), build, govulncheck
- **Frontend**: type-check, lint (ESLint + oxlint), test (vitest), build
- **Docker**: multi-stage build + cache (odvisen od backend + frontend)

## Arhitektura

```
main.go
  └── cmd/root.go (Cobra CLI + gorilla/mux)
       ├── Core modul
       │    ├── Auth (JWT/Zitadel/NoAuth)
       │    ├── Products/Recipes CRUD
       │    ├── Materials/Inventory CRUD
       │    ├── Orders (submit→start→finish→pay/cancel/refund)
       │    ├── Customers CRUD + stats
       │    ├── Categories CRUD
       │    ├── Sales reports + CSV export
       │    ├── Disposals + Waste
       │    ├── Settings
       │    ├── Languages (i18n)
       │    ├── WebSocket (real-time kitchen display)
       │    └── Background workers (fiscal retry)
       ├── HubSync modul
       │    └── Sinhronizacija z hub strežnikom
       ├── Fiscal modul (FURS ZAPOS)
       │    ├── ZOI, EOR, QR generation
       │    ├── mTLS + JWS RS256
       │    └── Offline retry worker
       ├── Fiscal HR modul (CIS eRačun)
       │    ├── ZKI, XML-DSig, SOAP
       │    └── Offline retry worker
       └── 44 dodatnih modulov (vsak sledi IBaseModule)
            ├── handlers/ → API endpoints
            ├── services/ → business logic
            └── models/   → data structures
```

Vsak modul implementira `IBaseModule` vmesnik (`OnStart`/`OnEnd`), se registrira v `cmd/root.go` prek `LoadModule()` in samodejno dobi:
- JWT avtentikacijo (ali NoAuth za javne endpoint-e)
- MongoDB dostop prek `common.GetDatabaseClient()`
- Cobra CLI ukaz

## Varnost

- **JWT secret** — samodejna generacija naključnega ključa ob zagonu (`crypto/rand`)
- **Rate limiting** — drseno okno na auth endpointih (10 prijav/min, 5 registracij/min)
- **NoSQL injection** — uporabniški vnos ekraniran z `regexp.QuoteMeta`
- **Error handling** — notranje napake niso izpostavljene strankam
- **Docker** — non-root uporabnik, healthcheck, strip binary
- **Kriptografija** — `crypto/rand` za varne tokene, bcrypt za gesla
- **PIN kode** — SHA-256 salted hash

## Dostopne vloge

| Vloga | Dovoljenja |
|-------|-----------|
| `superuser` | Vse |
| `admin` | Administracija, naročila, nastavitve |
| `cashier` | Blagajna, prodaja, stranke |
| `chef` | Kuhinjski zaslon, pregled materialov |

## Datoteke

| Datoteka | Namen |
|----------|-------|
| `README.md` | Ta datoteka |
| `LICENSE` | GNU GPL v2 |
| `AGENTS.md` | Navodila za AI agente |
| `.golangci.yml` | Konfiguracija lintinga |
| `.github/workflows/ci.yml` | CI/CD pipeline |
| `Dockerfile` | Multi-stage build |
| `docker-compose.yaml` | POS + frontend + MongoDB |
| `frontend/eslint.config.ts` | ESLint + Prettier + oxlint |
| `frontend/vitest.config.ts` | Vitest konfiguracija |

## Licenca

GNU General Public License v2 — glej [LICENSE](LICENSE)

> **Opozorilo:** NutrixPOS je v aktivnem razvoju. Nazaj-kompatibilnost ni zagotovljena do stabilne izdaje.
