package ast

import "lipsoft/token"

// Identifier 即`let x = 5;`中的x
type Identifier struct {
	Token token.Token // token.IDENT词法单元
	Value string
}

func (i *Identifier) expressionNode() {}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}
