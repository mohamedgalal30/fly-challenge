FROM golang:1.12.9-alpine
RUN apk add git
LABEL maintainer="Mohamed Galal <mohamed.galal30@gmail.com>"

ADD . /go/src/fly
WORKDIR /go/src/fly
ENV PORT=3000
# RUN go get
RUN go get
CMD ["go", "run", "main.go"]