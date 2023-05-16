.PHONY: test
test:
	go test -count=1 ./core/... -v

.PHONY: install
install:
	GOBIN=$(CURDIR)/bin go install ./...
