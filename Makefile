.PHONY: tests tidy build build down

tests:
	@find ./services -type f -name "go.mod" | while read gomod_file; do \
		module_dir=$$(dirname "$$gomod_file"); \
		echo "Testing module in: $$module_dir"; \
		(cd "$$module_dir" && go test -v ./...) || exit 1; \
	done

tidy:
	@find ./services -type f -name "go.mod" | while read gomod_file; do \
		module_dir=$$(dirname "$$gomod_file"); \
		echo "Tidying module in: $$module_dir"; \
		(cd "$$module_dir" && go mod tidy) || exit 1; \
	done

build:
	docker compose -f host/docker-compose.yml up -d

down:
	docker compose -f host/docker-compose.yml down
