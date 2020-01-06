package stone

import (
	"github.com/alecthomas/participle"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

var CODE = `
[This]
a = 1
b = 10.1
c = TRUE
d = "this is good"
e = ${SENTRY_DSN}
DELETE a

f = [1, 2, 3]
ff = []
g = [1, 10.1, c, FALSE, ${SENTRY_DSN}]
h = {
	"a": 1,
	"b": FALSE,
	"c": 10.1,
	"d": ${SENTRY_DSN},
	"e": "this is good",
	"f": g,
	"g": [1, 2, FALSE],
}
hh = {}
i = [
	[1, 2, 3],
	hh,
	{"a": TRUE},
	"string"
]
DELETE h["d"]
DELETE i[3]
h["d"] = 100
i[0] = "string"
i[1] = h["f"]

[That] < [This]
i[0] = [1, 2, 3]
`

func TestGrammar(t *testing.T) {
	parser, err := participle.Build(&Stone{})
	if err != nil {
		t.Fatal(err)
	}
	stone := Stone{}
	err = parser.Parse(strings.NewReader(CODE), &stone)
	if err != nil {
		t.Fatal(err)
	}

	assert.Len(t, stone.Sections, 2)
	assert.Len(t, stone.Sections[0].Stmts, 17)
	assert.Len(t, stone.Sections[1].Stmts, 1)
	assert.Equal(t, stone.Sections[0].Name, "This")
	assert.Equal(t, stone.Sections[0].ParentName, (*string)(nil))
	assert.Equal(t, stone.Sections[1].Name, "That")
	assert.Equal(t, *stone.Sections[1].ParentName, "This")
	assert.Equal(t, stone.Sections[0].Stmts[0].Assignment.Left.Identifier, "a")
	assert.Equal(t, *stone.Sections[0].Stmts[0].Assignment.Right.Int, 1)
}
