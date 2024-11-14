include .env
MIGRATION_PATH = ./cmd/migrate/migrations
autocannon: R ?= 4000
autocannon: D ?= 2
autocannon: C ?= 10

.PHONY: run
run:
	@go run ./cmd/api

.PHONY: test
test:
	@go test -v ./...

.PHONY: migrate-create
migration:
	@migrate create -seq -ext sql -dir $(MIGRATION_PATH) $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate-up
migrate-up:
	@migrate -path=$(MIGRATION_PATH) -database=$(DB_ADDR) up

.PHONY: migrate-down
migrate-down:
	@migrate -path=$(MIGRATION_PATH) -database=$(DB_ADDR) down

.PHONY: sonar-scan
sonar-scan:
	@${PATH_SONAR_SCAN} -X \
		-Dsonar.organization=${SONAR_ORGANIZATION} \
		-Dsonar.projectKey=${SONAR_PROJECT_KEY} \
		-Dsonar.sources=${SONAR_SOURCES} \
		-Dsonar.login=${SONAR_TOKEN} \
		-Dsonar.host.url=${SONAR_HOST}

.PHONY: seed
seed:
	@go run cmd/migrate/seed/main.go

.PHONY: gen-docs
gen-docs:
	@swag init -g ./api/main.go -d cmd,internal && swag fmt

.PHONY: test-autocannon
autocannon:
	@npx autocannon -r $(R) -d $(D) -c $(C) --renderStatusCodes $(filter-out $@,$(MAKECMDGOALS))

.PHONY: gen-proto-camunda
gen-proto-camunda:
	@protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative internal/proto/gateway.proto
