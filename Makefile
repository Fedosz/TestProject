APP_NAME=app

build:
	go build -o $(APP_NAME) ./cmd/app

run:
	go run ./cmd/app

test:
	go test ./...

lint:
	golangci-lint run ./...

docker-build:
	docker build -t rates_project:latest .

proto:
	protoc \
		--proto_path=api \
		--go_out=. \
		--go-grpc_out=. \
		api/rates.proto