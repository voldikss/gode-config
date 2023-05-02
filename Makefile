example:
	GODE_CONFIG_DIR=${PWD}/examples/yaml/config go run examples/yaml/main.go
	# GODE_CONFIG_DIR=${PWD}/examples/json/config go run examples/json/main.go
	# go run examples/json/main.go

test:
	go test -v ./...
