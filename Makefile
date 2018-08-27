APPNAME   := consul2ssh
DIST_DIR  := ./dist
PLATFORMS := linux-386 linux-amd64 linux-arm
RELEASE   := $(shell git describe --abbrev=0 --tags | tr -d '[:alpha:]')

deps:
	go generate
	go get -v -t ./...

build: get
	go build -v -o $(APPNAME)_$(RELEASE)

clean_build:
	rm -rfv $(APPNAME)_$(RELEASE)

build_dist: clean_dist $(PLATFORMS)
$(PLATFORMS):
	$(eval GOOS := $(firstword $(subst -, ,$@)))
	$(eval GOARCH := $(lastword $(subst -, ,$@)))
	env GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(DIST_DIR)/$(APPNAME)_$(RELEASE)_$(GOOS)_$(GOARCH)

clean_dist:
	rm -rfv $(DIST_DIR)
	mkdir -p $(DIST_DIR)

md5: dist
	cd $(DIST_DIR) && md5sum $(APPNAME)_$(RELEASE)_* | tee MD5SUM

release: deps build_dist md5

clean: clean_dist clean_build
