build:
	go build -o ${GOROOT}/bin/omgdd

build-prod:
	go build -o ${GOROOT}/bin/omgd

test:
	go test ./... -v
