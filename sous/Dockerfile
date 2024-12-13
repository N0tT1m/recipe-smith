FROM golang:1.18.4-alpine3.16
WORKDIR /app

ENV SRC_DIR=/go/src/github.com/Mardiniii/search_engine_client/
COPY . $SRC_DIR
COPY views /app/views

RUN cd $SRC_DIR; go build -o go_search_engine_client; cp search_engine_client /app/
ENTRYPOINT ["./search_engine_client"]