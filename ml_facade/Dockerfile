# Build environment
FROM golang:1.22-alpine AS build-env
WORKDIR /ml-facade

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -ldflags '-w -s' -a -o ./bin/api ./cmd

# Deployment environment
FROM alpine

COPY --from=build-env /ml-facade/bin/api /ml-facade/

EXPOSE 4000
CMD ["/ml-facade/api"]
