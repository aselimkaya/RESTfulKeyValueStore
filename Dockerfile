FROM golang:latest

RUN mkdir /build
WORKDIR /build

RUN export GO111MODULE=on
RUN go get github.com/aselimkaya/RESTfulKeyValueStore/src
RUN cd /build && git clone https://github.com/aselimkaya/RESTfulKeyValueStore.git

RUN cd /build/RESTfulKeyValueStore/src && go build

EXPOSE 80

ENTRYPOINT [ "/build/RESTfulKeyValueStore/src/src" ]