.PHONY: up
up:
	docker-compose up -d
	go run main.go

.PHONY: slqc
sqlc:
	go install github.com/kyleconroy/sqlc/cmd/sqlc@latest
	sqlc generate

.PHONY: k6
k6:
	k6 run tools/k6/transfer.js