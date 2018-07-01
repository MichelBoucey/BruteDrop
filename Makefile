build:
	go build -o brutedrop src/main.go

watch:
	CompileDaemon -build "go build -o brutedrop src/main.go"	

get-deps:
	go get "gopkg.in/yaml.v2"

install: build
	cp ./brutedrop /usr/bin/brutedrop
