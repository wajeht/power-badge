# Power Badge

[![CI](https://github.com/wajeht/power-badge/actions/workflows/build.yml/badge.svg?branch=main)](https://github.com/wajeht/power-badge/actions/workflows/build.yml)
[![Open Source Love](https://badges.frapsoft.com/os/v1/open-source.svg?v=103)](https://github.com/wajeht/power-badge)

Display a Home Assistant power sensor as a dynamic [Shields.io endpoint badge](https://shields.io/badges/endpoint-badge).

Power Badge reads a sensor state from the Home Assistant API and returns it in the Shields endpoint schema.

## Endpoints

| Endpoint   | Description                     |
| ---------- | ------------------------------- |
| `/`        | current power reading as JSON   |
| `/healthz` | health check                    |

## Configuration

| Environment variable | Description                              | Default |
| -------------------- | ---------------------------------------- | ------- |
| `HA_URL`             | Home Assistant base URL                  | required |
| `HA_TOKEN`           | Home Assistant long-lived access token   | required |
| `HA_SENSOR_ID`       | Power sensor entity ID                   | required |
| `PORT`               | HTTP listen port                         | `80` |

## Run with Docker

Set the required environment variables in your shell, then build and run the container:

```sh
docker build -t power-badge .
docker run --rm \
  --publish 8080:80 \
  --env HA_URL \
  --env HA_TOKEN \
  --env HA_SENSOR_ID \
  power-badge
```

Use the service URL with Shields.io:

```text
https://img.shields.io/endpoint?url=https%3A%2F%2Fexample.com%2F
```

## Development

```sh
go test ./...
```
