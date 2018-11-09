FROM ubuntu:18.04

ENV GOROOT "/usr/lib/go-1.10"

ENV PATH "${PATH}:$GOROOT/bin"

ENV PORT 8080

RUN apt-get update && \
    apt-get install -y golang-1.10-go

COPY . "${GOROOT}/src/github.com/domdom82/go-backpressure"

RUN cd "${GOROOT}/src/github.com/domdom82/go-backpressure" && \
    go build

CMD cd "${GOROOT}/src/github.com/domdom82/go-backpressure" && \
    ./go-backpressure -server server/config.yml

EXPOSE 8080