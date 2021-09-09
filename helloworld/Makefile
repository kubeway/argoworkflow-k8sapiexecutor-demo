all: helloworld

helloworld:
	@echo "build helloworld"
	go build -o helloworld helloworld.go

image:
	@echo "build helloworld image"
	docker build -t "helloworld:latest" -f ./Dockerfile ./