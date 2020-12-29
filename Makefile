.PHONY: mod
mod:
	go mod download

.PHONY: vendor
vendor:
	go mod vendor

.PHONY: test
test:
	go test ./...
