BIN := "./bin/socialnet"

export SOCIALNET_DB_USER=socialnet
export SOCIALNET_DB_PASS=socialnet


build:
	go build -v -o $(BIN)  ./cmd/server

run-mysql: build
	$(BIN) -config ./configs/server.mysql.yaml

run-pgsql: build
	$(BIN) -config ./configs/server.pgsql.yaml


test-mysql:
	go test -v -race  ./... -tags mysql

test-pgsql:
	go test -v -race  ./... -tags pgsql