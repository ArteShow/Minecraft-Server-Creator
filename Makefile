.PHONY: tests deps tidy lint lint-fix build up down

#test all packages
tests:
	go test -v ./...

#download dependencies
deps:
	go mod download

#tidy
tidy:
	go mod tidy

# Run linter
lint:
	@find ./services -type f -name "go.mod" | while read gomod_file; do \
		module_dir=$$(dirname "$$gomod_file"); \
		echo "Linting module in: $$module_dir"; \
		(cd "$$module_dir" && golangci-lint run); \
		if [ $$? -ne 0 ]; then \
			echo "Linting failed for $$module_dir"; \
			exit 1; \
		fi; \
	done; \

# Run linter with fix
lint-fix:
	@find ./services -type f -name "go.mod" | while read gomod_file; do \
		module_dir=$$(dirname "$$gomod_file"); \
		echo "Linting module in: $$module_dir"; \
		(cd "$$module_dir" && golangci-lint run --fix); \
		if [ $$? -ne 0 ]; then \
			echo "Linting failed for $$module_dir"; \
			exit 1; \
		fi; \
	done; \

#build docker
build:
	docker-compose --env-file config/docker.env up --build -d


#run docker
up:
	docker-compose --env-file config/docker.env up -d


#shut down docker
down:
	docker-compose --env-file config/docker.env down
