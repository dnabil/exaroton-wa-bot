# Builder
FROM golang:1.26-alpine AS builder

RUN apk add --no-cache git gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -ldflags "-s -w" -o bin/app cmd/app/main.go

# Runner
FROM alpine:3.19 AS runner

WORKDIR /app/bin

COPY --from=builder /app/bin/app .

COPY --from=builder /app/public ../public
COPY --from=builder /app/pages ../pages

EXPOSE 8080

CMD ["./app", "--cfg=../config.yml"]
