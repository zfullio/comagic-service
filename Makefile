gen:
	protoc --go_out=. --go-grpc_out=.  api/grpc/*.proto

build:
	go build -o ./bin/server_app ./cmd/server/main.go
	go build -o ./bin/cli_app ./cmd/cli/main.go
	go build -o ./bin/schedule_app ./cmd/schedule/main.go
