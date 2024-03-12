.PHONY: test test-race coverage

test:
	go test -v ./...

test-race:
	go test -race -v ./...

coverage:
	go test -coverprofile=c.out ./...;\
	go tool cover -func=c.out;\
	rm c.out
