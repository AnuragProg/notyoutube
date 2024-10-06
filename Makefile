PROTO_SRC = ./proto

FILE_SERVICE_IDENT = file-service

FILE_SERVICE_PROTO_DIR = $(PROTO_SRC)/$(FILE_SERVICE_IDENT)
FILE_SERVICE_PROTO_MQ_FILE = $(FILE_SERVICE_PROTO_DIR)/mq.proto
FILE_SERVICE_PROTO_MQ_OUT_DIR = ./$(FILE_SERVICE_IDENT)/types/mq
FILE_SERVICE_PROTO_MQ_OUT_FILE = $(FILE_SERVICE_PROTO_MQ_OUT_DIR)/mq.pb.go


$(FILE_SERVICE_PROTO_MQ_OUT_FILE): $(FILE_SERVICE_PROTO_MQ_FILE)
	@echo "Building $@ because $< changed"
	@mkdir -p $(FILE_SERVICE_PROTO_MQ_OUT_DIR) # creating dir if not created already
	# @-rm $(FILE_SERVICE_PROTO_MQ_OUT_FILE) # deleting previously created file
	@protoc --proto_path=$(FILE_SERVICE_PROTO_DIR) \
		--go_opt=paths=source_relative \
		--go_out=$(FILE_SERVICE_PROTO_MQ_OUT_DIR) \
		$(FILE_SERVICE_PROTO_MQ_FILE)
