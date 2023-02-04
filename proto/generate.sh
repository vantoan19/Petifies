#!/bin/bash
MODULE_PREFIX="github.com/vantoan19/Petifies/proto"
FOLDERS=("auth-gateway" "public-gateway" "auth-service")


for folder in "${FOLDERS[@]}"
do
    mkdir -p $folder
    protoc --proto_path=$folder $folder/*proto --go_out=$folder --go_opt=module=$MODULE_PREFIX/$folder --go-grpc_out=$folder --go-grpc_opt=module=$MODULE_PREFIX/$folder --go-grpc_opt=require_unimplemented_servers=false
done;
