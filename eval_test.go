package stone

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestEval(t *testing.T) {
	sm, err := Parse(strings.NewReader(CODE))
	if err != nil {
		t.Fatal(err)
	}
	this := (*sm)["This"]
	that := (*sm)["That"]
	assert.Equal(t, (*this)["b"], 10.1)
	assert.Equal(t, (*that)["b"], 10.1)
	h := (*this)["h"].(map[string]interface{})
	g := h["g"].([]interface{})
	assert.Equal(t, g[0], 1)
	h = (*that)["h"].(map[string]interface{})
	g = h["g"].([]interface{})
	assert.Equal(t, g[0], 1)

	i := (*this)["i"].([]interface{})
	assert.Equal(t, i[0], "string")
	i = (*that)["i"].([]interface{})
	assert.Equal(t, i[0], []interface{}{1, 2, 3})
	_, ok := (*this)["a"]
	assert.False(t, ok)
}
