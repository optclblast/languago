run:
	cd back && \
	go mod tidy && \
	LANGUAGO_CONFIG_DIR="./cfg/" go run ./cmd/main.go

build:
	cd back && \
	go mod tidy && \
    cd ./cmd && go build -o ../build/languago_main -buildmode=default && \
    chmod +x ./build/languago_main

netup:
	sudo docker network create --driver bridge --subnet=192.168.18.1/24 --attachable caddy;
test:
	// TODO
svc.proto:
	protoc --proto_path=./authorizer/internal/proto --go_out=./authorizer/internal/proto/authapi.pb ./authorizer/internal/proto/pb/auth_api.proto
	protoc --proto_path=./back/internal/proto --go_out=./back/internal/pbauth ./back/internal/proto/auth_api.proto
