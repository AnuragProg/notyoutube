##################################################
FROM golang:1.23-alpine as builder

WORKDIR /usr/app

RUN apk add --no-cache make
RUN apk add --no-cache protobuf
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

COPY ./file-service/go.mod ./file-service/go.sum .
RUN go mod download && go mod verify
COPY ./file-service .

##################################################
FROM builder as server-build
RUN make build BIN_NAME=app

##################################################
FROM scratch
COPY --from=server-build /usr/app/app .
ENTRYPOINT ["/app"]
