FROM golang:1.26-alpine AS build
WORKDIR /app
COPY go.mod ./
COPY main.go ./
RUN CGO_ENABLED=0 go build -o power-badge .

FROM alpine
RUN apk add --no-cache curl
COPY --from=build /app/power-badge /power-badge
HEALTHCHECK --interval=30s --timeout=5s --retries=3 CMD curl -f http://localhost/healthz
ENTRYPOINT ["/power-badge"]
