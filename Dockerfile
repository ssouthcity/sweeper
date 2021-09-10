FROM golang:latest

WORKDIR /sweeper
COPY . .

RUN go get -d -v ./...
RUN go install ./cmd/...

CMD [ "sweeperbot" ]