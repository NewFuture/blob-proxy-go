# build stage
FROM golang:alpine AS builder
RUN apk --no-cache add build-base git gcc ca-certificates
COPY app.go /src/
RUN cd /src && go build -o app

# final stage
FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /src/app .
CMD ["./app"]  
