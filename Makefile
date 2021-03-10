test: download
	go test -v -race ./...

cover:
	go test -v -coverprofile=coverage.txt -covermode=atomic ./...

dep:
	go get -u ./...

download:
	go mod download
