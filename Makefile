BIN := "./bin/socialnet"

export SOCIALNET_DB_USER=socialnet
export SOCIALNET_DB_PASS=socialnet
export SOCIALNET_DB=snet
export SOCIALNET_DB_ADDRESS=localhost



build:
	go build -v -o $(BIN)  ./cmd/server

run-mysql: build
	SOCIALNET_DB_PORT=3306 $(BIN) -config ./configs/server.mysql.yaml

run-pgsql: build
	SOCIALNET_DB_PORT=5432 $(BIN) -config ./configs/server.pgsql.yaml


test-mysql:
	SOCIALNET_DB_PORT=3306 go test -v -race  ./... -tags mysql

test-pgsql:
	SOCIALNET_DB_PORT=5432 go test -v -race  ./... -tags pgsql