package ast

import "lipsoft/token"

type ReturnStatement struct {
	Token       token.Token // return的词法单元
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}

func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}
