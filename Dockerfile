FROM golang:1.21.0

WORKDIR /app

COPY . .

ENV GO111MODULE=on

COPY go.mod .
COPY go.sum .

RUN go mod download

RUN go build -o app cmd/main.go

CMD ["/app/app"]