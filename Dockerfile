FROM golang:1.18.1

WORKDIR /app

COPY . ./
RUN go mod download


RUN CGO_ENABLED=0 GOOS=linux go build -o /diskit-server

EXPOSE 8080

CMD ["/diskit-server"]