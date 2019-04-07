# build binary
FROM golang:1.12-alpine3.9 AS build
ENV GO111MODULE=on
RUN apk add git
WORKDIR /go/mod/github.com/pocoz/drone-tg
COPY . /go/mod/github.com/pocoz/drone-tg
RUN go mod download
RUN CGO_ENABLED=0 go build -o /out/drone-tg github.com/pocoz/drone-tg/cmd/drone-tg-d

# copy to alpine image
FROM alpine:3.9 AS prod
WORKDIR /app
COPY --from=build /out/drone-tg /app
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
CMD ["/app/drone-tg"]
