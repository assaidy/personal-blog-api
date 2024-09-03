build:
	@echo "> start building project..."
	@go build -o ./bin/server cmd/main.go
	@echo "> build finished successfully"

run: build
	@echo "> running the project..."
	@./bin/server

all: build

clean:
	@echo "> start cleanup..."
	@rm -rf bin
	@rm -rf db/data.db
	@echo "> cleanup finished successfully"

