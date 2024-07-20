COMMIT_HASH := $(shell git rev-parse HEAD)

build:
	go build -ldflags "-X 'main.commitHash=$(COMMIT_HASH)'"

watch:
	CompileDaemon -build "go build"

get-deps:
	go get "gopkg.in/yaml.v2"

install: build
	sudo mv ./brutedrop /sbin/brutedrop && chmod 0755 /sbin/brutedrop

