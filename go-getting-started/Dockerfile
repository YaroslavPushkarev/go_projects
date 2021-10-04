FROM golang:1.16

WORKDIR /app

COPY . .

RUN go get -d -v
RUN go build -v 

CMD ["./go-getting-started"]