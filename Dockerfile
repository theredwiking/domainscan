FROM golang:1.22 AS builder

WORKDIR /usr/src/domainscan

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o dist/app

FROM ubuntu:24.04 AS server

WORKDIR /usr/src/app

RUN apt update && apt install nmap -y

COPY --from=builder /usr/src/domainscan/dist .

#CMD [ "app" ]
