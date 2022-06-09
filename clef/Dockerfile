FROM golang:1.15-alpine as build

ARG TAG="v1.9.24"

RUN apk add --no-cache make gcc musl-dev linux-headers git

RUN git clone https://github.com/ethereum/go-ethereum.git && cd /go/go-ethereum && \
    if [[ -n $TAG ]]; then git checkout $TAG; fi && env GO111MODULE=on go run build/ci.go install ./cmd/clef

FROM alpine:latest as runtime

RUN apk add --no-cache ca-certificates
COPY --from=build /go/go-ethereum/build/bin/clef /usr/local/bin/clef
RUN touch /test

COPY rules.js /app/config/rules.js
COPY 4byte.json /app/config/4byte.json
COPY password /app/config/password
COPY entrypoint.sh /entrypoint.sh

RUN /entrypoint.sh init

EXPOSE 8550

ENTRYPOINT ["/entrypoint.sh"]