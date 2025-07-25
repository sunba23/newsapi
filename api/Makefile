MAIN_PACKAGE_PATH := ./cmd/news
BINARY_NAME := news

build:
	GOOS=linux go build -o bin/${BINARY_NAME} ${MAIN_PACKAGE_PATH}

localrun:
	go run ${MAIN_PACKAGE_PATH}

devserver:
	air -c .air.toml

test:
	go test ./...

test-cov:
	go test -race ./... -coverprofile=coverage.txt -covermode=atomic

view-report:
	go tool cover -html=coverage.txt
