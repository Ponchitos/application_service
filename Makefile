.PHONY: all compile install build clean run test get migration_create up_db down_db

MIGRATION_PATH=migrations/postgres
GOCMD=go
BINARYNAME=appservice

all: compile

compile: clean get build

install:
	@$(GOCMD) install -v ./...

build:
	@echo "  >  Building binary..."
	@$(GOCMD) build -o $(BINARYNAME) -v ./server/cmd/*.go

clean:
	@echo "  >  Cleaning build cache"
	@$(GOCMD) clean ./...
	@rm -f $(BINARYNAME)

run:
	@./$(BINARYNAME)

test:
	@echo "  >  Running tests"
	@$(GOCMD) test -v -race -timeout 30s ./...

get:
	@echo "  >  Checking dependencies..."
	@$(GOCMD) mod tidy

migration_create:
	migrate create -ext sql -dir $(MIGRATION_PATH) ${name}

up_db:
	docker-compose -f dev.server.docker-compose.yml up -d --build

down_db:
	docker-compose -f dev.server.docker-compose.yml down