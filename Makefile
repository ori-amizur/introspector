TAG := $(or $(TAG),latest)
AGENT := $(or ${AGENT},quay.io/ocpmetal/agent:$(TAG))
CONNECTIVITY_CHECK := $(or ${CONNECTIVITY_CHECK},quay.io/ocpmetal/connectivity_check:$(TAG))
INVENTORY := $(or ${INVENTORY},quay.io/ocpmetal/inventory:$(TAG))
FREE_ADDRESSES := $(or ${FREE_ADDRESSES},quay.io/ocpmetal/free_addresses:$(TAG))
LOGS_SENDER := $(or ${LOGS_SENDER},quay.io/ocpmetal/logs_sender:$(TAG))

all: build

.PHONY: build clean build-image push subsystem agent-build hardware-info-build connectivity-check-build inventory-build logs-sender-build
build: agent-build hardware-info-build connectivity-check-build inventory-build free-addresses-build logs-sender-build

agent-build : src/agent/main/main.go
	mkdir -p build
	CGO_ENABLED=0 go build -o build/agent src/agent/main/main.go

connectivity-check-build : src/connectivity_check/main/main.go
	mkdir -p build
	CGO_ENABLED=0 go build -o build/connectivity_check src/connectivity_check/main/main.go

inventory-build : src/inventory
	mkdir -p build
	CGO_ENABLED=0 go build -o build/inventory src/inventory/main/main.go

free-addresses-build: src/free_addresses
	mkdir -p build
	CGO_ENABLED=0 go build -o build/free_addresses src/free_addresses/main/main.go

logs-sender-build: src/logs_sender
	mkdir -p build
	CGO_ENABLED=0 go build -o build/logs_sender src/logs_sender/main/main.go

clean:
	rm -rf build subsystem/logs

build-image: unittest build
	docker build -f Dockerfile.agent . -t $(AGENT)
	docker build -f Dockerfile.connectivity_check . -t $(CONNECTIVITY_CHECK)
	docker build -f Dockerfile.inventory . -t $(INVENTORY)
	docker build -f Dockerfile.free_addresses . -t $(FREE_ADDRESSES)
	docker build -f Dockerfile.logs_sender . -t $(LOGS_SENDER)

push: build-image subsystem
	docker push $(AGENT)
	docker push $(CONNECTIVITY_CHECK)
	docker push $(INVENTORY)
	docker push $(FREE_ADDRESSES)
	docker push $(LOGS_SENDER)

unittest:
	go test -v $(shell go list ./... | grep -v subsystem) -cover

subsystem: build-image
	cd subsystem; docker-compose up -d
	go test -v ./subsystem/... -count=1 -ginkgo.focus=${FOCUS} -ginkgo.v -ginkgo.skip="system-test" || ( cd subsystem; docker-compose down && /bin/false)
	cd subsystem; docker-compose down

generate:
	go generate $(shell go list ./...)

go-import:
	goimports -w -l .
