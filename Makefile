env?=dev

all: serve

get-go-dep:
	go get

build-go:
	go build

build-ui:
	cd ./views;\
	yarn;yarn build;\

build-all: get-go-dep build-go build-ui

serve: build-all
	./basket-products -env=$(env)

