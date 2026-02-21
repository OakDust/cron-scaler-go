.PHONY: proto migrate-up migrate-down migrate-status migrate-create clean

# Генерация proto
proto:
	@echo "Generating Go code from proto files..."
	@mkdir -p pkg/api/proto/scale-handler
	protoc -I api/proto/scale-handler \
		--go_out=pkg/api/proto/scale-handler \
		--go_opt=paths=source_relative \
		--go-grpc_out=pkg/api/proto/scale-handler \
		--go-grpc_opt=paths=source_relative \
		api/proto/scale-handler/*.proto
	@echo "Code generation completed"

# Миграции
DB_STRING?=postgres://postgres:postgres@localhost:5432/scale_handler?sslmode=disable

migrate-up:
	goose -dir migrations postgres "$(DB_STRING)" up

migrate-down:
	goose -dir migrations postgres "$(DB_STRING)" down

migrate-status:
	goose -dir migrations postgres "$(DB_STRING)" status

migrate-create:
	@read -p "Enter migration name: " name; \
	goose -dir migrations create $$name sql

clean:
	rm -rf pkg/api/proto/scale-handler/*.go

# Development
dev:
	go run cmd/main.go