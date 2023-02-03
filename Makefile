GO_BUILD_DIR=./build/server

MOBILE_API_GATEWAY_BINARY=mobileApiGateway


MOBILE_API_GATEWAY_MAIN=./server/services/mobile-api-gateway/cmd/grpc

## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	cd server; docker-compose up 

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_mobile_api_gateway
	@echo "Stopping docker images"
	cd server; docker-compose down
	@echo "Building and starting docker images..."
	cd server; docker-compose up --build 

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	cd server; docker-compose down
	@echo "Done!"

## build_broker: builds the broker binary as a linux executable
build_mobile_api_gateway:
	@echo "Building gateway binary..."
	export SERVER_MODE=development
	@echo ${SERVER_MODE}
	env GOOS=linux CGO_ENABLED=0 go build -o ${GO_BUILD_DIR}/${MOBILE_API_GATEWAY_BINARY} ${MOBILE_API_GATEWAY_MAIN}
	@echo "Done!"

## gen: generates TLS certificates 
gen-cert:
	cd cert; ./dev-cert-gen.sh; cd ..