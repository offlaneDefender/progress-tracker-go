run-client:
	go run ./cmd/client/main.go

run-server:
	go run ./cmd/server/main.go

run-test:
	go test ./test/... -v