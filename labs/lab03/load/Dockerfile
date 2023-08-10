FROM golang:latest
WORKDIR /app/


COPY ./app  /app/
RUN go get
RUN go build -v 
ENTRYPOINT [ "/app/load"]
