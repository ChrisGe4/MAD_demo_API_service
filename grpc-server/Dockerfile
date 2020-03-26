FROM golang:stretch AS build-env
ADD . /go-src
WORKDIR /go-src
RUN go build -o /grpc-server cmd/main.go

FROM gcr.io/distroless/base
COPY --from=build-env /grpc-server /grpc-server

EXPOSE 8082

ENTRYPOINT ["/grpc-server"]