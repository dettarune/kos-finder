APP_NAME=kos-finder
BINARY=bin/$(APP_NAME)

DB_USERNAME=detarune
DB_NAME=kos_finder

build:
	go build -o $(BINARY) ./cmd/api/main.go

run:
	go run ./cmd/api/main.go

start: build
	./$(BINARY)

clean:
	rm -rf $(BINARY)

test:
	go test ./... -v

migrate:
	cat ./db/migrations/*.sql | mysql -u$(DB_USERNAME) -p$(DB_PASSWORD) $(DB_NAME)
