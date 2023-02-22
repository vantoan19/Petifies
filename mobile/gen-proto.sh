OUT_DIR_DART=./lib/src/proto
PROTO_DIR=../proto
PROTO_FILES=("common/common.proto" "public-gateway/v1/public-gateway.v1.proto" "auth-gateway/v1/auth-gateway.v1.proto"
 "google/protobuf/timestamp.proto" "google/protobuf/duration.proto" "validate/validate.proto")

for proto_file in "${PROTO_FILES[@]}"
do 
    protoc --proto_path=$PROTO_DIR $proto_file \
    --dart_out=grpc:$OUT_DIR_DART 
done;
