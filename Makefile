gen:
	protoc --go_out=. --go-grpc_out=.  api/grpc/*.proto

gen_python:
	python -m grpc_tools.protoc -I ./ --python_out=./python --pyi_out=./python --grpc_python_out=./python api/grpc/*.proto

build_all:
	go build -o ./bin/server_app ./cmd/server/main.go
	go build -o ./bin/cli_app ./cmd/cli/main.go
	go build -o ./bin/schedule_app ./cmd/schedule/main.go

build_server:
	go build -o ./bin/server_app ./cmd/server/main.go