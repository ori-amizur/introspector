TAG := $(or $(TAG),latest)
AGENT := $(or ${AGENT},quay.io/ocpmetal/agent:$(TAG))
CONNECTIVITY_CHECK := $(or ${CONNECTIVITY_CHECK},quay.io/ocpmetal/connectivity_check:$(TAG))
INVENTORY := $(or ${INVENTORY},quay.io/ocpmetal/inventory:$(TAG))
HARDWARE_INFO := $(or ${HARDWARE_INFO},quay.io/ocpmetal/hardware_info:$(TAG))
FREE_ADDRESSES := $(or ${FREE_ADDRESSES},quay.io/ocpmetal/free_addresses:$(TAG))
CONTAINER_RUNTIME := $(shell command -v podman 2> /dev/null || echo docker)
UID = $(shell id -u)

all: build

.PHONY: build clean build-image push subsystem agent-build hardware-info-build connectivity-check-build inventory-build
build: deps generate-from-swagger generate agent-build hardware-info-build connectivity-check-build inventory-build free-addresses-build

deps:
	GOSUMDB=off go mod download

agent-build : src/agent/main/main.go
	mkdir -p build
	CGO_ENABLED=0 go build -o build/agent src/agent/main/main.go

hardware-info-build : src/hardware_info/main/main.go
	mkdir -p build
	CGO_ENABLED=0 go build -o build/hardware_info src/hardware_info/main/main.go

connectivity-check-build : src/connectivity_check/main/main.go
	mkdir -p build
	CGO_ENABLED=0 go build -o build/connectivity_check src/connectivity_check/main/main.go

inventory-build : src/inventory
	mkdir -p build
	CGO_ENABLED=0 go build -o build/inventory src/inventory/main/main.go

free-addresses-build: src/free_addresses
	mkdir -p build
	CGO_ENABLED=0 go build -o build/free_addresses src/free_addresses/main/main.go

clean:
	-rm -rf build subsystem/logs generated

build-image: build unittest
	$(CONTAINER_RUNTIME) build -f Dockerfile.agent . -t $(AGENT)
	$(CONTAINER_RUNTIME) build -f Dockerfile.connectivity_check . -t $(CONNECTIVITY_CHECK)
	$(CONTAINER_RUNTIME) build -f Dockerfile.inventory . -t $(INVENTORY)
	$(CONTAINER_RUNTIME) build -f Dockerfile.hardware_info . -t $(HARDWARE_INFO)
	$(CONTAINER_RUNTIME) build -f Dockerfile.free_addresses . -t $(FREE_ADDRESSES)

push: subsystem
	$(CONTAINER_RUNTIME) push $(AGENT)
	$(CONTAINER_RUNTIME) push $(CONNECTIVITY_CHECK)
	$(CONTAINER_RUNTIME) push $(INVENTORY)
	$(CONTAINER_RUNTIME) push $(HARDWARE_INFO)
	$(CONTAINER_RUNTIME) push $(FREE_ADDRESSES)

unittest:
	go test -v $(shell go list ./... | grep -v subsystem) -cover

subsystem: build-image
	cd subsystem; docker-compose up -d
	go test -v ./subsystem/... -count=1 -ginkgo.focus=${FOCUS} -ginkgo.v -ginkgo.skip="system-test" || ( cd subsystem; docker-compose down && /bin/false)
	cd subsystem; docker-compose down

generate:
	go generate $(shell go list ./...)

generate-from-swagger: clean
	mkdir -p ./generated/bm-inventory
	cp $(shell go list -m -f={{.Dir}} github.com/filanov/bm-inventory)/swagger.yaml ./generated/bm-inventory/swagger.yaml
	chown $(UID):$(UID) ./generated/bm-inventory/swagger.yaml
	$(CONTAINER_RUNTIME) run -u $(UID):$(UID) -v $(PWD):$(PWD):rw,Z -v /etc/passwd:/etc/passwd -w $(PWD) \
		quay.io/goswagger/swagger:v0.24.0 generate client --template=stratoscale -f ./generated/bm-inventory/swagger.yaml \
		--template-dir=/templates/contrib -t $(PWD)/generated/bm-inventory

go-import:
	goimports -w -l .

