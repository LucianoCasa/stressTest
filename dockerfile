# Etapa de build
FROM golang:1.24 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /dist/stress ./cmd/app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /dist/server ./cmd/server

# Etapa final
FROM gcr.io/distroless/base-debian12 AS stress
WORKDIR /app
COPY --from=build /dist/stress /app/stress
ENTRYPOINT ["/app/stress"]

FROM gcr.io/distroless/base-debian12 AS server
WORKDIR /app
COPY --from=build /dist/server /app/server
ENTRYPOINT ["/app/server"]

