FROM golang:latest
RUN mkdir /app
ADD . /app/
WORKDIR /app

EXPOSE 9000:9000
RUN go build -o main .
CMD ["/app/main"]