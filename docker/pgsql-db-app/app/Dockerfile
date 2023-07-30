FROM golang:latest
WORKDIR /usr/src/app/
RUN git clone https://github.com/filatkinen/socialnet
WORKDIR /usr/src/app/socialnet
COPY server* /usr/src/app/socialnet/configs/
ENV SOCIALNET_DB_USER=socialnet
ENV SOCIALNET_DB_PASS=socialnet
RUN go build -v -o ./bin/socialnet  ./cmd/server
ENTRYPOINT [ "./bin/socialnet"]
CMD ["-config", "./configs/server.pgsql.yaml"] 
EXPOSE 8800
