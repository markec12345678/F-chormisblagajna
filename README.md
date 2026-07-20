<p align="center">
  <img src="https://elmawardy.sirv.com/nutrixpos-docs/nutrixdocs0.png" alt="NutrixPOS" height="400" />
</p>

<p align="center">
  <a href="https://pkg.go.dev/github.com/nutrixpos/pos"><img src="https://pkg.go.dev/badge/github.com/nutrixpos/pos.svg" alt="Go Reference"></a>
  <a href="https://github.com/markec12345678/F-chormisblagajna/actions"><img src="https://github.com/markec12345678/F-chormisblagajna/actions/workflows/ci.yml/badge.svg" alt="CI"></a>
  <a href="https://nutrixpos.com/userguide/installation.html"><img src="https://img.shields.io/badge/docs-nutrixpos.com-teal" alt="Docs"></a>
</p>

# NutrixPOS

Blagajniški sistem za restavracije in trgovine. Go backend (MongoDB) + Vue 3 SPA frontend.

## Funkcionalnosti

- 🍳 **Kuhinjski zaslon** — real-time naročila z WebSocket
- 📦 **Inventar** — materiali, zaloge, dnevnik porabe
- 📋 **Recepture** — produkti z materiali, pod-recepture
- 🛒 **Naročila** — submit, start, finish, pay, cancel, refund, tips
- 👥 **Stranke** — CRUD zgodovina naročil
- 🔄 **HubSync** — sinhronizacija z oddaljenim strežnikom
- 🌍 **Več jezikov** — dinamični paketi (SLO, EN, AR)
- 🖨️ **Thermal printer** — Handlebars template za ESC/POS
- 🔐 **Auth** — JWT / Zitadel OIDC / NoAuth
- 👤 **Več vlog** — superuser, admin, cashier, chef
- 🌙 **Dark mode** + **RTL** podpora

## Tech Stack

| Komponenta | Tehnologija |
|------------|-------------|
| Backend | Go 1.25 + gorilla/mux + Cobra CLI |
| Baza | MongoDB |
| Frontend | Vue 3 + TypeScript + PrimeVue 4 |
| Build | Vite 6 + viteSingleFile |
| State | Pinia |
| Auth | JWT (HS256) / Zitadel OIDC |
| CI/CD | GitHub Actions |
| Container | Docker (multi-stage) |

## Hitri začetek

### Docker (priporočeno)

```bash
docker compose up
```

Ob prvič bo Setup Wizard zagnas na `http://localhost:8000` za nastavitev MongoDB povezave.

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
# Backend
go test -race ./...

# Frontend
cd frontend && npm run test:unit
```

### CI/CD

GitHub Actions pipeline (.github/workflows/ci.yml):
- **Backend**: vet, test, build, govulncheck
- **Frontend**: type-check, lint, test, build
- **Docker**: multi-stage build + cache

## Arhitektura

```
main.go
  └── cmd/root.go (Cobra CLI + gorilla/mux)
       ├── Core modul
       │    ├── Auth (JWT/Zitadel/NoAuth)
       │    ├── Products/Recipes CRUD
       │    ├── Materials/Inventory CRUD
       │    ├── Orders (submit→start→finish→pay)
       │    ├── Customers CRUD
       │    ├── Categories CRUD
       │    ├── Sales reports + CSV export
       │    ├── Disposals
       │    ├── Settings
       │    ├── Languages (i18n)
       │    ├── WebSocket (real-time)
       │    └── Background workers
       └── HubSync modul
            └── Sinhronizacija z hub strežnikom
```

## Dostopne vloge

| Vloga | Dovoljenja |
|-------|-----------|
| `superuser` | Vse |
| `admin` | Administracija, naročila, nastavitve |
| `cashier` | Blagajna, prodaja, stranke |
| `chef` | Kuhinjski zaslon, pregled materialov |

## Licenca

GNU General Public License v2 - glej [LICENSE](LICENSE)

> **Opozorilo:** NutrixPOS je v aktivnem razvoju. Nazaj-kompatibilnost ni zagotovljena do stabilne izdaje.
