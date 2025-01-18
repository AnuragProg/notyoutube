

.PHONY: help
help:
	@echo "proto-generate: Generate/Update proto files"
	@echo "proto-clean: Clean generated proto files"


.PHONY: proto-generate
proto-generate:
	./manage-proto.sh --generate --service file-service
	./manage-proto.sh --generate --service preprocessor-service
	./manage-proto.sh --generate --service dag-scheduler-service

.PHONY: proto-clean
proto-clean:
	./manage-proto.sh --clean --service file-service
	./manage-proto.sh --clean --service preprocessor-service
	./manage-proto.sh --clean --service dag-scheduler-service

