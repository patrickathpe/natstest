# (C) Copyright 2023-2024 Hewlett Packard Enterprise Development LP

# Project directories

PROJ_DIR                 := .
CMD_DIR                  := ./cmd
INTERNAL_DIR             := ./internal

# Go variables

LINT_DIRS                := $(PROJ_DIR)/...
FMT_DIRS                 := $(CMD_DIR)/ $(INTERNAL_DIR)/
COVERAGE_OUT             := coverage.out
UNIT_TEST_REPORT         := unit_test_coverage.html
LDFLAGS                  := -X main.Version=$(VERSION)
GCFLAGS                  := all=-N -l

# Development NATS variables
NATS_DIR                 := ${PROJ_DIR}/nats

## TARGET: RUN IN: DESCRIPTION
## help: any: Output this message and exit.
.PHONY: help
help:
	@fgrep -h '##' $(MAKEFILE_LIST) | fgrep -v fgrep | column -t -s ':' | sed -e 's/## //'

## clean: any: remove build and test output artifacts
.PHONY: clean
clean:
	rm -rf \
		$(UNIT_TEST_REPORT) \
		$(COVERAGE_OUT)

## lint: cds-dev: run all project linters
.PHONY: lint
lint: golint

## leaks: host: Run gitleaks to detect unintended hard coded secrets
.PHONY: leaks
leaks:
	@gitleaks detect -v .

## test: cds-dev: run all project unit tests
.PHONY: test
test: gounit

#----------------------------------------------------------------------
# Go targets
#----------------------------------------------------------------------

## vendor: cds-dev: download vendored dependencies
vendor: go.mod go.sum
	go mod vendor

## tidy: cds-dev: tidy dependencies
.PHONY: tidy
tidy:
	go mod tidy

## fmt: cds-dev: format Go code
.PHONY: fmt
fmt:
	gofumpt -w $(FMT_DIRS)
	gci write $(FMT_DIRS)

## golint: cds-dev: run the linter on Go code
.PHONY: golint
golint: vendor
	golangci-lint -j 2 run $(LINT_DIRS) --timeout=3m

## test: cds-dev: run Go unit tests
.PHONY: gounit
gounit: vendor
	go test $(UNIT_TEST_DIRS) -cover -coverprofile=$(COVERAGE_OUT) -count=1 -p=1
	go tool cover -func=$(COVERAGE_OUT)
	go tool cover -html=$(COVERAGE_OUT) -o $(UNIT_TEST_REPORT)
#----------------------------------------------------------------------
# NATS management targets
#----------------------------------------------------------------------

## nats-up: host: Starts and configures NATS containers
.PHONY: nats-up
nats-up: ${NATS_DIR}/keys/operator.jwt
	docker compose --profile nats up -d nats-box nats-server
	sleep 1
	${NATS_DIR}/configure-nats-host.sh

## nats-shell: host: Start an Alpine shell session in nats-box
.PHONY: nats-shell
nats-shell:
	docker compose --profile nats exec -it -u nats -w /home/nats nats-box /bin/ash

## nats-down: host: Stop NATS containers
.PHONY: nats-down
nats-down:
	docker compose --profile nats rm --stop --force --volumes nats-box nats-server

## nats-clean: host: Stop NATS containers and delete generated keys
.PHONY: nats-clean
nats-clean: nats-down
	rm -rf ${NATS_DIR}/nsc ${NATS_DIR}/keys

${NATS_DIR}/keys/operator.jwt: ${NATS_DIR}/configure-nsc-host.sh ${NATS_DIR}/configure-nsc-box.sh
	mkdir -p "${NATS_DIR}/nsc" "${NATS_DIR}/keys"
	docker compose --profile nats up -d nats-box
	${NATS_DIR}/configure-nsc-host.sh

#----------------------------------------------------------------------
# Development environment targets
# Start, stop & manage a Docker container based dev environment
#----------------------------------------------------------------------

## dev-up: host: Start Webhooks development environment
.PHONY: dev-up
dev-up: dev-containers-up nats-up

## dev-containers-up: host: Start Webhooks development containers
.PHONY: dev-containers-up
dev-containers-up:
	@docker compose up -d

## dev-down: host: Stop Webhooks development environment and preserve persistent volumes
.PHONY: dev-down
dev-down:
	@docker compose --profile nats --profile pgadmin down

## dev-clean: host: Stop Webhooks development environment and delete persistent volumes
.PHONY: dev-clean
dev-clean:
	@docker compose --profile nats --profile pgadmin down --volumes

## dev-shell: host: Open bash shell in Webhooks development environment
.PHONY: dev-shell
dev-shell:
	@docker compose exec cds-dev bash -l
