# Stage 1: Build
FROM golang:1.25-alpine AS build
WORKDIR /go/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN apk add --no-cache build-base
RUN CGO_ENABLED=1 GOOS=linux go build -trimpath -ldflags="-s -w" -o ./pos

# Stage 2: Final
FROM alpine:3.20
RUN apk --no-cache add ca-certificates tzdata
RUN addgroup -S nutrix && adduser -S nutrix -G nutrix

WORKDIR /app
COPY --from=build /go/src/app/pos .
COPY --from=build /go/src/app/assets ./assets/
COPY --from=build /go/src/app/config.example.yaml ./config.yaml

RUN chown -R nutrix:nutrix /app
USER nutrix

EXPOSE 8000
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget -qO- http://localhost:8000/ || exit 1

CMD ["./pos"]
