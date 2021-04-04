all: clean test build
clean:
		go clean ./...
		rm -f ./go-ai-behavior
test:
		go test -v ./...
build:
		go build -v
run: build
		./go-ai-behavior $(ARGS)
pprof:
		go tool pprof -http=:8080 go-ai-behavior.prof