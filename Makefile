GO_PROJECT_NAME := poh-golang

# GO commands
go_build:
	@echo "\n....Building $(GO_PROJECT_NAME)"
	go build -ldflags "-s -w" -o ./bin/ ./cmd/poh-golang
	go build -ldflags "-s -w" -o ./bin/perry main.go

go_dep_install:
	@echo "\n....Installing dependencies for $(GO_PROJECT_NAME)...."
	go get .

go_run:
	@echo "\n....Running $(GO_PROJECT_NAME)...."
	$(GOPATH)/bin/$(GO_PROJECT_NAME)

test:
	@echo "\n....Running tests for $(GO_PROJECT_NAME)...."
	LOG_IGNORE=1 go test ./pkg/poh_hash

# Project rules
build:
	$(MAKE) go_build

bench:
	LOG_IGNORE=1 go test -bench=. ./pkg/poh_hash -count 5 -benchmem
	
benchstat:
	LOG_IGNORE=1 go test -bench=. ./pkg/poh_hash -count 5 -benchmem | tee benchmark.out
	benchstat benchmark.out

prof:
	LOG_IGNORE=1 go test -cpuprofile cpu.prof -memprofile mem.prof -bench=. ./pkg/poh_hash

race:
	LOG_IGNORE=1 go test -race ./pkg/poh_hash

run:
	$(MAKE) go_build
	$(MAKE) go_run

clean:
	rm -f ./bin/*

docker:
	@echo "\n....Building latest docker image and uploading to GCR ...."
	$(MAKE) test
	docker buildx build --push --platform linux/arm/v7,linux/arm64/v8,linux/amd64 --tag calacode/$(GO_PROJECT_NAME):latest .

.PHONY: docker db_seed go_build go_dep_install go_prep_install go_run build run restart test
