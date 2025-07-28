package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENT = "IDENT"
	MONEY = "MONEY" 
	TWEET = "TWEET" 

	BANG     = "!"
	PLUS     = "+"
	MINUS    = "-"
	ASTERISK = "*"
	SLASH    = "/"
	MODULO   = "%"

	COMMA    = ","
	COLON    = ":"
	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	LT = "<"
	GT = ">"

	EQ       = "=="
	NOT_EQ   = "!="
	LTEQ     = "<="
	GTEQ     = ">="
	PLUS_EQ  = "+="
	MINUS_EQ = "-="

	ASSIGN = "="

	FUNCTION = "FUNCTION"
	RETURN   = "RETURN"
	LOOP     = "LOOP"
	LET      = "LET"
	IF       = "IF"
	ELSE_IF  = "ELSE_IF"
	ELSE     = "ELSE"
	BREAK    = "BREAK"
	CONTINUE = "CONTINUE"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	NULL     = "NULL"
	AND      = "AND"
	OR       = "OR"

	BING  = "BING"
	SADLY = "SADLY"

	FAKE_CODE = "FAKE_CODE"
)

var keywords = map[string]TokenType{
	"MAKE_IT_BIG":      FUNCTION,
	"GIVE_ME":          RETURN,
	"KEEP_WINNING":     LOOP,
	"LOOK":             LET,
	"BELIEVE_ME":       IF,
	"BUT_MAYBE":        ELSE_IF,
	"FAKE_NEWS":        ELSE,
	"IT_WAS_RIGGED":    BREAK,
	"TIRED_OF_WINNING": CONTINUE,
	"YUGE":             TRUE,
	"LOSER":            FALSE,
	"COVFEFE":          NULL,
	"IS":               ASSIGN,
	"AND":              AND,
	"OR":               OR,
	"BING":             BING,
	"SADLY":            SADLY,
	"FAKE_CODE":        FAKE_CODE,
}

func CheckKeyword(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
