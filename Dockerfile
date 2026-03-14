FROM golang:1.26-alpine@sha256:2389ebfa5b7f43eeafbd6be0c3700cc46690ef842ad962f6c5bd6be49ed82039 AS build
WORKDIR /app
COPY go.mod ./
COPY main.go ./
RUN CGO_ENABLED=0 go build -o power-badge .

FROM alpine
COPY --from=build /app/power-badge /power-badge
HEALTHCHECK --interval=30s --timeout=5s --retries=3 CMD wget -qO /dev/null http://localhost/healthz
ENTRYPOINT ["/power-badge"]
