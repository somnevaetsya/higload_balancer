FROM golang:alpine as builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go mod download
RUN go build -o main .

FROM alpine
WORKDIR /app
EXPOSE 80
EXPOSE 9090
COPY --from=builder /app/main /app/main
CMD ["./main"]
