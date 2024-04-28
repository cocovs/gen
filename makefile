GOPATH:=$(shell go env GOPATH)
APP_NAME="gen"

build:
	go build -o ${APP_NAME} .

buildall:buildlinux buildmac buildwin

buildlinux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${APP_NAME}-linux .

buildmac:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ${APP_NAME}-mac .

buildwin:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ${APP_NAME}-win .