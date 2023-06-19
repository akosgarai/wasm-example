FROM golang:1.20-alpine AS builder

# copy the source code into the container
COPY ./cmd /app/cmd
COPY ./pkg /app/pkg
COPY ./bin /app/bin
COPY ./assets /app/assets
COPY ./go.mod /app/go.mod
COPY ./go.sum /app/go.sum

# set working directory
WORKDIR /app

# get the dependencies
RUN go mod download

# build the backend binary to the /app/bin directory
RUN go build -o bin/router cmd/server/main.go

# build the frontend wasm application to the /app/assets directory
RUN GOOS=js GOARCH=wasm go build -o assets/app.wasm cmd/wasm/main.go

# deploy the application to a scratch container
FROM scratch

WORKDIR /app

# copy the binary from the builder container
COPY --from=builder /app/bin/router /app/bin/router
# copy the assets from the builder container
COPY --from=builder /app/assets /app/assets

# expose port 9090
EXPOSE 9090

# run the binary
ENTRYPOINT ["/app/bin/router"]
