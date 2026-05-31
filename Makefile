main_package_path = ./cmd/server/main.go
web_package_path = ./cmd/web/main.go

migrate_package_path = ./cmd/migrate/main.go
seed_package_path = ./cmd/seed/*.go
binary_name = cashout
web_binary_name = cashout-web
linux_binary_name = ${binary_name}-linux
linux_web_binary_name = ${web_binary_name}-linux

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
# TESTS
# ==================================================================================== #

.PHONY: lint
lint:
	golangci-lint cache clean
	golangci-lint run ./...

.PHONY: lint-fix
lint-fix:
	golangci-lint run --fix ./...

.PHONY: test-ci
test-ci: lint
	go test -v -race -buildvcs ./...

.PHONY: test-coverage
test-coverage: lint
	go test -v -race -buildvcs -coverprofile=/tmp/coverage.out ./...
	go tool cover -html=/tmp/coverage.out

.PHONY: test
test: lint
	gotestsum --format dots -- -race -buildvcs ./...

.PHONY: test/live
test/live:
	gotcha watch # --fast

.PHONY: vet
vet:
	go vet ./...

.PHONY: sec
sec:
	go list -json -deps ./... | nancy sleuth
	govulncheck ./...


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
	go build -o=/tmp/bin/${binary_name} ${main_package_path}

## build-web: build the web server
.PHONY: build-web
build-web:
	go build -o=/tmp/bin/${web_binary_name} ${web_package_path}

## build-all: build both the bot and web server
.PHONY: build-all
build-all: build build-web

## build-linux: build the application for linux x86_64 (CentOS)
.PHONY: build-linux
build-linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o=/tmp/bin/${linux_binary_name} ${main_package_path}

## build-linux-web: build the web server for linux x86_64 (CentOS)
.PHONY: build-linux-web
build-linux-web:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o=/tmp/bin/${linux_web_binary_name} ${web_package_path}

## build-linux-all: build both applications for linux
.PHONY: build-linux-all
build-linux-all: build-linux build-linux-web

## run: run the application
.PHONY: run
run: build
	/tmp/bin/${binary_name}

## run-web: run the web server
.PHONY: run-web
run-web: build-web
	/tmp/bin/${web_binary_name}

## run-all: run both the bot and web server (requires two terminals)
.PHONY: run-all
run-all:
	@echo "Starting both services..."
	@echo "Run 'make run' in one terminal and 'make run-web' in another"
	@echo "Or use 'make run/live-all' for live reloading of both"

## run/live: run the application with reloading on file changes
.PHONY: run/live
run/live:
	go run github.com/cosmtrek/air@v1.43.0 \
		--build.cmd "make build" --build.bin "/tmp/bin/${binary_name}" --build.delay "100" \
		--build.exclude_dir "" \
		--build.include_ext "go, tpl, tmpl, html, css, scss, js, ts, sql, jpeg, jpg, gif, png, bmp, svg, webp, ico" \
		--misc.clean_on_exit "true"

## run/live-web: run the web server with reloading on file changes
.PHONY: run/live-web
run/live-web:
	go run github.com/cosmtrek/air@v1.43.0 \
		--build.cmd "make build-web" --build.bin "/tmp/bin/${web_binary_name}" --build.delay "100" \
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

## db/status: show migration status
.PHONY: db/status
db/status: db/build
	/tmp/bin/migrate -command status

## db/rollback: rollback the last migration (if supported)
.PHONY: db/rollback
db/rollback: db/build confirm
	/tmp/bin/migrate -command down


# ==================================================================================== #
# DATABASE SEEDING
# ==================================================================================== #

## db/seed/build: build the seed tool
.PHONY: db/seed/build
db/seed/build:
	go build -o=/tmp/bin/seed ${seed_package_path}

## db/seed: seed the database with random transactions for a user (requires SEED_USER_TG_ID env var)
.PHONY: db/seed
db/seed: db/seed/build
	/tmp/bin/seed


# ==================================================================================== #
# OPENAPI + SDKS
# ==================================================================================== #

SWAG_VERSION       := v1.16.6
OPENAPI_WRAPPER    := @openapitools/openapi-generator-cli@2.20.5
OPENAPI_GEN_JAR    := 7.10.0
# Use the npm wrapper around openapi-generator. It downloads the matching JAR
# and runs it via local Java. Requires `npx` and Java 11+ on PATH.
OPENAPI_GEN := OPENAPI_GENERATOR_VERSION=$(OPENAPI_GEN_JAR) npx --yes $(OPENAPI_WRAPPER)

# Module identity used by the Go SDK so it resolves under `go get`.
# Must match the actual GitHub owner/repo.
SDK_GIT_USER := alainrk
SDK_GIT_REPO := cashout

## openapi: generate api/swagger.{yaml,json} from swag annotations
.PHONY: openapi
openapi:
	go run github.com/swaggo/swag/cmd/swag@$(SWAG_VERSION) init \
	  --generalInfo cmd/server/main.go \
	  --output api \
	  --parseDependency --parseInternal \
	  --outputTypes yaml,json

## sdk-python: generate the Python SDK into sdks/python
.PHONY: sdk-python
sdk-python:
	$(OPENAPI_GEN) generate \
	  -i api/swagger.yaml -g python \
	  -o sdks/python \
	  --git-user-id=$(SDK_GIT_USER) --git-repo-id=$(SDK_GIT_REPO) \
	  --additional-properties=packageName=cashout_sdk,projectName=cashout-sdk
	# The generator emits `license = "NoLicense"`, a bare string that setuptools
	# validates as an SPDX expression and rejects (it is not valid SPDX), breaking
	# `pip`/`uv` installs. Rewrite it to the table form, which installs cleanly.
	@sed -i.bak 's|^license = "NoLicense"$$|license = { text = "NoLicense" }|' sdks/python/pyproject.toml
	@rm -f sdks/python/pyproject.toml.bak

## sdk-go: generate the Go SDK into sdks/go
.PHONY: sdk-go
sdk-go:
	$(OPENAPI_GEN) generate \
	  -i api/swagger.yaml -g go \
	  -o sdks/go \
	  --git-user-id=$(SDK_GIT_USER) --git-repo-id=$(SDK_GIT_REPO) \
	  --additional-properties=packageName=cashout,isGoSubmodule=true
	# The Go generator hard-codes the module path to {gitUserId}/{gitRepoId}/{packageName}
	# which does not match our on-disk layout (sdks/go). Rewrite the module path so
	# `go get github.com/$(SDK_GIT_USER)/$(SDK_GIT_REPO)/sdks/go` resolves correctly.
	@sed -i.bak 's|^module .*|module github.com/$(SDK_GIT_USER)/$(SDK_GIT_REPO)/sdks/go|' sdks/go/go.mod && rm sdks/go/go.mod.bak

## sdk-ts: generate the TypeScript (fetch) SDK into sdks/typescript
.PHONY: sdk-ts
sdk-ts:
	$(OPENAPI_GEN) generate \
	  -i api/swagger.yaml -g typescript-fetch \
	  -o sdks/typescript \
	  --git-user-id=$(SDK_GIT_USER) --git-repo-id=$(SDK_GIT_REPO) \
	  --additional-properties=npmName=@cashout/sdk,npmVersion=1.0.0,supportsES6=true

## sdks: regenerate the OpenAPI spec and all SDKs
.PHONY: sdks
sdks: openapi sdk-python sdk-go sdk-ts
