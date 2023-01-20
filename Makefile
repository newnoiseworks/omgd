build:
	go build -o ~/go/bin/omgdd

build-prod:
	go build -o ~/go/bin/omgd

test:
	go test ./... -v
