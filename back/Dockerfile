FROM golang:1.21.1-alpine3.18 AS languago_back

ENV LANGUAGO_CONFIG_DIR=/app/cfg/

RUN mkdir /app
WORKDIR /app
COPY . .

RUN go mod tidy && \
    go build -o build/languago_main -buildmode=default && \
    chmod +x ./build/languago_main

CMD [ "./build/languago_main" ]
