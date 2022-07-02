FROM golang:alpine as builder

#Install git
RUN apk add --no-cache git

WORKDIR /app 

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go install -ldflags="-w -s" ./cmd/auth_server
RUN CGO_ENABLED=0 GOOS=linux go install -ldflags="-w -s" ./cmd/dashboard_server
RUN CGO_ENABLED=0 GOOS=linux go install -ldflags="-w -s" ./cmd/toggles_server

FROM scratch

WORKDIR /app

COPY --from=builder /go/bin/auth_server /usr/bin/
COPY --from=builder /go/bin/dashboard_server /usr/bin/
COPY --from=builder /go/bin/toggles_server /usr/bin/