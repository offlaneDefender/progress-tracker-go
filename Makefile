BINARY_NAME=progress-tracker
CLIENT_PATH=./cmd/client/main.go
SERVER_PATH=./cmd/server/main.go

run-client-dev:
	go run ${CLIENT_PATH}

run-server-dev:
	go run ${SERVER_PATH}

build-server:
	go build -o bin/${BINARY_NAME}-server ${SERVER_PATH}

build-client:
	go build -o bin/${BINARY_NAME}-client ${CLIENT_PATH}

clean:
	go clean
	rm -r bin/

run-test:
	go test ./test/... -v

run-test-coverage:
	go test ./test/... -coverprofile=coverage.out
