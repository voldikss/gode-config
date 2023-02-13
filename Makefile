example-json:
	go run examples/json/main.go

example-yaml:
	go run examples/yaml/main.go

clean:
	rm -rf ./gitlab-merge-bot || echo

test:
	go test -v ./...
