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

Celovit blagajniški sistem za restavracije in trgovine. **Go backend (MongoDB) + Vue 3 SPA frontend** z modularno arhitekturo — 45 modulov, ki jih poljubno vklapljaš.

---

## Screenshots

<p align="center">
  <img src="https://elmawardy.sirv.com/nutrixpos-docs/nutrixdocs0.png" alt="Dashboard" width="45%" />
  <img src="https://elmawardy.sirv.com/nutrixpos-docs/nutrixdocs0.png" alt="Orders" width="45%" />
</p>
<p align="center">
  <img src="https://elmawardy.sirv.com/nutrixpos-docs/nutrixdocs0.png" alt="Kitchen Display" width="45%" />
  <img src="https://elmawardy.sirv.com/nutrixpos-docs/nutrixdocs0.png" alt="Floor Plan" width="45%" />
</p>

> Slike so placeholder — zamenjaj z dejanskimi posnetki zaslona.

---

## Moduli

| Kategorija | Moduli |
|------------|--------|
| **POS & Naročila** | Core Orders, Multi-Payment, Split Bill, Online Orders, Delivery Management, Tableside Ordering, Self-Service Kiosk |
| **Kuhinja** | Kitchen Display (WebSocket), Menu Engineering (profitability matrix), Receipt Customization |
| **Stranke** | Customer CRUD + stats, Customer Feedback, Loyalty, Gift Cards (2 modula) |
| **Inventar & Nabava** | Materials/Inventory, Inventory Alerts, Inventory Transfers, Purchase Orders, Suppliers, Waste Tracking |
| **Finance** | Accounting (journal/revenue/expense), Expenses, Employee Tips, FURS ZAPOS (SI fiscal), CIS eRačun (HR fiscal) |
| **Osebje** | Auth & Roles (JWT/OIDC/NoAuth), Employee Performance, Scheduling, Time Clock, Staff Chat, Staff Training, Branch |
| **Prostor** | Floor Plan (visual table editor), Table management + QR codes, Reservations (2 modula), Queue/Waitlist |
| **Marketing** | Marketing Campaigns, Promotions (auto-discount rules) |
| **Operativno** | Multi-Location Dashboard, Audit Log, Reports (CSV export), Notifications, HubSync (remote sync), Languages (i18n), AI |
| **Nastavitve** | App Settings, Admin Setup Wizard |

---

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
| Testi | 58 backend paketov + 40 frontend test datotek |

---

## Hitri začetek

### Zahteve

| Komponenta | Minimalno | Priporočeno |
|------------|-----------|-------------|
| CPU | 2 jedri | 4+ jeder |
| RAM | 2 GB | 4+ GB |
| Disk | 10 GB | 50+ GB (SSD) |
| MongoDB | 5.0+ | 7.0+ |
| Go (build) | 1.25+ | 1.25+ |
| Node.js (build) | 20 LTS | 22 LTS |
| Thermal printer | ESC/POS | Epson TM-T20/T88 |
| OS | Linux / Windows / macOS | Linux (deploy) |

### Docker (priporočeno)

```bash
docker compose up
```

Ob prvem zagonu bo Setup Wizard na `http://localhost:8000` za nastavitev MongoDB povezave.

### Ročno

**Backend:**
```bash
cp config.example.yaml config.yaml
# uredi config.yaml (MongoDB nastavitve)
go run .
```

**Frontend:**
```bash
cd frontend
cp .env.example .env
npm install
npm run dev
```

API strežnik: `http://localhost:8000`
Frontend dev: `http://localhost:3000`

---

## Konfiguracija

```yaml
databases:
  - name: core
    type: mongo
    host: 127.0.0.1
    port: 27017
    database: nutrix
    username: ""
    password: ""

auth:
  jwt_secret: ""          # prazen = naključno ob zagonu
  jwt_expire_hrs: 24
  enabled: true

zitadel:
  enabled: false          # vključi za OIDC

serve_frontend: true      # streži Vue SPA iz Go backend-a
```

---

## Okoljske spremenljivke (frontend)

`.env` datoteka v `frontend/`:

