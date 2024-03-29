FROM golang:1.18.0-alpine3.15 AS builder

RUN apk update
RUN apk add git

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /app/mqtt ./cmd


FROM alpine

WORKDIR /
COPY --from=builder /app/mqtt .

# tcp
EXPOSE 1883

# websockets
EXPOSE 1882

# dashboard
EXPOSE 8080

ENTRYPOINT [ "/mqtt" ]
