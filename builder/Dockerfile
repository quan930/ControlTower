FROM golang as build

WORKDIR /go/src/builder
COPY . /go/src/builder/
WORKDIR /go/src/builder/cmd/builder
RUN go build

FROM docker:20.10.12-git as runner
RUN apk add --no-cache netcat-openbsd
COPY --from=build  /go/src/builder/cmd/builder/builder /run/
WORKDIR run
ENV REPO = ""
ENV BRANCH = ""
ENV DOCKERFILE = ""
ENV IMAGE = ""
ENV USER = ""
ENV PASSWORD = ""
CMD ["./builder"]