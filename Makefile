build:
	go build -o ${GOPATH}/bin/omgdd

build-prod:
	go build -o ${GOPATH}/bin/omgd

test:
	go test -v ./... | sed ''/PASS/s//$$(printf "\033[32mPASS\033[0m")/'' | sed ''/FAIL/s//$$(printf "\033[31mFAIL\033[0m")/''

todo:
	egrep TODO **/*.go