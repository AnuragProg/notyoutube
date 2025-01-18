#!/bin/bash

readonly PROTO_ROOT_DIR="proto"
readonly FILE_SERVICE_IDENT="file-service"
readonly PREPROCESSOR_SERVICE_IDENT="preprocessor-service"
readonly DAG_SCHEDULER_SERVICE_IDENT="dag-scheduler-service"

# INFO: Associative maps for mapping proto files to their generated code directories.
#       - Key: Proto file name (e.g., "raw_video_metadata.proto").
#       - Value: Semicolon-separated paths for `go_out` and `go-grpc_out` directories.
#       - Note: Currently, all generated code is placed in one directory to avoid import errors.

# File service specific mappings.
# Proto files in following services are translated to paths as follows:
#   File-Service:
#       - proto-file[key] (raw_video_metadata.proto) => proto/file-service/raw_video_metadata.proto
#       - gen-dir[oneof value] (types/mq) => file-service/types/mq
#   Same goes for every other service with dir changing as per service...
declare -A file_service_proto_file_to_generated_dir=(
    ["raw_video_metadata.proto"]="types/mq;types/mq"
    # ["raw_video_service.proto"]="types/raw_video_service;repository_impl/raw_video_service" # NOTE: refer to 1.
    ["raw_video_service.proto"]="repository_impl/raw_video_service;repository_impl/raw_video_service" # NOTE: refer to 1.
)
declare -A preprocessor_service_proto_file_to_generated_dir=(
    ["raw_video_metadata.proto"]="types/mq;types/mq"
    ["dag.proto"]="types/mq;types/mq"
    ["raw_video_service.proto"]="repository_impl/raw_video_service;repository_impl/raw_video_service"
)
declare -A dag_scheduler_service_proto_file_to_generated_dir=(
    ["dag.proto"]="types/mq;types/mq"
    ["dag_service.proto"]="repository_impl/dag_service;repository_impl/dag_service"
)

readonly file_service_proto_to_generated
readonly preprocessor_service_proto_to_generated
readonly dag_scheduler_service_proto_file_to_generated_dir

generate_proto(){
    local proto_path=$1
    local go_out=$2
    local go_grpc_out=$3
    local proto_file=$4

    echo "Generating: $proto_file at $go_out and $go_grpc_out..."
    protoc --proto_path=$proto_path \
        --go_opt=paths=source_relative \
        --go_out=$go_out \
        --go-grpc_opt=paths=source_relative \
        --go-grpc_out=$go_grpc_out \
        $proto_file
    echo "Generated: $proto_file!"
}

generate_proto_for_service() {
    local service=$1

    # different directories for storing types and grpc methods
    local go_out_dir=""
    local go_grpc_out_dir=""
    # ident are just the values stored in the associative maps for go and grpc output
    local go_out_ident=""
    local go_grpc_out_ident=""
    case $service in
        $FILE_SERVICE_IDENT)
            for proto_file in ${!file_service_proto_file_to_generated_dir[@]}; do
                IFS=";" read -r go_out_ident go_grpc_out_ident <<< "${file_service_proto_file_to_generated_dir[$proto_file]}"
                go_out_dir="$FILE_SERVICE_IDENT/$go_out_ident"
                go_grpc_out_dir="$FILE_SERVICE_IDENT/$go_grpc_out_ident"
                generate_proto "./$PROTO_ROOT_DIR/$FILE_SERVICE_IDENT" $go_out_dir $go_grpc_out_dir "./$PROTO_ROOT_DIR/$FILE_SERVICE_IDENT/$proto_file" &
            done
            wait
            ;;
        $PREPROCESSOR_SERVICE_IDENT)
            for proto_file in ${!preprocessor_service_proto_file_to_generated_dir[@]}; do
                IFS=";" read -r go_out_ident go_grpc_out_ident <<< "${preprocessor_service_proto_file_to_generated_dir[$proto_file]}"
                go_out_dir="$PREPROCESSOR_SERVICE_IDENT/$go_out_ident"
                go_grpc_out_dir="$PREPROCESSOR_SERVICE_IDENT/$go_grpc_out_ident"
                generate_proto "./$PROTO_ROOT_DIR/$PREPROCESSOR_SERVICE_IDENT" $go_out_dir $go_grpc_out_dir "./$PROTO_ROOT_DIR/$PREPROCESSOR_SERVICE_IDENT/$proto_file" &
            done
            wait
            ;;
        $DAG_SCHEDULER_SERVICE_IDENT)
            for proto_file in ${!dag_scheduler_service_proto_file_to_generated_dir[@]}; do
                IFS=";" read -r go_out_ident go_grpc_out_ident <<< "${dag_scheduler_service_proto_file_to_generated_dir[$proto_file]}"
                go_out_dir="$DAG_SCHEDULER_SERVICE_IDENT/$go_out_ident"
                go_grpc_out_dir="$DAG_SCHEDULER_SERVICE_IDENT/$go_grpc_out_ident"
                generate_proto "./$PROTO_ROOT_DIR/$DAG_SCHEDULER_SERVICE_IDENT" $go_out_dir $go_grpc_out_dir "./$PROTO_ROOT_DIR/$DAG_SCHEDULER_SERVICE_IDENT/$proto_file" &
            done
            wait
            ;;
        *)
            echo "Error: unknown service" $servce "requested" >&2
            exit 1
            ;;
    esac
}

