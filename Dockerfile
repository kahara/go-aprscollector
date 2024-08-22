FROM golang:1.22.5-bullseye AS build

RUN mkdir /workdir
COPY go.* /workdir/
COPY *.go /workdir/

WORKDIR /workdir
RUN go build -o aprscollector .

FROM gcr.io/distroless/base-debian12 AS production

COPY --from=build /workdir/aprscollector /

CMD ["/aprscollector"]
