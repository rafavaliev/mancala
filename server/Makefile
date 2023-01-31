NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m

BINARY_NAME?=mancala
GO_LINKER_FLAGS=-ldflags "-s"
DIR_OUT=$(CURDIR)/out

.PHONY: all clean deps build install

all: clean deps build install

clean:
	@printf "$(OK_COLOR)==> Cleaning project$(NO_COLOR)\n"
	rm -f ${DIR_OUT}/${BINARY_NAME}

deps:
	@printf "$(OK_COLOR)==> Installing deps$(NO_COLOR)\n"
	@go mod tidy
	@go mod vendor

build:
	@printf "$(OK_COLOR)==> Building binary$(NO_COLOR)\n"
	@go build ./cmd/main/main.go


#---------------
#-- tests
#---------------
.PHONY: tests test-unit
tests: test-unit

test-unit: tools.format tools.vet
	@printf "$(OK_COLOR)==> Unit Testing$(NO_COLOR)\n"
	@go test -v -race -cover ./...

#---------------
#-- tools
#---------------
.PHONY: tools tools.golint tools.goimports tools.format tools.vet
tools: tools.goimports tools.format tools.lint tools.vet

tools.goimports:
	@command -v goimports >/dev/null ; if [ $$? -ne 0 ]; then \
		echo "--> installing goimports"; \
		@go get golang.org/x/tools/cmd/goimports; \
	fi
	@echo "$(OK_COLOR)==> checking imports 'goimports' tool$(NO_COLOR)"
	@goimports -l -w *.go cmd pkg internal &>/dev/null | grep ".*\.go"; if [ "$$?" = "0" ]; then exit 1; fi

tools.format:
	@echo "$(OK_COLOR)==> formatting code with 'gofmt' tool$(NO_COLOR)"
	@gofmt -l -s -w *.go cmd pkg internal | grep ".*\.go"; if [ "$$?" = "0" ]; then exit 1; fi

tools.lint:
	@command -v golint >/dev/null ; if [ $$? -ne 0 ]; then \
		echo "--> installing golint"; \
		@go get github.com/golang/lint/golint; \
	fi
	@echo "$(OK_COLOR)==> checking code style with 'golint' tool$(NO_COLOR)"
	@go list ./... | xargs -n 1 golint -set_exit_status

tools.vet:
	@echo "$(OK_COLOR)==> checking code correctness with 'go vet' tool$(NO_COLOR)"
	@go vet ./...