FROM golang:1.12.9-alpine

ADD ./src /go/src/fly
WORKDIR /go/src/fly
ENV PORT=3000
CMD ["go", "run", "main.go"]