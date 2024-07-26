SHORT_SHA := $(shell git rev-parse --short HEAD)

build:
	go build -ldflags "-X 'main.commitHash=$(SHORT_SHA)'"

watch:
	CompileDaemon -build "go build"

get-deps:
	go get "gopkg.in/yaml.v2"

install: build
	sudo mv ./brutedrop /sbin/brutedrop && chmod 0755 /sbin/brutedrop

clean:
	rm -f ./brutedrop
