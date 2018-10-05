FROM golang:alpine

WORKDIR /go/src/github.com/kdevar

ADD . basket-products

CMD ["go", "run", "./basket-products/main.go"]