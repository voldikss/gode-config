example:
	# GO_CONFIG_DIR=${PWD}/examples/yaml/config go run examples/yaml/main.go
	GO_CONFIG_DIR=${PWD}/examples/json/config go run examples/json/main.go
	# go run examples/json/main.go

clean:
	rm -rf ./gitlab-merge-bot || echo

test:
	go test -v ./...
