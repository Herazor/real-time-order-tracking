FROM alpine:3.22

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY server .

EXPOSE 8080

CMD ["./server"]