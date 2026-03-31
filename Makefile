.PHONY: tests tidy build down

tests:
	@find ./host/services ./hub/services -type f -name "go.mod" | while read gomod_file; do \
		module_dir=$$(dirname "$$gomod_file"); \
		echo "Testing module in: $$module_dir"; \
		(cd "$$module_dir" && go test -v ./...) || exit 1; \
	done

tidy:
	@find ./host/services ./hub/services -type f -name "go.mod" | while read gomod_file; do \
		module_dir=$$(dirname "$$gomod_file"); \
		echo "Tidying module in: $$module_dir"; \
		(cd "$$module_dir" && go mod tidy) || exit 1; \
	done

build:
	docker compose --env-file host/config/docker.env -f host/docker-compose.yaml up -d --build
	docker compose --env-file hub/config/docker.env -f hub/docker-compose.yaml up -d --build
	docker compose --env-file website/config/docker.env -f website/docker-compose.yaml up -d --build

down:
	docker compose --env-file host/config/docker.env -f host/docker-compose.yaml down
	docker compose --env-file hub/config/docker.env -f hub/docker-compose.yaml down
	docker compose --env-file website/config/docker.env -f website/docker-compose.yaml up -d --build
