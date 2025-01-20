FROM golang:1.23

WORKDIR /app

COPY . .

RUN make build_me

CMD ["./bin/vexalServer"]
