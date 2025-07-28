package ast

import (
	"strings"

	"github.com/yurikdotdev/covfefescript/internal/token"
)

type LookStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls *LookStatement) statementNode()       {}
func (ls *LookStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LookStatement) String() string {
	return ls.TokenLiteral() + " " + ls.Name.String() + " IS " + ls.Value.String() + "!"
}

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	return rs.TokenLiteral() + " " + rs.ReturnValue.String() + "!"
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string       { return es.Expression.String() }

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	return statementsToString(bs.Statements)
}

type ForLoopStatement struct {
	Token     token.Token
	Condition Expression
	Body      *BlockStatement
}

func (fls *ForLoopStatement) statementNode()       {}
func (fls *ForLoopStatement) TokenLiteral() string { return fls.Token.Literal }
func (fls *ForLoopStatement) String() string {
	var out strings.Builder

	out.WriteString(fls.TokenLiteral() + " ")
	out.WriteString("(")
	out.WriteString(fls.Condition.String())
	out.WriteString(") ")
	out.WriteString(fls.Body.String())

	return out.String()
}

type BreakStatement struct {
	Token token.Token
}

func (br *BreakStatement) statementNode()       {}
func (br *BreakStatement) TokenLiteral() string { return br.Token.Literal }
func (br *BreakStatement) String() string {
	return br.TokenLiteral() + "!"
}

type ContinueStatement struct {
	Token token.Token
}

func (cs *ContinueStatement) statementNode()       {}
func (cs *ContinueStatement) TokenLiteral() string { return cs.Token.Literal }
func (cs *ContinueStatement) String() string {
	return cs.TokenLiteral() + "!"
}
