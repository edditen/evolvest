FROM golang:1.14.3-alpine as builder
LABEL maintainer="tzz1002@gmail.com"
ENV GO111MODULE=on
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
ADD . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o evolvestd ./cmd/evolvestd

FROM scratch
WORKDIR /app
ADD conf conf
COPY --from=builder /app/evolvestd .
EXPOSE 8762 8080
ENTRYPOINT ["./evolvestd", "-c", "conf/config.yaml"]
