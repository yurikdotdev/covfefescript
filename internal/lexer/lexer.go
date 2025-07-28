package lexer

import "github.com/yurikdotdev/covfefescript/internal/token"

type Lexer struct {
	input        string
	currPosition int
	nextPosition int
	char         byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.nextPosition >= len(l.input) {
		l.char = 0
	} else {
		l.char = l.input[l.nextPosition]
	}
	l.currPosition = l.nextPosition
	l.nextPosition += 1
}

func (l *Lexer) peekChar() byte {
	if l.nextPosition >= len(l.input) {
		return 0
	}
	return l.input[l.nextPosition]
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.char {
	case '=':
		tok = l.newCompoundToken('=', token.EQ, token.ILLEGAL)
	case '+':
		tok = l.newCompoundToken('=', token.PLUS_EQ, token.PLUS)
	case '-':
		tok = l.newCompoundToken('=', token.MINUS_EQ, token.MINUS)
	case '!':
		tok = l.newCompoundToken('=', token.NOT_EQ, token.BANG)
	case '<':
		tok = l.newCompoundToken('=', token.LTEQ, token.LT)
	case '>':
		tok = l.newCompoundToken('=', token.GTEQ, token.GT)
	case '/':
		if l.peekChar() == '/' {
			l.skipSingleLineComment()
			return l.NextToken()
		}
		tok = newToken(token.SLASH, l.char)
	case '"':
		tok.Type = token.TWEET
		tok.Literal = l.readString()
	case '(':
		tok = newToken(token.LPAREN, l.char)
	case ')':
		tok = newToken(token.RPAREN, l.char)
	case '{':
		tok = newToken(token.LBRACE, l.char)
	case '}':
		tok = newToken(token.RBRACE, l.char)
	case ',':
		tok = newToken(token.COMMA, l.char)
	case ':':
		tok = newToken(token.COLON, l.char)
	case '[':
		tok = newToken(token.LBRACKET, l.char)
	case ']':
		tok = newToken(token.RBRACKET, l.char)
	case '*':
		tok = newToken(token.ASTERISK, l.char)
	case '%':
		tok = newToken(token.MODULO, l.char)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.char) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.CheckKeyword(tok.Literal)

			if tok.Type == token.FAKE_CODE {
				l.skipMultiLineComment()
				return l.NextToken()
			}

			return tok
		} else if isDigit(l.char) {
			tok.Type = token.MONEY
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.char)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) readIdentifier() string {
	position := l.currPosition
	for isLetter(l.char) || isDigit(l.char) {
		l.readChar()
	}
	return l.input[position:l.currPosition]
}

func (l *Lexer) readNumber() string {
	position := l.currPosition
	for isDigit(l.char) {
		l.readChar()
	}
	return l.input[position:l.currPosition]
}

func (l *Lexer) readString() string {
	position := l.currPosition + 1
	for {
		l.readChar()
		if l.char == '"' || l.char == 0 {
			break
		}
	}
	return l.input[position:l.currPosition]
}

func (l *Lexer) skipWhitespace() {
	for l.char == ' ' || l.char == '\t' || l.char == '\n' || l.char == '\r' {
		l.readChar()
	}
}

func (l *Lexer) skipSingleLineComment() {
	for l.char != '\n' && l.char != 0 {
		l.readChar()
	}
}

func (l *Lexer) skipMultiLineComment() {
	for {
		if l.char == 'S' && l.peekChar() == 'A' {
			if l.nextPosition+1 < len(l.input) && l.input[l.nextPosition+1] == 'D' {
				l.readChar()
				l.readChar()
				l.readChar()
				break
			}
		}

		l.readChar()
		if l.char == 0 {
			break
		}
	}
}

func (l *Lexer) newCompoundToken(peek byte, twoCharType token.TokenType, oneCharType token.TokenType) token.Token {
	if l.peekChar() == peek {
		ch := l.char
		l.readChar()
		literal := string(ch) + string(l.char)
		return token.Token{Type: twoCharType, Literal: literal}
	}
	return newToken(oneCharType, l.char)
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
