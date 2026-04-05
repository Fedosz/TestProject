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

mocks:
	minimock -i rates_project/internal/rates.ExchangeClient,rates_project/internal/rates.RateRepository -o ./internal/rates/mocks -s "_mock.go"
	minimock -i rates_project/internal/app/health.Pinger -o ./internal/app/health/mocks -s "_mock.go"