APPNAME        := consul2ssh
DIST_DIR       := ./dist
PLATFORMS      := linux-386 linux-amd64 linux-arm
RELEASE        := $(shell git describe --abbrev=0 --tags | tr -d '[:alpha:]')
DOCKER_USER    := $(or ${DOCKER_USERNAME}, root)
DOCKER_REPO    := $(DOCKER_USER)/$(APPNAME)
COVERAGEF_FILE := coverage.out

#
# General.
deps:
	go generate
	go get -v -t ./...

unittest:
	go test -v -covermode=count -coverprofile=$(COVERAGEF_FILE) ./...

build:
	go build -v -o $(APPNAME)_$(RELEASE)

clean:
	rm -rfv $(APPNAME)_$(RELEASE)

#
# Unit test.
unittest_coverage: deps unittest
	# Install required tools.
	go get -v golang.org/x/tools/cmd/cover github.com/mattn/goveralls

        # Get coverage report, and push it to coveralls.
	${GOPATH}/bin/goveralls -coverprofile=$(COVERAGEF_FILE) -service=travis-ci

#
# Docker image.
docker_build:
	# Build tagged and latest images.
	docker build -t $(DOCKER_REPO):$(RELEASE) -t $(DOCKER_REPO) .

docker_list:
	docker images

docker_push:
	# Push images.
	docker push $(DOCKER_REPO)

docker_release: docker_build docker_list docker_push

#
# CLI bin.
cli_build: cli_clean $(PLATFORMS)
$(PLATFORMS):
	$(eval GOOS := $(firstword $(subst -, ,$@)))
	$(eval GOARCH := $(lastword $(subst -, ,$@)))
	$(eval BIN_FILE := $(APPNAME)_$(RELEASE)_$(GOOS)_$(GOARCH))
	env GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(DIST_DIR)/$(BIN_FILE)
	gzip $(DIST_DIR)/$(BIN_FILE)

cli_md5:
	cd $(DIST_DIR) && md5sum $(APPNAME)_$(RELEASE)_* | tee MD5SUM

cli_clean:
	rm -rfv $(DIST_DIR)
	mkdir -p $(DIST_DIR)

cli_release: deps cli_build cli_md5
