



.PHONY: proto-generate
proto-generate:
	./manage-proto.sh --generate --service file-service
	./manage-proto.sh --generate --service preprocessor-service

.PHONY: proto-clean
proto-clean:
	./manage-proto.sh --clean --service file-service
	./manage-proto.sh --clean --service preprocessor-service
