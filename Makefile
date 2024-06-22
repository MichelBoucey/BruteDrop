build:
	go build -o brutedrop main.go

watch:
	CompileDaemon -build "go build -o brutedrop main.go"

get-deps:
	go get "gopkg.in/yaml.v2"

install: build
	mv ./brutedrop /sbin/brutedrop && chmod 0700 /sbin/brutedrop

test: install
	brutedrop -version | grep -q  License
