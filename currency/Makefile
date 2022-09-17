.PHONY: protos

protos:
	 rm -f pb/*.go
	 protoc --proto_path=protos  --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
    protos/*.proto