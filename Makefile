gen:
	protoc --go_out=. --go-grpc_out=.  api/grpc/*.proto

gen_python:
	python -m grpc_tools.protoc -I ./ --python_out=./python --pyi_out=./python --grpc_python_out=./python api/grpc/*.proto

build_server:
	go build -o ./bin/server_app ./cmd/server/main.go