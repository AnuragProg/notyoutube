ARG GO_VERSION=1.23

##################################################
FROM golang:${GO_VERSION}-alpine AS builder

ARG PROTOC_GEN_GO_VERSION=v1.36.1
ARG PROTOC_GEN_GO_GRPC_VERSION=v1.5.1

WORKDIR /usr/app

RUN apk add --no-cache make
RUN apk add --no-cache protobuf
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@${PROTOC_GEN_GO_VERSION}
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@${PROTOC_GEN_GO_GRPC_VERSION}

COPY ./preprocessor-service/go.mod ./preprocessor-service/go.sum .
RUN go mod download && go mod verify
COPY ./preprocessor-service .

##################################################
FROM builder AS server-build
RUN make build BIN_NAME=app

##################################################
FROM scratch
COPY --from=server-build /usr/app/app .
ENTRYPOINT ["/app"]
