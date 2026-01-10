.PHONY: proto clean

PROTO_DIR := api/proto/scale-handler
OUT_DIR := pkg/api/proto/scale-handler

proto:
	@echo "Cleaning old generated files..."
	@rm -rf $(OUT_DIR)/*.go
	@mkdir -p $(OUT_DIR)
	
	@echo "Generating with protoc..."
	protoc -I $(PROTO_DIR) \
		--go_out=$(OUT_DIR) \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(OUT_DIR) \
		--go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/*.proto
	
	@echo "Done! Generated files in $(OUT_DIR)/"
	@ls -la $(OUT_DIR)/

clean:
	rm -rf $(OUT_DIR)/*.go

proto-easyp:
	@echo "Generating easyp code..."
	@which protoc-gen-easyp >/dev/null 2>&1 || go install github.com/monzo/typhon/protoc-gen-easyp@latest
	protoc -I $(PROTO_DIR) \
		--plugin=protoc-gen-easyp=$$(which protoc-gen-easyp) \
		--easyp_out=$(OUT_DIR) \
		--easyp_opt=paths=source_relative \
		$(PROTO_DIR)/*.proto

proto: proto-go proto-easyp