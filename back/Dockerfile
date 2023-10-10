FROM fedora:38 AS languago_builder 

RUN curl -OL https://golang.org/dl/go1.21.1.linux-amd64.tar.gz && \
                                    sudo rm -rf /usr/local/go && \
                                    sudo tar -C /usr/local -xzf go*.linux-amd64.tar.gz

ENV LANGUAGO_CONFIG_DIR=/app/cfg/

RUN mkdir /app
RUN mkdir /app/build
WORKDIR /app
COPY . .

RUN /usr/local/go/bin/go mod tidy 
RUN cd cmd && /usr/local/go/bin/go build -o /app/build/languago_main -buildmode=default
RUN chmod +x ./build/languago_main

EXPOSE 3300

CMD [ "./build/languago_main" ]
