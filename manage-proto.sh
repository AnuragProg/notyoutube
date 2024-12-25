#!/bin/bash

readonly PROTO_ROOT_DIR="proto"
readonly FILE_SERVICE_IDENT="file-service"
readonly PREPROCESSOR_SERVICE_IDENT="preprocessor-service"

# key=<proto file name>; value=<location inside of service for generation>
declare -A file_service_proto_file_to_generated_dir=(
    ["raw_video_metadata.proto"]="types/mq"
)
declare -A preprocessor_service_proto_file_to_generated_dir=(
    ["raw_video_metadata.proto"]="types/mq"
    ["dag.proto"]="types/mq"
)

readonly file_service_proto_to_generated
readonly preprocessor_service_proto_to_generated


generate_proto(){
    local proto_path=$1
    local out_dir=$2
    local proto_file=$3

    echo "Generating: $proto_file..."
    protoc --proto_path=$proto_path \
        --go_opt=paths=source_relative \
        --go_out=$out_dir \
        $proto_file
    echo "Generated: $proto_file!"
}

generate_proto_for_service() {
    local service=$1

    local out_dir=""
    case $service in
        $FILE_SERVICE_IDENT)
            for proto_file in ${!file_service_proto_file_to_generated_dir[@]}; do
                out_dir="$FILE_SERVICE_IDENT/${file_service_proto_file_to_generated_dir[$proto_file]}"
                generate_proto "./$PROTO_ROOT_DIR/$FILE_SERVICE_IDENT" $out_dir "./$PROTO_ROOT_DIR/$FILE_SERVICE_IDENT/$proto_file" &
            done
            wait
            ;;
        $PREPROCESSOR_SERVICE_IDENT)
            for proto_file in ${!preprocessor_service_proto_file_to_generated_dir[@]}; do
                out_dir="$PREPROCESSOR_SERVICE_IDENT/${preprocessor_service_proto_file_to_generated_dir[$proto_file]}"
                generate_proto "./$PROTO_ROOT_DIR/$PREPROCESSOR_SERVICE_IDENT" $out_dir "./$PROTO_ROOT_DIR/$PREPROCESSOR_SERVICE_IDENT/$proto_file" &
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

    local out_dir=""
    case $service in
        $FILE_SERVICE_IDENT)
            for proto_file in ${!file_service_proto_file_to_generated_dir[@]}; do
                generated_proto_files="$FILE_SERVICE_IDENT/${file_service_proto_file_to_generated_dir[$proto_file]}/*.pb.go"
                echo "Deleting: $generated_proto_files contents"
                rm -r $generated_proto_files 2>/dev/null || echo "No contents to remove or some files could not be deleted in $generated_proto_files"
            done
            ;;
        $PREPROCESSOR_SERVICE_IDENT)
            for proto_file in ${!preprocessor_service_proto_file_to_generated_dir[@]}; do
                generated_proto_files="$PREPROCESSOR_SERVICE_IDENT/${preprocessor_service_proto_file_to_generated_dir[$proto_file]}/*.pb.go"
                echo "Deleting: $generated_proto_files contents"
                rm -r $generated_proto_files 2>/dev/null || echo "No contents to remove or some files could not be deleted in $generated_proto_files"
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