| Spremenljivka | Privzeto | Opis |
|--------------|----------|------|
| `VITE_APP_BACKEND_HOST` | `localhost:8000` | Naslov Go API strežnika |
| `VITE_APP_MODULE_CORE_API_PREFIX` | `/core` | Core module API prefix |
| `VITE_APP_MODULE_FISCAL_API_PREFIX` | `/fiscal` | FURS ZAPOS API prefix |
| `VITE_APP_MODULE_FISCAL_HR_API_PREFIX` | `/fiscal_hr` | CIS eRačun API prefix |
| `VITE_APP_APP_VERSION` | `v0.0.1` | Version string v UI |

---

## Struktura projekta

```
nutrixpos/
├── main.go                  # Entrypoint
├── cmd/
│   ├── root.go              # CLI root + module registration (45 modulov)
│   └── seed.go              # Seeder za test podatke
├── common/
│   ├── config/              # Config parsing (yaml)
│   ├── customerrors/        # Standardizirane napake
│   ├── helpers/             # Helper funkcije
│   ├── logger/              # JSON logger
│   ├── middlewares/         # Rate limiting, X-Forwarded-For
│   └── userio/              # CLI prompter
├── modules/
│   ├── imodules.go          # IBaseModule, IHttpModule, IBackgroundWorkerModule
│   ├── modulebuilder.go     # Builder pattern za module
│   ├── appmanager.go        # App manager (start/stop/register)
│   ├── backgroundworkers.go # Worker management
│   ├── core/                # Core module (orders, products, inventory, customers, etc.)
│   ├── auth/                # JWT, bcrypt, OIDC
│   ├── fiscal/              # FURS ZAPOS (Slovenia)
│   ├── fiscal_hr/           # CIS eRačun (Croatia)
│   ├── hubsync/             # Remote sync
│   ├── accounting/          # Financial journal
│   ├── ai/                  # AI integration
│   ├── auditlog/            # Event history
│   ├── branch/              # Branch management
│   ├── chat/                # Staff chat
│   ├── customerdisplay/     # Customer-facing slideshow
│   ├── delivery/            # Zone-based delivery
│   ├── employee/            # Performance metrics
│   ├── expense/             # Expense tracking
│   ├── feedback/            # Customer ratings
│   ├── floorplan/           # Visual table editor
│   ├── giftcard/ + giftcards/ # Gift cards
│   ├── inventoryalerts/     # Stock alerts
│   ├── inventorytransfer/   # Cross-branch transfer
│   ├── kiosk/               # Self-service kiosk
│   ├── kitchen/             # Kitchen display
│   ├── loyalty/             # Loyalty points
│   ├── marketing/           # Campaigns
│   ├── menuengineering/     # Profitability matrix
│   ├── multilocation/       # Multi-branch dashboard
│   ├── multipayment/        # Split payments
│   ├── notification/        # In-app notifications
│   ├── onlineorder/         # Online portal
│   ├── promotion/           # Discount rules
│   ├── purchase/            # Purchase orders
│   ├── queue/               # Waitlist
│   ├── receipt/             # Receipt templates
│   ├── report/              # CSV export
│   ├── reservation/ + reservations/ # Booking
│   ├── scheduling/          # Staff shifts
│   ├── splitbill/           # Bill splitting
│   ├── supplier/            # Vendor management
│   ├── table/               # Table + QR codes
│   ├── tableside/           # QR ordering
│   ├── timeclock/           # Check-in/out
│   ├── tips/                # Employee tips
│   ├── training/            # Staff training
│   └── waste/               # Waste tracking
├── frontend/                # Vue 3 SPA
│   ├── src/
│   │   ├── components/      # PrimeVue komponente
│   │   ├── pages/           # 40+ strani
│   │   ├── services/        # API klici
│   │   ├── classes/         # Business model classes
│   │   └── __tests__/       # 40 test datotek
│   ├── .env.example         # VITE_APP_* spremenljivke
│   └── package.json
├── config.example.yaml      # Konfiguracija
├── Dockerfile               # Multi-stage build
├── docker-compose.yaml      # POS + frontend + MongoDB
├── package.nsi              # Windows NSIS installer
└── .github/workflows/ci.yml # CI/CD pipeline
```

---

## Arhitektura modulov

Vsak modul implementira enega ali več vmesnikov iz `modules/imodules.go`:

