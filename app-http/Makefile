IMG ?= guwe/apiserver-load-tester-http:$(shell git rev-parse --short=7 HEAD)

docker-build:
	docker build -t ${IMG} . -f Dockerfile