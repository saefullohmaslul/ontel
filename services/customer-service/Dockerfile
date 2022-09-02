FROM golang:1.18-alpine AS dev

RUN apk add --no-cache make git gcc libc-dev
RUN go install github.com/cosmtrek/air@latest
RUN go install github.com/rubenv/sql-migrate/...@latest
RUN mkdir -p /app

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

EXPOSE 7150

CMD ["make", "start-dev"]