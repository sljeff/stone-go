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
import this

myvar = 10086  // comment
anthor = "the way to go"`))
	t.Log("read line 1")
	err := lex.readLine()
	if err != nil {
		t.Error(err)
	}
	if !lex.hasMore {
		t.Error("lex.hasMore should be true")
	}
	if len(lex.queue) != 1 { // EOL
		t.Errorf("error length: %v %v", len(lex.queue), lex.queue[0])
	}

	t.Log("read line 2")
	lex.readLine()
	if len(lex.queue) != 3 { // 2 EOLs
		t.Errorf("error length: %v", len(lex.queue))
	}
	t.Log(lex.queue[1].GetText())
	t.Log(lex.queue[1].lineNum)

	t.Log("read line 3")
	lex.readLine()
	queue := []Token{
		EOL,
		Token{tokenType: StrLiteral, value: `"hello\n"`, lineNum: 2},
		EOL,
		Token{tokenType: Identifier, value: "import", lineNum: 3},
		Token{tokenType: Identifier, value: "this", lineNum: 3},
		EOL,
	}

	t.Log("read line 4")
	lex.readLine()
	queue = append(queue, EOL)

	t.Log("read line 5")
	lex.readLine()
	queue = append(
		queue,
		Token{tokenType: Identifier, value: "myvar", lineNum: 5},
		Token{tokenType: Identifier, value: "=", lineNum: 5},
		Token{tokenType: NumLiteral, value: "10086", lineNum: 5},
		EOL,
	)

	t.Log("read line 6")
	lex.readLine()
	queue = append(
		queue,
		Token{tokenType: Identifier, value: "anthor", lineNum: 6},
		Token{tokenType: Identifier, value: "=", lineNum: 6},
		Token{tokenType: StrLiteral, value: `"the way to go"`, lineNum: 6},
		EOL,
	)

	for index, token := range lex.queue {
		if *token != queue[index] {
			text, err := token.GetText()
			t.Error("error", token.tokenType, text, err)
			text, err = queue[index].GetText()
			t.Error("error", queue[index].tokenType, text, err)
		}
	}
}

func TestLexer_Read(t *testing.T) {
	lex := NewLexer(strings.NewReader(
		`// comment
"hello\n"
import this

myvar = 10086  // comment
anthor = "the way to go"`))

	var token *Token
	token = lex.Read()
	if *token != EOL { // EOL
		t.Error("error")
	}
	token = lex.Read()
	if *token != (Token{tokenType: StrLiteral, value: `"hello\n"`, lineNum: 2}) {
		t.Error("error", token.value, token.tokenType, token.lineNum)
	}
	token = lex.Read()
	if *token != EOL { // EOL
		t.Error("error", token.value, token.tokenType, token.lineNum)
	}
	token = lex.Read()
	if *token != (Token{tokenType: Identifier, value: "import", lineNum: 3}) {
		t.Error("error", token.value, token.tokenType, token.lineNum)
	}
	token = lex.Read()
	if *token != (Token{tokenType: Identifier, value: "this", lineNum: 3}) {
		t.Error("error", token.value, token.tokenType, token.lineNum)
	}
	token = lex.Read()
	if *token != EOL { // EOL
		t.Error("error", token.value, token.tokenType, token.lineNum)
	}
	token = lex.Read()
	if *token != EOL { // EOL
		t.Error("error", token.value, token.tokenType, token.lineNum)
	}
	token = lex.Read()
	if *token != (Token{tokenType: Identifier, value: "myvar", lineNum: 5}) {
		t.Error("error", token.value, token.tokenType, token.lineNum)
	}
	token = lex.Read()
	if *token != (Token{tokenType: Identifier, value: "=", lineNum: 5}) {
		t.Error("error", token.value, token.tokenType, token.lineNum)
	}
	token = lex.Read()
	if *token != (Token{tokenType: NumLiteral, value: "10086", lineNum: 5}) {
		t.Error("error", token.value, token.tokenType, token.lineNum)
	}
	token = lex.Read()
	if *token != EOL { // EOL
		t.Error("error", token.value, token.tokenType, token.lineNum)
	}
	token = lex.Read()
	if *token != (Token{tokenType: Identifier, value: "anthor", lineNum: 6}) {
		t.Error("error", token.value, token.tokenType, token.lineNum)
	}
	token = lex.Read()
	if *token != (Token{tokenType: Identifier, value: "=", lineNum: 6}) {
		t.Error("error", token.value, token.tokenType, token.lineNum)
	}
	token = lex.Read()
	if *token != (Token{tokenType: StrLiteral, value: `"the way to go"`, lineNum: 6}) {
		t.Error("error", token.value, token.tokenType, token.lineNum)
	}
	token = lex.Read()
	if *token != EOL { // EOL
		t.Error("error", token.value, token.tokenType, token.lineNum)
	}
}
