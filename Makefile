run: build-compiler  build-container

build-compiler: ./main/compiler.go
	go build -o compiler ./main/compiler.go

build-container: ./main/container.go
	go build -o container ./main/container.go


