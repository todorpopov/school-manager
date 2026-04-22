FROM golang:1.25-alpine AS backend-builder

WORKDIR /backend

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o app ./cmd/api/main.go

FROM alpine:3.18 AS runner

WORKDIR /app

COPY --from=backend-builder /backend/app .

RUN mkdir -p /app/uploads

EXPOSE 8080

CMD ["./app"]