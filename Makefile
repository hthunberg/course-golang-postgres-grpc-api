
# The go command now disables cgo by default on systems without a C toolchain. 
# Enable statically linked binaries to make the application more portable. 
# It allows us to use the binary with source images that don't support shared 
# libraries when building our container images.
CGO_FLAGS 	:= CGO_ENABLED=1

# The name of the compilation architecture. 
GO_ARCH 	:= GOARCH=amd64

.PHONY: all
all: sqlc fmt lint test integrationtest gosec build

.PHONY: all-tools
all-tools: lint-install sqlc-install gosec-install goimports-install mockgen-install

.PHONY: fmt
fmt:
	go fmt ./... ./test/integration

.PHONY: build
build: fmt test gosec
	$(GO_ARCH) $(CGO_FLAGS) \
	go build -buildmode=pie -ldflags "-s -w" -o bin/bank cmd/bank/main.go

.PHONY: run
run: build
	$(GO_ARCH) $(CGO_FLAGS) \
	./bin/bank

.PHONY: test
test: fmt mock
	$(GO_ARCH) $(CGO_FLAGS) \
	go test -race -v -cover -short ./...

# Server requires a running PostgreSQL `make -C build up`
.PHONY: server
server:
	go run cmd/bank/main.go

.PHONY: integrationtest
integrationtest: fmt
	$(GO_ARCH) $(CGO_FLAGS) \
	go test -race -v -cover ./... -tags=integration
	
.PHONY: lint-install
lint-install:
	go install github.com/mgechev/revive@latest

.PHONY: lint
lint: fmt
	revive -config=revivie-lint.toml ./... 

.PHONY: mockgen-install
mockgen-install:
	go install go.uber.org/mock/mockgen@latest

# Generate mocks
.PHONY: mock
mock: mockgen-install
	mockgen -source internal/app/bank/bank.go -destination internal/app/bank/mock/bank.go  -package mockdb

# We use https://docs.sqlc.dev/en/stable/index.html for database queries and mapping. This library
# has support for PostgreSQL, MySQL and SQLite, no other DBs supported.
.PHONY: sqlc-install
sqlc-install:
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Generate DB code using sqlc based on sqlc.yaml
.PHONY: sqlc
sqlc:
	sqlc generate -f build/db/sqlc.yaml

# Automates security checks
.PHONY: gosec-install
gosec-install:
	go install github.com/securego/gosec/v2/cmd/gosec@latest

.PHONY: gosec
gosec: goimports
	gosec ./... 

# Fix issue with unused import in sqlc generated file
.PHONY: goimports-install
goimports-install:
	go install golang.org/x/tools/cmd/goimports@latest

.PHONY: goimports
goimports:
	goimports -w internal/app/db/querier.go 

.PHONY: docker-build
docker-build:
	DOCKER_DEFAULT_PLATFORM=linux/amd64 \
	DOCKER_BUILDKIT=1 \
	docker build \
		-f build/docker/Dockerfile \
		-t "bank" \
		.

