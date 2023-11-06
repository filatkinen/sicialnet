FROM golang:latest
WORKDIR /usr/src/app/

RUN  git  clone https://github.com/filatkinen/socialnet
WORKDIR /usr/src/app/socialnet
COPY dialog.yaml /usr/src/app/socialnet/configs/
RUN go build -v -o ./bin/dialog  ./cmd/dialog

# FROM scratch
# COPY --from=build /usr/src/app/socialnet/bin/socialnet /usr/src/app/socialnet/bin/socialnet
#WORKDIR /usr/src/app/socialnet

ENTRYPOINT [ "./bin/dialog"]
# CMD ["-config", "./configs/server.pgsql.yaml"] 
EXPOSE 8802
