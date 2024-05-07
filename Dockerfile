FROM golang as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY cmd/app ./cmd/app
COPY internal ./internal

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/app

FROM alpine:latest

RUN apk --no-cache add ca-certificates bash

WORKDIR /root/

COPY --from=builder /app/main .

CMD ["./main"]