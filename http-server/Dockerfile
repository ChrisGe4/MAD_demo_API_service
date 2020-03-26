FROM golang:stretch AS build-env
ADD . /go-src
WORKDIR /go-src
RUN go build -o /http-server cmd/main.go

FROM gcr.io/distroless/base
COPY --from=build-env /http-server /http-server

EXPOSE 8082

ENTRYPOINT ["/http-server"]