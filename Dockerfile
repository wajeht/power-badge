FROM golang:1.27-alpine@sha256:8a5910f31396cd4d89662f56c68b3ae31d374308270a1c3bd96672ee5ed43414 AS build
WORKDIR /app
COPY go.mod ./
COPY main.go ./
RUN CGO_ENABLED=0 go build -o power-badge .

FROM alpine@sha256:294b683cb724975bec92580e1e685676bd4b50bda910ddb8c51d4cabeaec77e6
RUN apk add --no-cache curl
COPY --from=build /app/power-badge /power-badge
HEALTHCHECK --interval=30s --timeout=5s --retries=3 CMD curl -f http://localhost/healthz
ENTRYPOINT ["/power-badge"]
