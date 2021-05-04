FROM golang:1.14 as builder

WORKDIR /go/src/pj
ADD . /go/src/pj

ENV GO111MODULE=on

RUN go build -o /go/bin/pj

FROM gcr.io/distroless/static-debian10
COPY --from=builder /go/bin/pj /usr/local/bin/
CMD ["/usr/local/bin/pj"]
