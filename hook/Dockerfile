FROM golang as build
WORKDIR /go/src/hook
COPY . /go/src/hook/
WORKDIR /go/src/hook/cmd/hook
RUN go build

FROM ubuntu as runner
COPY --from=build  /go/src/hook/cmd/hook/hook /run/
WORKDIR run
CMD ["./hook"]