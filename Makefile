.PHONY: run-local
run-local:
	ENV=local go run cmd/gochat/main.go

.PHONY: run-dev
run-dev:
	ENV=dev go run cmd/gochat/main.go

.PHONY: run-prod
run-prod:
	ENV=prod go run cmd/gochat/main.go

.PHONY: test 
test:
	go test -v ./...

.PHONY: test-race
test-race:
	go test -race -v ./...

.PHONY: coverage
coverage:
	go test -coverprofile=c.out ./...;\
	go tool cover -func=c.out;\
	rm c.out
