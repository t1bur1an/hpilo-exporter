FROM golang:1.12.5
RUN go get github.com/golang/dep/cmd/dep
ENV WORKDIRECTORY /go/src/github.com/t1bur1an/hpilo-exporter
WORKDIR ${WORKDIRECTORY}
ADD Gopkg.* ./
RUN dep ensure --vendor-only
ADD . .
RUN go build
CMD ${WORKDIRECTORY}/hpilo-exporter