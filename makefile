
# 生成 Swagger 文檔
swagger:
	swag init -g cmd/servers/main.go -o docs/swagger

goose_up:
	goose -dir migrations postgres "host=localhost user=postgres password=root dbname=ec_wallet port=5432 sslmode=disable TimeZone=UTC" up

goose_down:
	goose -dir migrations postgres "host=localhost user=postgres password=root dbname=ec_wallet port=5432 sslmode=disable TimeZone=UTC" down


# 初始化address pool
init_pool:
	go run cmd/cli/main.go init-pool --chain=ETH --count=200
