package ast

import "lipsoft/token"

// LetStatement 标识符、产生值的表达式、词法单元
type LetStatement struct {
	Token token.Token // token.LET词法单元
	Name  *Identifier // 绑定的标识符
	Value Expression  // 产生值的表达式
}

func (ls *LetStatement) statementNode() {}

func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}
