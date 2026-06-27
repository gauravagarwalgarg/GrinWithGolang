# Makefile for gRPC service
#
# Prerequisites:
#   brew install protobuf
#   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
#   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

.PHONY: proto server client clean

# Generate Go code from .proto files
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/user.proto

# Run the server
server:
	go run cmd/server/main.go

# Run the client
client:
	go run cmd/client/main.go

clean:
	rm -f pb/*.go
