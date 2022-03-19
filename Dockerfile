FROM golang:1.16-alpine

COPY . /go/src/github.com/bonczj/web-pub-sub
WORKDIR /go/src/github.com/bonczj/web-pub-sub
RUN go install
CMD [ "/go/bin/web-pub-sub" ]