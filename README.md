# Example wasm application

Written in golang.
Routing with gin.
Frontend wasm application.

## Docker system

- Dockerfile for the site container

Build:

```bash
docker build -t wasm-examle:latest -f Dockerfile .
```

Run:

```bash
docker run -p -e ASSETS_DIR=/app/assets 9090:9090 wasm-examle:latest
```

### Docker compose

- app contains the site

```bash
docker compose up -d
```
