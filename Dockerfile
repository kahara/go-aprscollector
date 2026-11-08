FROM golang:1.26.1-trixie AS build

RUN mkdir /workdir
COPY go.* /workdir/
COPY *.go /workdir/

WORKDIR /workdir
RUN go build -o aprscollector .

FROM gcr.io/distroless/base-debian13 AS production

COPY --from=build /workdir/aprscollector /

CMD ["/aprscollector"]
