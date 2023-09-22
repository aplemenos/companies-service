# ==============================================================================
# Tools commands

run-linter:
	echo "Starting linters"
	golangci-lint run ./... -c .golangci.yml

swaggo:
	echo "Starting swagger generating"
	swag init -g **/**/*.go


# ==============================================================================
# Main

run:
	go run ./cmd/companies-api/main.go

build:
	go build ./cmd/companies-api/main.go

test:
	go test -cover ./...

# ==============================================================================
# Modules support

tidy:
	go mod tidy

deps-reset:
	git checkout -- go.mod
	go mod tidy

deps-upgrade:
	go get -u -t -d -v ./...
	go mod tidy

deps-cleancache:
	go clean -modcache

# ==============================================================================
# Go migrate postgresql

force:
	migrate -database postgres://postgres:postgres@localhost:5432/company_db?sslmode=disable -path migrations force 1

version:
	migrate -database postgres://postgres:postgres@localhost:5432/company_db?sslmode=disable -path migrations version

migrate_create:
	migrate create -ext sql -dir ./migrations -seq create_initial_tables

migrate_up:
	migrate -database postgres://postgres:postgres@localhost:5432/company_db?sslmode=disable -path migrations up

migrate_down:
	migrate -database postgres://postgres:postgres@localhost:5432/company_db?sslmode=disable -path migrations down

# ==============================================================================
# Docker compose commands

# develop:
# 	echo "Starting docker environment"
# 	docker-compose -f docker-compose.dev.yml up --build

prod:
	echo "Starting docker prod environment"
	docker-compose -f docker-compose.yml up --build

# ==============================================================================
# Docker support

FILES := $(shell docker ps -aq)

down-local:
	docker stop $(FILES)
	docker rm $(FILES)

clean:
	docker system prune -f

logs-local:
	docker logs -f $(FILES)