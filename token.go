package stone

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
)

const (
	regexPat   = `\s*((//.*)|([0-9]+)|("(\\"|\\\\|\\n|[^"])*")|[A-Z_a-z][A-Z_a-z0-9]*|==|<=|>=|&&|\|\||[[:punct:]])?`
	NumLiteral = 1
	Identifier = 2
	StrLiteral = 3
)

var (
	EOL          = Token{value: `\n`}
	EOF          = Token{}
	Parttern     = regexp.MustCompile(regexPat)
	tokenTypeAll = map[TokenType]struct{}{
		NumLiteral: {},
		Identifier: {},
		StrLiteral: {},
	}
	tokenTypeText = map[TokenType]struct{}{
		Identifier: {},
		StrLiteral: {},
	}
	tokenTypeNum = map[TokenType]struct{}{NumLiteral: {}}
)

type TokenType uint8

type Token struct {
	lineNum   uint
	tokenType TokenType

	value string
}

func NewToken(lineNum uint, tokenType TokenType, value string) (*Token, error) {
	if _, ok := tokenTypeAll[tokenType]; !ok {
		return nil, fmt.Errorf("invalid token type: %v", tokenType)
	}
	return &Token{lineNum: lineNum, tokenType: tokenType, value: value}, nil
}

func (token *Token) IsStrLiteral() bool {
	return token.tokenType == StrLiteral
}

func (token *Token) IsIdentifier() bool {
	return token.tokenType == Identifier
}

func (token *Token) IsNumLiteral() bool {
	return token.tokenType == NumLiteral
}

func (token *Token) GetText() (string, error) {
	if _, ok := tokenTypeText[token.tokenType]; !ok {
		return "", fmt.Errorf("token has no text")
	}
	if token.tokenType == Identifier {
		return token.value, nil
	} else if token.tokenType == StrLiteral {
		return toStrLiteral(token.value), nil
	}
	return "", fmt.Errorf("token has no text")
}

func (token *Token) GetNumber() (int, error) {
	if _, ok := tokenTypeNum[token.tokenType]; !ok {
		return 0, fmt.Errorf("token has no number")
	}
	return strconv.Atoi(token.value)
}

func toStrLiteral(text string) string {
	value := []rune(text)
	// exclude first and last quote
	pos := 1
	lenSub1 := len(value) - 1

	realChars := []rune{}
	var c rune
	for pos < lenSub1 {
		c = value[pos]
		if c == '\\' && pos+1 < lenSub1 {
			nextC := value[pos+1]
			if nextC == '"' || nextC == '\\' {
				pos += 1
				c = nextC
			} else if nextC == 'n' {
				pos += 1
				c = '\n'
			}
		}
		realChars = append(realChars, c)
		pos += 1
	}
	return string(realChars)
}

/* Lexer */

type Lexer struct {
	reader io.Reader

	scanner     *bufio.Scanner
	queue       []*Token
	hasMore     bool
	scanLineNum uint
}

func NewLexer(reader io.Reader) *Lexer {
	return &Lexer{
		reader:  reader,
		scanner: bufio.NewScanner(reader),
		hasMore: true,
	}
}

func (lex *Lexer) readLine() error {
	ok := lex.scanner.Scan()
	if !ok {
		lex.hasMore = false
		return nil
	}
	lex.scanLineNum += 1
	lineText := lex.scanner.Text()
	pos := 0
	endPos := len(lineText)
	for pos < endPos {
		subLine := lineText[pos:endPos]
		indexes := Parttern.FindStringSubmatchIndex(subLine)
		if len(indexes) == 0 {
			return fmt.Errorf("invalid line: %s", subLine)
		}
		if err := lex.addToken(subLine, indexes, lex.scanLineNum); err != nil {
			return err
		}
		pos += indexes[1]
	}
	lex.queue = append(lex.queue, &EOL)
	return nil
}

func (lex *Lexer) addToken(text string, indexes []int, lineNum uint) error {
	if indexes[2] == indexes[3] { // total
		return nil
	}
	if indexes[4] != indexes[5] { // comment
		return nil
	}
	var tokenType TokenType
	var tokenStart, tokenEnd int

	if indexes[6] != indexes[7] { // NumLiteral
		tokenType = NumLiteral
		tokenStart = indexes[6]
		tokenEnd = indexes[7]
	} else if indexes[8] != indexes[9] { // StrLiteral
		tokenType = StrLiteral
		tokenStart = indexes[8]
		tokenEnd = indexes[9]
	} else { // Identifier
		tokenType = Identifier
		tokenStart = indexes[2]
		tokenEnd = indexes[3]
	}

	token, err := NewToken(lineNum, tokenType, text[tokenStart:tokenEnd])
	if err != nil {
		return err
	}
	lex.queue = append(lex.queue, token)
	return nil
}
