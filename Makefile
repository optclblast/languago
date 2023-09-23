run:
	cd back && \
	go mod tidy && \
	LANGUAGO_CONFIG_DIR="./cfg/" go run main.go

build:
	cd back && \
	go mod tidy && \
    go build -o ./build/languago_main -buildmode=default && \
    chmod +x ./build/languago_main

test:
	// TODO