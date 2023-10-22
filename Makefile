run:
	cd back && \
	go mod tidy && \
	LANGUAGO_CONFIG_DIR="./cfg/" go run ./cmd/main.go

build:
	cd back && \
	go mod tidy && \
    go build -o ./build/languago_main -buildmode=default && \
    chmod +x ./build/languago_main

svc.proto:
	protoc --proto_path=./authorizer/internal/proto --go_out=./authorizer/internal/proto/authapi.pb ./authorizer/internal/proto/pb/auth_api.proto
	protoc --proto_path=./back/internal/proto --go_out=./back/internal/pbauth ./back/internal/proto/auth_api.proto
