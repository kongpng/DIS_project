FROM golang:1.22.3-alpine3.20

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /myapp

EXPOSE 8080

CMD ["/myapp"]
