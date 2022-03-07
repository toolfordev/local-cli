FROM docker.io/library/golang:1.16

WORKDIR /toofordev

COPY . .

RUN go mod vendor; \
    go build .
