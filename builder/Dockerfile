FROM golang as build

WORKDIR /go/src/builder
COPY . /go/src/builder/
WORKDIR /go/src/builder/cmd/builder
RUN go build

FROM ubuntu as runner
RUN apt-get update && apt-get install -y ca-certificates openssl
ARG cert_location=/usr/local/share/ca-certificates
COPY --from=build  /go/src/builder/cmd/builder/builder /run/
WORKDIR run
#CMD ["./hook"]