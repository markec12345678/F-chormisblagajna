# Contributing

Hvala, da želiš prispevati k NutrixPOS!

## Razvojno okolje

### Zahteve

- Go 1.25+
- Node.js 22+
- MongoDB 5.0+
- Docker (opcijsko)

### Zagon

```bash
# Backend
cp config.example.yaml config.yaml
go run ./cmd/pos

# Frontend
cd frontend
npm install
npm run dev
```

## Pred spremembami

1. Forkej repozitorij
2. Ustvari vejico (`git checkout -b feature/nova-funkcija`)
3. Naredi spremembe
4. Zaženi teste (`go test -race ./...`)
5. Pošlji Pull Request

## Konvencije

### Go

- Sledi obstoječemu slogu kode
- Uporabljaj `gofmt` / `goimports`
- Vključi teste za novo funkcionalnost
- Uporabljaj `error` return namesto `panic`/`log.Fatal` v servisnih metodah

### Vue Frontend

- Functional components + hooks
- Tailwind CSS utility classes
- TypeScript za vse nove komponente

### Commits

- Jasni, opisni commit sporočila
- En commit = en logični blok sprememb

## Testiranje

```bash
# Backend
go test -race ./...

# Frontend
cd frontend
npm run test:unit
npm run type-check
npm run lint
```

## Code Review

Vsi PR-ji gredo skozi code review. Zagotovi, da:
- Build gre skozi (`go build ./...`)
- Testi delujejo (`go test ./...`)
- Ni novih warning-ov
- Koda je čista in berljiva
