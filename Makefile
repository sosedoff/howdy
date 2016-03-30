build:
	go build

setup:
	go get

release:
	gox -cgo -osarch="darwin/amd64 linux/amd64" -output="./bin/howdy_{{.OS}}_{{.Arch}}"