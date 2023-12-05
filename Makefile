before.build:
	go mod download

build.testnet:
	CGO_ENABLED=0 GOOS=linux go build testnet.go

run:
	./testnet