```go
type IBaseModule interface {
    OnStart() func() error
    OnEnd() func()
}

type IHttpModule interface {
    RegisterHttpHandlers(r *mux.Router, prefix string)
}

type IBackgroundWorkerModule interface {
    RegisterBackgroundWorkers() []Worker
}

type ISeederModule interface {
    Seed(entities []string, is_new_only bool) error
    GetSeedables() (entities []string, err error)
}
```

Moduli se registrirajo v `cmd/root.go`:

```go
root.RegisterModule("accounting", accounting.NewAccountingModule(logger, client))
root.RegisterModule("kitchen", kitchen.NewKitchenModule(logger, client))
// ... 45 modulov
```

Vsak modul samodejno dobi:
- **JWT avtentikacijo** (ali NoAuth za javne endpoint-e: kiosk, tableside, customer display)
- **MongoDB dostop** prek `common.GetDatabaseClient()` (singleton)
- **Cobra CLI ukaz** prek `modules/appmanager.go`
- **HTTP router** z lastnim prefixom

---

## API dokumentacija

OpenAPI 3.0 specifikacija za core module:  
📄 `modules/core/specs.api.yaml` (1772 vrstic)

Pokriva: `categories`, `products`, `materials`, `orders`, `customers`, `sales`, `disposals`, `settings`, `languages`

Zagon mock strežnika:
```bash
npx prism mock modules/core/specs.api.yaml
```

---

## Windows installer

NSIS installer v `package.nsi`:

```bash
makensis package.nsi    # generira nutrixpos-installer.exe
```

- Name: Nutrix POS
- Target: `$LocalAppData\NutrixPOS`
- License, komponente, Start Menu, bližnjice
- Uninstaller

---

## Razvoj

### Build

```bash
# Backend
go build -o pos .

# Frontend
cd frontend && npm run build
```

### Testi

```bash
# Backend (58 paketov z testi, vsi ok)
go test -count=1 -timeout 300s ./...

# Frontend (40 testnih datotek)
cd frontend && npx vitest run
```

### CI/CD

GitHub Actions pipeline (`.github/workflows/ci.yml`):
- **Backend**: vet, golangci-lint, test (race + coverage), build, govulncheck
- **Frontend**: type-check, lint (ESLint + oxlint), test (vitest), build
- **Docker**: multi-stage build + cache (odvisen od backend + frontend)

---

## Varnost

- **JWT secret** — samodejna generacija naključnega ključa ob zagonu (`crypto/rand`)
- **Rate limiting** — drseno okno na auth endpointih (10 prijav/min, 5 registracij/min)
- **NoSQL injection** — uporabniški vnos ekraniran z `regexp.QuoteMeta`
- **PIN kode** — SHA-256 salted hash
- **Gesla** — bcrypt hash
- **Error handling** — notranje napake niso izpostavljene strankam
- **Docker** — non-root uporabnik, healthcheck, stripped binary

---

## Dostopne vloge

| Vloga | Dovoljenja |
|-------|-----------|
| `superuser` | Vse |
| `admin` | Administracija, naročila, nastavitve |
| `cashier` | Blagajna, prodaja, stranke |
| `chef` | Kuhinjski zaslon, pregled materialov |

---

## Roadmap

- [x] Core POS (orders, inventory, products, customers)
- [x] FURS ZAPOS fiskalizacija (Slovenija)
- [x] CIS eRačun fiskalizacija (Hrvaška)
- [x] Kitchen display (WebSocket real-time)
- [x] Multi-payment / Split bill
- [x] Floor plan + QR tableside ordering
- [x] Reservations + Queue management
- [x] Delivery management
- [x] Self-service kiosk
- [x] Gift cards + Loyalty
- [x] Marketing campaigns + Promotions
- [x] Employee management (scheduling, time clock, tips, training)
- [x] Accounting + Expenses
- [x] Multi-location dashboard
- [x] Staff chat + Notifications
- [x] Purchase orders + Suppliers
- [x] AI integration
- [ ] Mobile apps (iOS/Android)
- [ ] Stripe/PayPal payment gateway
- [ ] POS hardware integration (barcode scanner, customer display)
- [ ] Data analytics dashboard
- [ ] API v2 (GraphQL)

---

## Licenca

GNU General Public License v2 — glej [LICENSE](LICENSE)

> **Opozorilo:** NutrixPOS je v aktivnem razvoju. Nazaj-kompatibilnost ni zagotovljena do stabilne izdaje.
