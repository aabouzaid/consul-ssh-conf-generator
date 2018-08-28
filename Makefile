APPNAME     := consul2ssh
DIST_DIR    := ./dist
PLATFORMS   := linux-386 linux-amd64 linux-arm
RELEASE     := $(shell git describe --abbrev=0 --tags | tr -d '[:alpha:]')
DOCKER_USER := $(or ${DOCKER_USERNAME}, root)
DOCKER_REPO := $(DOCKER_USER)/$(APPNAME)

#
# General.
deps:
	go generate
	go get -v -t ./...

build:
	go build -v -o $(APPNAME)_$(RELEASE)

clean:
	rm -rfv $(APPNAME)_$(RELEASE)

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
	env GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(DIST_DIR)/$(APPNAME)_$(RELEASE)_$(GOOS)_$(GOARCH)

cli_md5:
	cd $(DIST_DIR) && md5sum $(APPNAME)_$(RELEASE)_* | tee MD5SUM

cli_clean:
	rm -rfv $(DIST_DIR)
	mkdir -p $(DIST_DIR)

cli_release: deps cli_build cli_md5
