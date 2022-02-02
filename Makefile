PROTO_OUT ?= pkg/pb

# Current dir
DIR ?= $(shell pwd)

# Protoc command dockerized
PROTOC_CMD ?= docker run --rm -v $(DIR):$(DIR) -w $(DIR) thethingsindustries/protoc

.PHONY: compile-proto-go
compile-proto-go: ## Compile the protobuf contracts to generated Go code
	@rm -rf $(PROTO_OUT)
	@mkdir -p $(PROTO_OUT)
	@find proto -name "*.proto" -type f -exec $(PROTOC_CMD) -Iproto/ \
	     --go_opt=paths=source_relative \
	     --go_out=plugins=grpc:$(PROTO_OUT) {} \;

.PHONY: docker-up
docker-up:
	docker-compose -f docker/docker-compose.yaml up

.PHONY: docker-down
docker-down:
	docker-compose -f docker/docker-compose.yaml down
	docker system prune --volumes --force

.PHONY: up
up:
	go run -race cmd/ingester/main.go

.PHONY: docker-build
docker-build:
	docker build -t claudioed/ingester:latest .

.PHONY: docker-push
docker-push:
	docker push claudioed/ingester:latest