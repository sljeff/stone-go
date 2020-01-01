package stone

import (
	"strings"
	"testing"
)

func TestToken_GetText(t *testing.T) {
	type fields struct {
		lineNum   uint
		tokenType TokenType
		value     string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name:    "normal text",
			fields:  fields{tokenType: 3, value: `"normal text"`},
			want:    "normal text",
			wantErr: false,
		},
		{
			name:    `normal "text"`,
			fields:  fields{tokenType: 3, value: `"normal \"text\""`},
			want:    `normal "text"`,
			wantErr: false,
		},
		{
			name:    `normal \text\`,
			fields:  fields{tokenType: 3, value: `"normal \\text\\"`},
			want:    `normal \text\`,
			wantErr: false,
		},
		{
			name:    "normal \ntext",
			fields:  fields{tokenType: 3, value: `"normal \ntext"`},
			want:    "normal \ntext",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := &Token{
				lineNum:   tt.fields.lineNum,
				tokenType: tt.fields.tokenType,
				value:     tt.fields.value,
			}
			got, err := token.GetText()
			if (err != nil) != tt.wantErr {
				t.Errorf("Token.GetText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Token.GetText() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLexer_readLine(t *testing.T) {
	lex := NewLexer(strings.NewReader(
`// comment
"hello\n"
import this`,
	))
	t.Log("readLine 1")
	err := lex.readLine()
	if err != nil {
		t.Error(err)
	}
	if !lex.hasMore {
		t.Error("lex.hasMore should be true")
	}
	if len(lex.queue) != 1 {  // EOL
		t.Errorf("error length: %v %v", len(lex.queue), lex.queue[0])
	}
	t.Log("readLine 2")
	lex.readLine()
	if len(lex.queue) != 3 {  // 2 EOLs
		t.Errorf("error length: %v", len(lex.queue))
	}
	t.Log(lex.queue[1].GetText())
	t.Log(lex.queue[1].lineNum)
	t.Log("readLine 3")
	// TODO
}
