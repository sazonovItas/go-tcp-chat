FROM golang:alpine AS build-stage

WORKDIR /build

ADD go.mod .
ADD go.sum .
RUN go mod download

COPY . .

RUN GOOS=linux CGO_ENABLED=0 go build -o /app/chat ./cmd/gochat

FROM build-stage AS run-test-stage
RUN go test -v ./...

FROM alpine AS build-release-stage

WORKDIR /app

COPY --from=build-stage /app/chat /app/chat
COPY --from=build-stage /build/config /app/config

EXPOSE 5050