clean_proto_for_service(){
    local service=$1

    # different directories for storing types and grpc methods
    local go_out_dir=""
    local go_grpc_out_dir=""
    # ident are just the values stored in the associative maps for go and grpc output
    local go_out_ident=""
    local go_grpc_out_ident=""
    case $service in
        $FILE_SERVICE_IDENT)
            for proto_file in ${!file_service_proto_file_to_generated_dir[@]}; do
                IFS=";" read -r go_out_ident go_grpc_out_ident <<< "${file_service_proto_file_to_generated_dir[$proto_file]}"
                go_out_files="$FILE_SERVICE_IDENT/$go_out_ident/*.pb.go"
                go_grpc_out_files="$FILE_SERVICE_IDENT/$go_grpc_out_ident/*.pb.go"
                echo "Deleting: $go_out_files & $go_grpc_out_files contents"
                rm -r $go_out_files $go_grpc_out_files 2>/dev/null || echo "No contents to remove or some files could not be deleted in $go_out_files & $go_grpc_out_files"
            done
            ;;
        $PREPROCESSOR_SERVICE_IDENT)
            for proto_file in ${!preprocessor_service_proto_file_to_generated_dir[@]}; do
                IFS=";" read -r go_out_ident go_grpc_out_ident <<< "${preprocessor_service_proto_file_to_generated_dir[$proto_file]}"
                go_out_files="$PREPROCESSOR_SERVICE_IDENT/$go_out_ident/*.pb.go"
                go_grpc_out_files="$PREPROCESSOR_SERVICE_IDENT/$go_grpc_out_ident/*.pb.go"
                echo "Deleting: $go_out_files & $go_grpc_out_files contents"
                rm -r $go_out_files $go_grpc_out_files 2>/dev/null || echo "No contents to remove or some files could not be deleted in $go_out_files & $go_grpc_out_files"
            done
            ;;
        $DAG_SCHEDULER_SERVICE_IDENT)
            for proto_file in ${!dag_scheduler_service_proto_file_to_generated_dir[@]}; do
                IFS=";" read -r go_out_ident go_grpc_out_ident <<< "${dag_scheduler_service_proto_file_to_generated_dir[$proto_file]}"
                go_out_files="$DAG_SCHEDULER_SERVICE_IDENT/$go_out_ident/*.pb.go"
                go_grpc_out_files="$DAG_SCHEDULER_SERVICE_IDENT/$go_grpc_out_ident/*.pb.go"
                echo "Deleting: $go_out_files & $go_grpc_out_files contents"
                rm -r $go_out_files $go_grpc_out_files 2>/dev/null || echo "No contents to remove or some files could not be deleted in $go_out_files & $go_grpc_out_files"
            done
            ;;
        *)
            echo "Error: unknown service" $service "requested" >&2
            exit 1
            ;;
    esac
}

###### BELOW IS FOR USER INPUT EXTRACTION AND VALIDATION ONLY ######
generate=false
clean=false
service=""

while [[ $# -gt 0 ]]; do
    case $1 in
        -g|--generate)
            generate=true
            shift
            ;;
        -c|--clean)
            clean=true
            shift
            ;;
        -s|--service)
            service=$2
            shift
            shift
            ;;
        *)
            echo "Error: invalid option" $1 >&2 
            exit 1
    esac
done

if $generate && $clean; then
    echo "Error: conflicting commands, can't generate and clean at the same time" >&2
    exit 1
fi
if ! $generate && ! $clean; then
    echo "Error: conflicting commands, noop situation" >&2
    exit 1
fi

if $generate; then
    generate_proto_for_service $service
elif $clean; then
    clean_proto_for_service $service
fi
