main:
	go build -o alfa-node cmd/alfa/main.go 
	go build -o client-node cmd/node/main.go
	go build -o key-generator cmd/key-generator/main.go

blockchain:
	go build -o alfa-node cmd/alfa/main.go 
	go build -o client-node cmd/node/main.go

key-genrator:
	go build -o key-generator cmd/key-generator/main.go

clean:
	rm alfa-node client-node key-generator