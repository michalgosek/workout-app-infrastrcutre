# syntax=docker/dockerfile:1
FROM golang:1.19-buster as build

WORKDIR /app

ADD . ./

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /trainings-service cmd/srv/main.go

FROM scratch

COPY .env  /
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /trainings-service /
