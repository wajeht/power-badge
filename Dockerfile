FROM golang:1.26-alpine AS build
WORKDIR /app
COPY go.mod ./
COPY main.go ./
RUN CGO_ENABLED=0 go build -o power-badge .

FROM alpine
COPY --from=build /app/power-badge /power-badge
HEALTHCHECK --interval=30s --timeout=5s --retries=3 CMD wget -qO /dev/null http://localhost/healthz
ENTRYPOINT ["/power-badge"]
