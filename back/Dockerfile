FROM fedora:38 AS languago_builder 

RUN dnf update && dnf upgrade && \
    curl -OL https://golang.org/dl/go1.19.5.linux-amd64.tar.gz && \
                                    sudo rm -rf /usr/local/go && \
                                    sudo tar -C /usr/local -xzf go*.linux-amd64.tar.gz && \
    dnf install gcc 

ENV LANGUAGO_CONFIG_DIR=/app/cfg/

RUN mkdir /app
WORKDIR /app
COPY . .

RUN go mod tidy && \
    go build -o build/languago_main -buildmode=default && \
    chmod +x ./build/languago_main

CMD [ "./build/languago_main" ]
