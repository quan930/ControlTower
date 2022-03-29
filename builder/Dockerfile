FROM golang:1.17

WORKDIR /go/src/builder
COPY . /go/src/builder/
WORKDIR /go/src/builder/cmd/builder
RUN apt-get update -y && apt-get install -y libgpgme-dev && apt-get install -y libbtrfs-dev
RUN go mod tidy && go build

FROM quay.io/podman/stable:v3.4.4 as runner
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