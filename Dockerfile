FROM golang:1.18

RUN go version
ENV GOPATH=/

COPY ./ ./

# install psql
RUN apt-get update
RUN apt-get -y install postgresql-client

# build go app
RUN go mod download
RUN go build -o go-aut-registration-user-telegram ./cmd/tgbotapi/main.go

CMD ["./rpg-api-telegram"]