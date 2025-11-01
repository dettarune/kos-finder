APP_NAME=go-Commerce
BINARY=bin/$(APP_NAME)

DB_USERNAME=detarune
DB_NAME=go_commerce

# Build binary
build:
	go build -o $(BINARY) ./cmd/api/main.go

# Run server
run:
	go run ./cmd/api/main.go

# Jalankan binary yang udah dibuild
start: build
	./$(BINARY)

# Clean binary
clean:
	rm -rf $(BINARY)

# Jalankan unit test
test:
	go test ./... -v

# Contoh migrate database (kalau nanti dipake)
migrate:
	cat ./db/migrations/*.sql | mysql -u$(DB_USERNAME) -p$(DB_PASSWORD) $(DB_NAME)
