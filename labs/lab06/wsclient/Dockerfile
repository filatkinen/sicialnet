FROM golang:1.20-alpine
WORKDIR /usr/src/app/

RUN  git  clone https://github.com/filatkinen/socialnet
WORKDIR /usr/src/app/socialnet
RUN go build -v -o ./bin/wsclient  ./cmd/ws-client

ENTRYPOINT ["tail"]
CMD ["-f","/dev/null"]
