#!/bin/bash
MODULE_PREFIX="github.com/vantoan19/Petifies/proto"
FOLDERS=("auth-gateway/v1" "public-gateway/v1" "user-service/v1" "media-service/v1" "post-service/v1" "relationship-service/v1")

protoc --proto_path=. common/common.proto \
    --go_out=paths=source_relative:. \
    --validate_out=lang=go,paths=source_relative:. 
for folder in "${FOLDERS[@]}"
do
    mkdir -p $folder
    protoc --proto_path=. $folder/*proto \
        --go_out=paths=source_relative:. \
        --go-grpc_out=paths=source_relative:. \
        --go-grpc_opt=require_unimplemented_servers=false \
        --validate_out=lang=go,paths=source_relative:. 
done;
