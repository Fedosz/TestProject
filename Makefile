proto:
	protoc \
		--proto_path=api \
		--go_out=. \
		--go-grpc_out=. \
		api/rates.proto