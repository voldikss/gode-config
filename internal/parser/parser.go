package parser

type Parser interface {
	Parse(filename string, content []byte) (interface{}, error)
}

func ParseConfigFile(filepath string) (any, error) {
	return "", nil
}

func Parse() {

}
