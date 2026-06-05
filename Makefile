APP_NAME=checker
BUILD_DIR=dist

.PHONY: build build-linux run clean tidy

build:
	go mod tidy
	go build -o $(APP_NAME) ./cmd/checker
	go build -o domains-helper ./cmd/domains-helper

build-helper:
	go mod tidy
	go build -o domains-helper ./cmd/domains-helper

build-linux:
	GOOS=linux GOARCH=amd64 go build -o $(APP_NAME) ./cmd/checker

run:
	go run ./cmd/checker -config configs/domains.example.json -interval 30s

clean:
	rm -f $(APP_NAME)
	rm -rf $(BUILD_DIR)

tidy:
	go mod tidy