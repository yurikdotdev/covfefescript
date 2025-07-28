package ast

import "strings"

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	return statementsToString(p.Statements)
}

func statementsToString(stmts []Statement) string {
	var out strings.Builder
	for _, s := range stmts {
		out.WriteString(s.String())
	}
	return out.String()
}
