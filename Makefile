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