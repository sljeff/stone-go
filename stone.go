package stone

import (
	"bytes"
	"github.com/alecthomas/participle"
	"io"
	"io/ioutil"
)

func Parse(r io.Reader) (*SectionMap, error) {
	parser, err := participle.Build(&Stone{})
	if err != nil {
		return nil, err
	}
	stone := Stone{}
	if err = parser.Parse(r, &stone); err != nil {
		return nil, err
	}
	return stone.eval(), nil
}

func ParseFile(f string) (*SectionMap, error) {
	b, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}
	return Parse(bytes.NewReader(b))
}
