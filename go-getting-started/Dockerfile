FROM golang:alpine

WORKDIR /app

COPY . .

RUN go mod download -x
RUN go build -v 

CMD ["./go-getting-started"]