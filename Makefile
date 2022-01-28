.PHONY: build
build:
	go build -o build/maven-resource main.go

test:
	go test ./...

lint:
	golangci-lint run

int: build
	./build/maven-resource in --groupId org.postgresql --artifactId postgresql --type jar --repository https://repo1.maven.org/maven2
