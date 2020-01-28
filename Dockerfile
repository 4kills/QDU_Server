FROM golang:1.13 AS builder
WORKDIR $GOPATH/src/github.com/4kills/qdu_server
COPY . . 
RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /main .

FROM alpine:3.11
RUN apk --no-cache add ca-certificates
COPY --from=builder /main ./
COPY /gallery.html ./
RUN chmod +x ./main
ENTRYPOINT [ "./main" ]
EXPOSE 8080
EXPOSE 1337