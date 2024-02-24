BINARY_NAME=progress-tracker
CLIENT_PATH=./cmd/client/main.go
SERVER_PATH=./cmd/server/main.go

run-client-dev:
	go run ${CLIENT_PATH}

run-server-dev:
	go run ${SERVER_PATH}

build-server:
	go build -o ${BINARY_NAME}-server ${SERVER_PATH}

build-client:
	go build -o ${BINARY_NAME}-client ${CLIENT_PATH}

clean:
	go clean
	rm ${BINARY_NAME}-server
	rm ${BINARY_NAME}-client

test:
	go test ./test/... -v

test-coverage:
	go test ./test/... -coverprofile=coverage.out