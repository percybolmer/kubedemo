FROM golang:alpine as builder

WORKDIR /app
COPY . .

RUN CGO_ENBALED=0 GOOS=linux GOARCH=amd64 go build -o hellogopher -ldflags="-w -s"

ENTRYPOINT [ "./hellogopher" ]