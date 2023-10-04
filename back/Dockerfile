FROM fedora:38 AS languago_builder 

RUN curl -OL https://golang.org/dl/go1.21.1.linux-amd64.tar.gz && \
                                    sudo rm -rf /usr/local/go && \
                                    sudo tar -C /usr/local -xzf go*.linux-amd64.tar.gz

ENV LANGUAGO_CONFIG_DIR=/app/cfg/

RUN mkdir /app
WORKDIR /app
COPY . .

RUN /usr/local/go/bin/go mod tidy && \
    /usr/local/go/bin/go build -o build/languago_main -buildmode=default && \
    chmod +x ./build/languago_main

CMD [ "./build/languago_main" ]
