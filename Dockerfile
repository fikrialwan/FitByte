ARG GO_VERSION=1.24
FROM golang:${GO_VERSION}-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o app ./cmd/app

# Creating the smallest possible Docker image for production
FROM gcr.io/distroless/static-debian12:debug-nonroot

# Accept build arg for config file path
ARG CONFIG_FILE_PATH=.env

WORKDIR /app
COPY --from=builder --chown=nonroot:nonroot /app/app ./app

# Set environment variable for config file path
ENV CONFIG_FILE_PATH=${CONFIG_FILE_PATH}

ENTRYPOINT ["./app"]
