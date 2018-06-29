build:
	go build -o brutedrop src/brutedrop.go

watch:
	CompileDaemon -build "go build -o brutedrop src/brutedrop.go"	

get-deps:
	go get "gopkg.in/yaml.v2"
