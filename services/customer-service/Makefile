.PHONY: start-dev
start-dev:
	air -c .air.toml

build:
	go build -o ./dist/main .

start:
	chmod +x ./dist/main
	./dist/main

migrate:
	sql-migrate up

migrate-undo:
	sql-migrate down