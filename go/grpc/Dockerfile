FROM golang:alpine as builder

WORKDIR /app 

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install -ldflags="-w -s" ./cmd/greeter_server
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install -ldflags="-w -s" ./cmd/global_server


FROM scratch

WORKDIR /app

COPY --from=builder /go/bin/greeter_server /usr/bin/
COPY --from=builder /go/bin/global_server /usr/bin/

EXPOSE 50051
