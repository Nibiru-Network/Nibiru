# Build nbn in a stock Go builder container
FROM golang:1.16-alpine as builder

RUN apk add --no-cache make gcc musl-dev linux-headers git

ADD . /nbn
RUN cd /nbn && make nbn

# Pull nbn into a second stage deploy alpine container
FROM alpine:latest

RUN apk add --no-cache ca-certificates
COPY --from=builder /nbn/build/bin/nbn /usr/local/bin/

EXPOSE 8545 8546 30303 30303/udp
ENTRYPOINT ["nbn"]
