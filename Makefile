build:
	@go build -o server main.go

run: build
	@./server

all: build

clean:
	@rm -rf server

