FROM golang:1.18

RUN go version
ENV GOPATH=/
ARG GITLAB_USER
ARG GITLAB_TOKEN
ENV GOPRIVATE=git.andersenlab.com

COPY ./ ./

# install psql
RUN apt-get update
RUN apt-get -y install postgresql-client

# set credential to privat repo
RUN echo "machine git.andersenlab.com login ${GITLAB_USER} password ${GITLAB_TOKEN}" > ~/.netrc

# build go app
RUN go mod download
RUN go build -o go-aut-registration-user-telegram ./cmd/tgbotapi/main.go

CMD ["./go-aut-registration-user-telegram"]
