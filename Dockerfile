FROM golang:latest AS build

WORKDIR /go/src/github.com/iofq/ip
COPY . .

ENV GO111MODULE=on CGO_ENABLED=0 GOOS=linux
RUN make

FROM scratch
EXPOSE 8080

COPY --from=build /go/bin/ip /opt/ip
ENTRYPOINT ["/opt/ip"]
