FROM golang:latest as builder

LABEL maintainer = "Matheus Carmo (a.k.a Carmel) <mematheuslc@gmail.com>"

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o authorizer .

# Final image
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /authorizer
COPY --from=builder /app/authorizer .
COPY --from=builder /app/operations .
