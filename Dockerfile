FROM golang:1.23

WORKDIR /app

COPY . .

RUN make build

CMD ["./bin/aurumServer-linux-amd64"]
