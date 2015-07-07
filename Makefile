build:
	go build

release:
	gox -osarch="darwin/amd64 linux/amd64" -output="./bin/howdy_{{.OS}}_{{.Arch}}"