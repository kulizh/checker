APP_NAME=checker
BUILD_DIR=dist

.PHONY: build run clean tidy deploy

build:
	go mod tidy
	go build -o $(APP_NAME) ./cmd/checker
	
build-linux:
	GOOS=linux GOARCH=amd64 go build -o checker ./cmd/checker

run:
	go run ./cmd/checker -config configs/domains.example.json -interval 30s

clean:
	rm -rf $(APP_NAME) $(BUILD_DIR)

tidy:
	go mod tidy

dist:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/checker
	cp -r configs $(BUILD_DIR)/ || true