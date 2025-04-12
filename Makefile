main_package_path = ./cmd/server/main.go
migrate_package_path = ./cmd/migrate/main.go
binary_name = happypoor

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

.PHONY: no-dirty
no-dirty:
	@test -z "$(shell git status --porcelain)"


# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## lint: run all linters
.PHONY: lint
lint:
	golangci-lint run --timeout 5m --allow-parallel-runners --sort-results

## lint-fix: try to fix all lint errors
.PHONY: lint-fix
lint-fix:
	golangci-lint run --fix

## test: run all tests colorized
# .PHONY: test
test: lint
	./test.sh

## test-ci: run all tests in CI
.PHONY: test-ci
test-ci:
	go test -v -race -buildvcs ./...

## test/cover: run all tests and display coverage
.PHONY: test/cover
test/cover:
	go test -v -race -buildvcs -coverprofile=/tmp/coverage.out ./...
	go tool cover -html=/tmp/coverage.out

.PHONY: vet
vet:
	go vet ./...


# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## tidy: tidy modfiles and format .go files
.PHONY: tidy
tidy:
	go mod tidy -v
	go fmt ./...

## build: build the application
.PHONY: build
build:
	# Include additional build steps, like TypeScript, SCSS or Tailwind compilation here...
	go build -o=/tmp/bin/${binary_name} ${main_package_path}

## run: run the application
.PHONY: run
run: build
	/tmp/bin/${binary_name}

## run/live: run the application with reloading on file changes
.PHONY: run/live
run/live:
	go run github.com/cosmtrek/air@v1.43.0 \
		--build.cmd "make build" --build.bin "/tmp/bin/${binary_name}" --build.delay "100" \
		--build.exclude_dir "" \
		--build.include_ext "go, tpl, tmpl, html, css, scss, js, ts, sql, jpeg, jpg, gif, png, bmp, svg, webp, ico" \
		--misc.clean_on_exit "true"


# ==================================================================================== #
# MIGRATIONS
# ==================================================================================== #

## db/build: build the migration tool
.PHONY: db/build
db/build:
	go build -o=/tmp/bin/migrate ${migrate_package_path}

## db/migrate: run all pending migrations
.PHONY: db/migrate
db/migrate: db/build
	/tmp/bin/migrate -command up

## db/migrate/new name=$1: create a new migration file
.PHONY: db/migrate/new
db/migrate/new:
	@test $(name) || (echo "name parameter is required, e.g. make db/migrate/new name=add_users_table"; exit 1)
	@echo "Creating migration $(name)..."
	@mkdir -p ./internal/migrations/versions
	@next_version=$$(find ./internal/migrations/versions -name "*.go" | sort -V | tail -n 1 | sed -E 's/.*\/([0-9]+)_.*/\1/' | awk '{printf "%03d", $$1+1}'); \
	if [ -z "$$next_version" ]; then next_version="001"; fi; \
	filename="./internal/migrations/versions/$${next_version}_$$(echo $(name) | tr ' ' '_').go"; \
	sed 's/{{VERSION}}/'"$$next_version"'/g; s/{{NAME}}/'"$(name)"'/g; s/{{FUNC_NAME}}/'"$$(echo $(name) | tr ' ' '_' | tr '[:upper:]' '[:lower:]')"'/g' ./scripts/migration_template.go > $$filename; \
	echo "Created $$filename"

## db/status: show migration status
.PHONY: db/status
db/status: db/build
	/tmp/bin/migrate -command status

## db/rollback: rollback the last migration (if supported)
.PHONY: db/rollback
db/rollback: db/build confirm
	/tmp/bin/migrate -command down
