GO_BUILD_DIR=./build/server

MOBILE_API_GATEWAY_BINARY=mobileApiGateway
USER_SERVICE_BINARY=userService

MOBILE_API_GATEWAY_MAIN=./server/services/mobile-api-gateway/cmd/grpc
USER_SERVICE_MAIN=./server/services/user-service/cmd/grpc

COMPOSE_FILES=-f common.yaml -f mobile-gateway.yaml -f user-service.yaml 

## up: starts all containers in the background without forcing build
up: format
	@echo "Starting Docker images..."
	cd server/infrastructure/docker-compose; docker compose ${COMPOSE_FILES} up

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: format gen_cert gen_proto build_mobile_api_gateway build_user_service
	@echo "Stopping docker images"
	cd server/infrastructure/docker-compose; docker compose ${COMPOSE_FILES} down
	@echo "Building and starting docker images..."
	cd server/infrastructure/docker-compose; docker compose ${COMPOSE_FILES} up --build

ci_up_build: gen_cert gen_proto build_mobile_api_gateway build_user_service
	cd server/infrastructure/docker-compose; docker compose ${COMPOSE_FILES} up --build -d

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	cd server/infrastructure/docker-compose; docker compose ${COMPOSE_FILES} down
	@echo "Done!"

## build_broker: builds the broker binary as a linux executable
build_mobile_api_gateway:
	@echo "Building gateway binary..."
	env GOOS=linux CGO_ENABLED=0 go build -o ${GO_BUILD_DIR}/${MOBILE_API_GATEWAY_BINARY} ${MOBILE_API_GATEWAY_MAIN}
	@echo "Done!"


## build_broker: builds the broker binary as a linux executable
build_user_service:
	@echo "Building auth service binary..."
	env GOOS=linux CGO_ENABLED=0 go build -o ${GO_BUILD_DIR}/${USER_SERVICE_BINARY} ${USER_SERVICE_MAIN}
	@echo "Done!"

## gen: generates TLS certificates 
gen_cert:
	@echo "Generating certs"
	cd cert; ./dev-cert-gen.sh; cd ..

# gen_proto: generates proto files
gen_proto:
	@echo "Generating proto stubs"
	cd proto; ./generate.sh;

format:
	gofmt -s -w .