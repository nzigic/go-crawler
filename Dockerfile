FROM golang:1.7.3 AS builder
WORKDIR /Users/nemanjazigic/go/src/crawler
RUN go get -d -v golang.org/x/net/html  
COPY app.go    .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /Users/nemanjazigic/go/src/crawler/main .
CMD ["./app"]  