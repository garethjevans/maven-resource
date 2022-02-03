.PHONY: build
build:
	go build -o build/maven-resource main.go

test:
	go test ./...

lint:
	golangci-lint run

check-%: build
	cat cmd/testdata/check-$*.json
	cat cmd/testdata/check-$*.json | ./build/maven-resource check

in-%: build
	mkdir -p test-output
	cat cmd/testdata/in-$*.json | ./build/maven-resource in test-output
	tree test-output
	find test-output -type file | grep -v jar | xargs cat -A
	rm -fr test-output

test-all: test check-postgres check-first-attempt check-lifecycle-versions in-postgres
