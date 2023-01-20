build:
	go build -o ${GOPATH}/bin/omgdd

build-prod:
	go build -o ${GOPATH}/bin/omgd

test:
	go test ./... -v
