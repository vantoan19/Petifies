GO_BUILD_DIR=./build/server

MOBILE_API_GATEWAY_BINARY=mobileApiGateway
AUTH_SERVICE_BINARY=authService

MOBILE_API_GATEWAY_MAIN=./server/services/mobile-api-gateway/cmd/grpc
AUTH_SERVICE_MAIN=./server/services/user-services/auth-service/cmd/grpc

## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	cd server; docker-compose up 

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_mobile_api_gateway build_auth_service
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
	env GOOS=linux CGO_ENABLED=0 go build -o ${GO_BUILD_DIR}/${MOBILE_API_GATEWAY_BINARY} ${MOBILE_API_GATEWAY_MAIN}
	@echo "Done!"

## build_broker: builds the broker binary as a linux executable
build_auth_service:
	@echo "Building auth service binary..."
	env GOOS=linux CGO_ENABLED=0 go build -o ${GO_BUILD_DIR}/${AUTH_SERVICE_BINARY} ${AUTH_SERVICE_MAIN}
	@echo "Done!"

## gen: generates TLS certificates 
gen-cert:
	cd cert; ./dev-cert-gen.sh; cd ..