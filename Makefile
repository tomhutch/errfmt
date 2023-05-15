.PHONY: test
test:
	go test -count=1 ./core/... -v

test-fixer:
	go test -count=1 ./errfmtfixer/... -v

install:
	GOBIN=$(CURDIR)/bin go install ./...
