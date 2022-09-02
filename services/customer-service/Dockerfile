FROM golang:1.18-stretch AS dev

LABEL MAINTAINER="Saefulloh Maslul"
LABEL APP="Customer Service"

RUN apt-get update && \
    apt-get install -y git make build-essential mercurial &&\
    apt-get install -y make
RUN go install github.com/cosmtrek/air@latest
RUN go install github.com/rubenv/sql-migrate/...@latest
RUN mkdir -p /app

RUN echo "Asia/Jakarta" > /etc/timezone && rm /etc/localtime && dpkg-reconfigure --frontend noninteractive tzdata

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

EXPOSE 7150

CMD ["make", "start-dev"]