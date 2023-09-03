# syntax=docker/dockerfile:1

FROM golang:1.20

WORKDIR /ally

# COPY go.mod ./

# RUN go mod download

# # Build
# RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-gs-ping

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/engine/reference/builder/#expose
EXPOSE 3000

# Run
RUN go install github.com/githubnemo/CompileDaemon@latest
ENTRYPOINT CompileDaemon --build="go build docker/containers/ally/main.go" --command="./main" --directory="./"
# CMD ["go", "run", "example/docker/containers/ally/main.go"]
