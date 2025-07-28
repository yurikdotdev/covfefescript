package lexer

import (
	"testing"

	"github.com/yurikdotdev/covfefescript/internal/token"
)

func TestSimpleTokens(t *testing.T) {
	input := `+-%*/<>!(),:[]{}`
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.PLUS, "+"},
		{token.MINUS, "-"},
		{token.MODULO, "%"},
		{token.ASTERISK, "*"},
		{token.SLASH, "/"},
		{token.LT, "<"},
		{token.GT, ">"},
		{token.BANG, "!"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.COMMA, ","},
		{token.COLON, ":"},
		{token.LBRACKET, "["},
		{token.RBRACKET, "]"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.EOF, ""},
	}
	runLexerTest(t, input, tests)
}

func TestMultiCharOperators(t *testing.T) {
	input := `
10 == 10
10 != 9
10 <= 11
11 >= 10
0 == 0
five += 1
ten -= 1
`
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.MONEY, "10"},
		{token.EQ, "=="},
		{token.MONEY, "10"},
		{token.MONEY, "10"},
		{token.NOT_EQ, "!="},
		{token.MONEY, "9"},
		{token.MONEY, "10"},
		{token.LTEQ, "<="},
		{token.MONEY, "11"},
		{token.MONEY, "11"},
		{token.GTEQ, ">="},
		{token.MONEY, "10"},
		{token.MONEY, "0"},
		{token.EQ, "=="},
		{token.MONEY, "0"},
		{token.IDENT, "five"},
		{token.PLUS_EQ, "+="},
		{token.MONEY, "1"},
		{token.IDENT, "ten"},
		{token.MINUS_EQ, "-="},
		{token.MONEY, "1"},
		{token.EOF, ""},
	}
	runLexerTest(t, input, tests)
}

func TestKeywordsAndLiterals(t *testing.T) {
	input := `
LOOK five IS 5!
YUGE LOSER COVFEFE
"foobar"
"a great string"
`
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "LOOK"},
		{token.IDENT, "five"},
		{token.ASSIGN, "IS"},
		{token.MONEY, "5"},
		{token.BANG, "!"},
		{token.TRUE, "YUGE"},
		{token.FALSE, "LOSER"},
		{token.NULL, "COVFEFE"},
		{token.TWEET, "foobar"},
		{token.TWEET, "a great string"},
		{token.EOF, ""},
	}
	runLexerTest(t, input, tests)
}

func TestCommentSkipping(t *testing.T) {
	input := `
// This is fake news.
LOOK // another one
FAKE_CODE
  This code? Totally fake.
  Nobody even runs this stuff.
SAD
five
`
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "LOOK"},
		{token.IDENT, "five"},
		{token.EOF, ""},
	}
	runLexerTest(t, input, tests)
}

func TestFullStatement(t *testing.T) {
	input := `
LOOK add IS MAKE_IT_BIG(x, y) {
  GIVE_ME x + y!
}!

BELIEVE_ME add(5, 5) == 10 {
	BING("success")!
} FAKE_NEWS {
	SADLY("failure")!
}
`
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "LOOK"},
		{token.IDENT, "add"},
		{token.ASSIGN, "IS"},
		{token.FUNCTION, "MAKE_IT_BIG"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "GIVE_ME"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.BANG, "!"},
		{token.RBRACE, "}"},
		{token.BANG, "!"},

		{token.IF, "BELIEVE_ME"},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.MONEY, "5"},
		{token.COMMA, ","},
		{token.MONEY, "5"},
		{token.RPAREN, ")"},
		{token.EQ, "=="},
		{token.MONEY, "10"},
		{token.LBRACE, "{"},
		{token.BING, "BING"},
		{token.LPAREN, "("},
		{token.TWEET, "success"},
		{token.RPAREN, ")"},
		{token.BANG, "!"},
		{token.RBRACE, "}"},
		{token.ELSE, "FAKE_NEWS"},
		{token.LBRACE, "{"},
		{token.SADLY, "SADLY"},
		{token.LPAREN, "("},
		{token.TWEET, "failure"},
		{token.RPAREN, ")"},
		{token.BANG, "!"},
		{token.RBRACE, "}"},
		{token.EOF, ""},
	}
	runLexerTest(t, input, tests)
}

func runLexerTest(t *testing.T, input string, tests []struct {
	expectedType    token.TokenType
	expectedLiteral string
}) {
	t.Helper()

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q. literal was %q",
				i, tt.expectedType, tok.Type, tok.Literal)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
