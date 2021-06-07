FROM golang:alpine

WORKDIR /app

ADD go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o appbuild

ENTRYPOINT ["/app/appbuild","run"]
