FROM golang:1.18

RUN go version
ENV GOPATH=/
ARG GITLAB_USER
ARG GITLAB_TOKEN
ENV GOPRIVATE=

COPY ./ ./

# install psql
RUN apt-get update
RUN apt-get -y install postgresql-client

# set credential to privat repo
RUN echo "machine  login ${GITLAB_USER} password ${GITLAB_TOKEN}" > ~/.netrc

# build go app
RUN go mod download
RUN go build -o rpg-api-telegram ./cmd/tgbotapi/main.go

EXPOSE 5012

CMD ["./rpg-api-telegram"]
