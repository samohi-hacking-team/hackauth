FROM golang:1.14-buster
WORKDIR /app
RUN chown nobody /app
RUN go get -u github.com/gobuffalo/packr/packr
COPY go.mod go.sum ./
RUN go mod download
COPY main.go .
RUN go build -o t main.go

FROM debian:buster-slim
WORKDIR /app
COPY --from=0 /app/ .
USER nobody
EXPOSE 3000
EXPOSE 80
STOPSIGNAL SIGINT
ENTRYPOINT ["./t"]
