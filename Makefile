.PHONY: build deps #release run test
	
# directory to output build
DIST_DIR=./dist
# get the date and time to use as a buildstamp
DATE=$$(date '+%Y-%m-%d')
TIME=$$(date '+%I:%M:%S%p')
LDFLAGS="-s -w -X main.buildDate=$(DATE) -X main.buildTime=$(TIME)"

build:
	@go build --ldflags=$(LDFLAGS) -o $(DIST_DIR)/log_agg main.go
	
deps:
	@go get github.com/kardianos/govendor
	@govendor sync
	@go get -v github.com/mitchellh/gox
	# @go get -t -v ./...
