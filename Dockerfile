FROM golang:1.9-alpine3.7

LABEL maintainer="Austin J. Alexander <aja@design-brain.com>"

RUN apk add --update --no-cache build-base git

COPY . /go/src/github.com/design-brain/images-api
WORKDIR /go/src/github.com/design-brain/images-api

RUN make

EXPOSE 8080

CMD ["make", "run"]
