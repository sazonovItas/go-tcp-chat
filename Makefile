.PHONY: run-local
run-local:
	CONFIG_PATH=./configs/local.yaml go run cmd/gochat/main.go

.PHONY: run-dev
run-dev:
	CONFIG_PATH=./configs/dev.yaml go run cmd/gochat/main.go

.PHONY: run-prod
run-prod:
	CONFIG_PATH=./configs/prod.yaml go run cmd/gochat/main.go

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
