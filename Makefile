
ifdef BUILDID
LDFLAGS=-ldflags "-X newtopia.BuildId $(BUILDID)"
endif

all: run

run: install
	$(GOPATH)/bin/genetics_loader

build: copy
	cd $(GOPATH)/src/genetics_loader; GOPATH=$(GOPATH) go build ./...

install: copy
	@echo Ensure you have called \'make updatedeps\'. Proceeding with install.
	GOPATH=$(GOPATH) GOBIN=$(GOPATH)/bin go install $(LDFLAGS) $(GOPATH)/src/genetics_loader/genetics_loader.go

copy: clean
	cp -R src/genetics_loader $(GOPATH)/src/genetics_loader;

updatedeps:
	# locally you can use `make copy; go get newtopia/...` Add -t to also fetch test dependencies.
	go list -f '{{join .Deps "\n"}}' ./... \
		| grep -v newtopia \
		| sort -u \
		| xargs go get -f -u -v

clean:
	rm -rf $(GOPATH)/src/genetics_loader

test: copy
	GOPATH=$(GOPATH) go test genetics_loader/... -v -p 1

format: fmt
fmt:
	go fmt ./src/...
