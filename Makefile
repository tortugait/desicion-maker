# suppress output, run `make XXX V=` to be verbose
V := @

-include .env
export $(shell sed 's/=.*//' .env 2>/dev/null)

# common
TARGETS = $(notdir $(shell find ./cmd/* -maxdepth 0 -type d))

# build
OUT_DIR = ./bin
LD_FLAGS = -ldflags "-s -w"
ifeq ($(shell uname),Linux)
	ifneq ($(shell ls -la /proc/self/map_files | grep musl),)
		GO_TAGS = musl # add musl tag for musl-based linux only
	endif
endif

.PHONY: build
build: clean $(TARGETS)

$(TARGETS):
	@echo BUILDING $@
	$(V)CGO_ENABLED=1 go build -o ${OUT_DIR}/$@ ${LD_FLAGS} ./cmd/$@
	@echo DONE

.PHONY: clean
clean:
	@echo REMOVING ${OUT_DIR}
	$(V)rm -rf ${OUT_DIR}

.PHONY: lint
lint:
	$(V)golangci-lint run

# note: one may see linker warning, see https://github.com/golang/go/issues/61229#issuecomment-1988965927
.PHONY: test
test: GO_TEST_FLAGS += -race
test:
	$(V)go test ${GO_TEST_FLAGS} ./...

.PHONY: gen
gen:
	$(V)go run github.com/vektra/mockery/v2@v2.45.0
	$(V)go generate -x ./...
