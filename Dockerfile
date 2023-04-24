FROM golang:alpine AS builder
WORKDIR /build
COPY . .
RUN go build -o comagic-service /build/cmd/server/main.go

FROM alpine
WORKDIR /app
COPY --from=builder /build/comagic-service ./
RUN chmod +x comagic-service
RUN apk add --no-cache tzdata
ENV TZ=Europe/Moscow
ENTRYPOINT ["./comagic-service", "--env"]