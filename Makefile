env?=dev

build-go:
	go build

build-ui:
	cd ./views;\
	yarn build

build-all: build-go build-ui

serve: build-all
	./basket-products -env=$(env)

