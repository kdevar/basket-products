# STEP 1 build executable binary
FROM golang:alpine as builder
COPY . $GOPATH/src/github.com/kdevar/basket-products
WORKDIR $GOPATH/src/github.com/kdevar/basket-products

RUN apk add --no-cache git make

RUN go get -d -v

RUN GOARCH=amd64 CGO_ENABLED=0 GOOS=linux go build -o /go/bin/basket-products

FROM scratch
# Copy our static executable
COPY --from=builder /go/bin/basket-products /go/bin/basket-products
ENTRYPOINT ["/go/bin/basket-products"]