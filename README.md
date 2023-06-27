# Example wasm application

Written in golang.
Routing with gin.
Frontend wasm application.
Websocket with gorilla.

The frontend contains a form, where a couple of parameter could be set. On case of form submission, the backend validates the data and executes a script on a remote machine. The validated data is used as the arguments of the executed script.

## Docker system

- Dockerfile for the site container

Build:

```bash
docker build -t wasm-example:latest -f Dockerfile .
```

Run:

```bash
docker run -p -e ASSETS_DIR=/app/assets 9090:9090 wasm-examle:latest
```

### Docker compose

- app contains the backend and frontend application
- staging represents the staging remote machine
- production represents the production remote machine.

```bash
docker compose up -d
```

When the go code changes the app image has to be recreated. It could be achived with adding the `--build` flag to the compose up command.

```bash
docker compose up -d --build
```
