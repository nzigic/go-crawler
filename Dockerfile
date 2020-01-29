FROM golang:alpine
RUN apk update && apk add --no-cache git
WORKDIR /go/src/
COPY . .
RUN go get -d -v
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/crawler

FROM scratch  
COPY --from=0 /go/bin/crawler /crawler

# Run the binary.
ENTRYPOINT ["/crawler"]