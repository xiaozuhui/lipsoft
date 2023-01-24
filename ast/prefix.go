package ast

import (
	"bytes"
	"lipsoft/token"
)

type PrefixExpression struct {
	Token    token.Token // 前缀词法单元，如!
	Operator string      // 是包含！或-的字符串
	Right    Expression  // 右侧的表达式
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}
