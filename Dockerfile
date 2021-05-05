FROM golang:1.16-alpine3.13
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go mod tidy
RUN go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /app .
CMD ["/root/main"]