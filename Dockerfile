FROM golang:1.22
ENV CGO_ENABLED=0 GOOS=linux GO111MODULE=on
WORKDIR  /app
COPY .  .
COPY go.mod go.sum ./


RUN go mod download
RUN go build -o main .
# WORKDIR /
EXPOSE 9000
ENTRYPOINT ["./main"]