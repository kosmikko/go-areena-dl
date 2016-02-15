COMMIT = $$(git describe --always)
DEBUG_FLAG = $(if $(DEBUG),-debug)

install:
	@echo "====> Install go-areena-dl in $(GOPATH)/bin ..."
	@go install -ldflags "-X main.GitCommit=\"$(COMMIT)\""
