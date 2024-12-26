package ast

import "monkey/token"

type Node interface {
	TokenLiteral() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

// Root node of every AST the parser produces.
type Program struct {
	// Every valid Monkey program is a series of statements
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

type Identifier struct {
	Token token.Token
	Value string
}

type LetStatement struct {
	Token token.Token // Specifically the token.LET token
	// Hold the identifier of the binding
	Name *Identifier
	// Expression that produces the value
	Value Expression
}

type ReturnStatement struct {
	Token       token.Token // Specifically the token.RETURN token
	ReturnValue Expression
}

func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}
