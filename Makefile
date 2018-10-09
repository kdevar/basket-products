env?=dev

all: serve

build-go:
	go build

build-ui:
	cd ./views;\
	yarn;yarn build;\

build-all: build-go build-ui

serve: build-all
	./basket-products -env=$(env)

