GO_BUILD_DIR=./build/server

MOBILE_API_GATEWAY_BINARY=mobileApiGateway
USER_SERVICE_BINARY=userService
MEDIA_SERVICE_BINARY=mediaService
POST_SERVICE_BINARY=postService
RELATIONSHIP_SERVICE_BINARY=relationshipService
NEWFEED_SERVICE_BINARY=newfeedService

MOBILE_API_GATEWAY_MAIN=./server/services/mobile-api-gateway/cmd/grpc
USER_SERVICE_MAIN=./server/services/user-service/cmd/grpc
MEDIA_SERVICE_MAIN=./server/services/media-service/cmd/grpc
POST_SERVICE_MAIN=./server/services/post-service/cmd/grpc
RELATIONSHIP_SERVICE_MAIN=./server/services/relationship-service/cmd/grpc
NEWFEED_SERVICE_MAIN=./server/services/newfeed-service/cmd/grpc

COMPOSE_DATABASES=-f common.yaml -f databases.yaml
COMPOSE_KAFKA=-f common.yaml -f kafka.yaml -f init-kafka.yaml
COMPOSE_FILES=-f common.yaml -f mobile-gateway.yaml -f user-service.yaml -f media-service.yaml -f post-service.yaml -f relationship-service.yaml -f newfeed-service.yaml

## up: starts all containers in the background without forcing build
up: format
	@echo "Starting Docker images..."
	cd server/infrastructure/docker-compose; docker compose ${COMPOSE_FILES} up

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: down format gen_cert gen_proto_server build_mobile_api_gateway build_user_service build_media_service build_post_service build_relationship_service build_newfeed_service
	@echo "Building and starting docker images..."
	cd server/infrastructure/docker-compose; docker compose ${COMPOSE_KAFKA} up --build -d
	cd server/infrastructure/docker-compose; docker compose ${COMPOSE_DATABASES} up --build -d
	sleep 60
	cd server/infrastructure/docker-compose; docker compose ${COMPOSE_FILES} up --build

ci_up_build: gen_cert gen_proto_server build_mobile_api_gateway build_user_service build_media_service build_post_service build_relationship_service build_newfeed_service
	cd server/infrastructure/docker-compose; docker compose ${COMPOSE_KAFKA} up --build -d
	cd server/infrastructure/docker-compose; docker compose ${COMPOSE_DATABASES} up --build -d
	sleep 30
	cd server/infrastructure/docker-compose; docker compose ${COMPOSE_FILES} up --build -d

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	cd server/infrastructure/docker-compose; docker compose ${COMPOSE_FILES} down --remove-orphans
	@echo "Done!"

## build_mobile_api_gateway: builds the mobile api gateway binary as a linux executable
build_mobile_api_gateway:
	@echo "Building gateway binary..."
	env GOOS=linux CGO_ENABLED=0 go build -o ${GO_BUILD_DIR}/${MOBILE_API_GATEWAY_BINARY} ${MOBILE_API_GATEWAY_MAIN}
	@echo "Done!"


## build_user_service: builds the user service binary as a linux executable
build_user_service:
	@echo "Building auth service binary..."
	env GOOS=linux CGO_ENABLED=0 go build -o ${GO_BUILD_DIR}/${USER_SERVICE_BINARY} ${USER_SERVICE_MAIN}
	@echo "Done!"

## build_media_service: builds the media service binary as a linux executable
build_media_service:
	@echo "Building media service binary..."
	env GOOS=linux CGO_ENABLED=0 go build -o ${GO_BUILD_DIR}/${MEDIA_SERVICE_BINARY} ${MEDIA_SERVICE_MAIN}
	@echo "Done!"

## build_post_service: builds the post service binary as a linux executable
build_post_service:
	@echo "Building post service binary..."
	env GOOS=linux CGO_ENABLED=0 go build -o ${GO_BUILD_DIR}/${POST_SERVICE_BINARY} ${POST_SERVICE_MAIN}
	@echo "Done!"

## build_relationship_service: builds the post service binary as a linux executable
build_relationship_service:
	@echo "Building relationship service binary..."
	env GOOS=linux CGO_ENABLED=0 go build -o ${GO_BUILD_DIR}/${RELATIONSHIP_SERVICE_BINARY} ${RELATIONSHIP_SERVICE_MAIN}
	@echo "Done!"

## build_newfeed_service: builds the post service binary as a linux executable
build_newfeed_service:
	@echo "Building newfeed service binary..."
	env GOOS=linux CGO_ENABLED=0 go build -o ${GO_BUILD_DIR}/${NEWFEED_SERVICE_BINARY} ${NEWFEED_SERVICE_MAIN}
	@echo "Done!"

## gen: generates TLS certificates 
gen_cert:
	@echo "Generating certs"
	cd cert; ./dev-cert-gen.sh; cd ..

# gen_proto: generates proto files
gen_proto: gen_proto_server gen_proto_client

gen_proto_server:
	@echo "Generating proto stubs for server"
	cd proto; ./generate.sh;

gen_proto_client:
	@echo "Generating proto stubs for mobile"
	cd mobile; ./gen-proto.sh;	

format:
	gofmt -s -w .

lint: lint-server

lint-server:
	golangci-lint run --out-format=github-actions --timeout=5m -- $$(go work edit -json | jq -c -r '[.Use[].DiskPath] | map_values(. + "/...")[]')

user-db-up:
	cd server/services/user-service; migrate -path db/migrations -database "postgresql://postgres:password@localhost:5433/users?sslmode=disable" up

user-db-down:
	cd server/services/user-service; migrate -path db/migrations -database "postgresql://postgres:password@localhost:5433/users?sslmode=disable" down

newfeed-db-up:
	cd server/services/newfeed-service; migrate -path db/migrations -database "cassandra://localhost:9042/newfeed?x-multi-statement=true&username=cassandra&password=cassandra" up

newfeed-db-down: 
	cd server/services/newfeed-service; migrate -path "db/migrations" -database "cassandra://localhost:9042/newfeed?x-multi-statement=true&username=cassandra&password=cassandra" down

count:
	git ls-files | grep -v .sum | grep -v .lock | grep -v asset | xargs wc -l