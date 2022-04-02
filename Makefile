GOOS ?= $(shell go env GOOS)

run-serial:
	go run . --config=${CONFIG}

run-parallels:
	go run . --parallels --config=${CONFIG}

build:
	CGO_ENABLED=0 GOOS=$(GOOS) GO111MODULE=on go build -a -o goporter .

build-multi:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o goporter-linux-amd64 .
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 GO111MODULE=on go build -a -o goporter-linux-arm64 .
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 GO111MODULE=on go build -a -o goporter-darwin-amd64 .
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 GO111MODULE=on go build -a -o goporter-darwin-arm64 .
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 GO111MODULE=on go build -a -o goporter-windows-amd64 .
