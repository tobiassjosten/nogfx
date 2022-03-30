PKGS := $(shell go list ./... | grep -v /vendor)
VERSION := $(shell git describe --always --long --dirty)

.PHONY: build
build: silence test
	@go build -ldflags="-X main.version=$(VERSION)" -o nogfx cmd/main.go

.PHONY: silence
silence:
	$(eval SILENT = "yes")

.PHONY: test
test: lint
	@[ "${SILENT}" = "yes" ] && (go test -race $(PKGS) || true) || go test -race $(PKGS)

.PHONY: lint
lint:
	@[ "${SILENT}" = "yes" ] && (golangci-lint run || true) || golangci-lint run

.PHONY: watch
watch:
	@clear; make test; fswatch -o -e ".*" -i ".*/[^.]*\\.go$$" -i ".*/testdata/.*" . | xargs -n1 -I{} sh -c 'clear; make test'

.EXPORT_ALL_VARIABLES: run
.PHONY: run
run:
	@./nogfx

.EXPORT_ALL_VARIABLES: debug
.PHONY: debug
debug:
	@dlv exec ./nogfx

.PHONY: docker
docker:
	@docker build . -t nogfx --build-arg VERSION=${VERSION}
