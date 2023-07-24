package ast

import (
	"bytes"
	"monkey/token"
	"strings"
)

type Node interface {
    TokenLiteral() string       // returns literal value of the token which the node is associated with (for debugging and testing purposes)
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

// Program node is root node of AST
type Program struct {
    Statements []Statement
}

func (p *Program) TokenLiteral() string {
    if len(p.Statements) > 0 {
        return p.Statements[0].TokenLiteral()
    } else {
        return ""
    }
}

func (p *Program) String() string {
    var output bytes.Buffer

    for _, stmt := range p.Statements {
        output.WriteString(stmt.String())
    }

    return output.String()
}

// ---------------------------------------------------------------
//              STATEMENTS
// ---------------------------------------------------------------

type LetStatement struct {
    Token token.Token   // token.LET token
    Name *Identifier
    Value Expression
}

func (ls *LetStatement) statementNode() {}

func (ls *LetStatement) TokenLiteral() string {
    return ls.Token.Literal
}

func (ls *LetStatement) String() string {
    var output bytes.Buffer

    output.WriteString(ls.TokenLiteral() + " ")
    output.WriteString(ls.Name.String())
    output.WriteString(" = ")

    if ls.Value != nil {
        output.WriteString(ls.Value.String())
    }

    output.WriteString(";")

    return output.String()
}


type ReturnStatement struct {
    Token token.Token   // token.RETURN token
    ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}

func (rs *ReturnStatement) TokenLiteral() string {
    return rs.Token.Literal
}

func (rs *ReturnStatement) String() string {
    var output bytes.Buffer

    output.WriteString(rs.TokenLiteral() + " ")

    if rs.ReturnValue != nil {
        output.WriteString(rs.ReturnValue.String())
    }

    output.WriteString(";")

    return output.String()
}


// for statements like `x + 10;`
type ExpressionStatement struct {
    Token token.Token
    Expression Expression
}

func (es *ExpressionStatement) statementNode() {}

func (es *ExpressionStatement) TokenLiteral() string {
    return es.Token.Literal
}

func (es *ExpressionStatement) String() string {
    var output bytes.Buffer

    if es.Expression != nil {
        output.WriteString(es.Expression.String())
    }

    return output.String()
}

type BlockStatement struct {
    Token token.Token       // the { token
    Statements []Statement
}

func (bs *BlockStatement) statementNode() {}

func (bs *BlockStatement) TokenLiteral() string {
    return bs.Token.Literal
}

func (bs *BlockStatement) String() string {
    var out bytes.Buffer

    for _, s := range bs.Statements {
        out.WriteString(s.String())
    }

    return out.String()
}

// ---------------------------------------------------------------
//              EXPRESSIONS
// ---------------------------------------------------------------

type Identifier struct {
    Token token.Token   // token.IDENT token
    Value string
}

func (i *Identifier) expressionNode() {}

func (i *Identifier) TokenLiteral() string {
    return i.Token.Literal
}

func (i *Identifier) String() string {
    return i.Value
}

type IntegerLiteral struct {
    Token token.Token
    Value int64
}

func (il *IntegerLiteral) expressionNode() {}

func (il *IntegerLiteral) TokenLiteral() string {
    return il.Token.Literal
}

func (il *IntegerLiteral) String() string {
    return il.Token.Literal
}

type StringLiteral struct {
    Token token.Token
    Value string
}

func (sl *StringLiteral) expressionNode() {}

func (sl *StringLiteral) TokenLiteral() string {
    return sl.Token.Literal
}

func (sl *StringLiteral) String() string {
    return sl.Token.Literal
}

type PrefixExpression struct {
    Token token.Token       // prefix token e.g. "!"
    Operator string
    Right Expression
}

func (pe *PrefixExpression) expressionNode() {}

func (pe *PrefixExpression) TokenLiteral() string {
    return pe.Token.Literal
}

func (pe *PrefixExpression) String() string {
    var out bytes.Buffer

    out.WriteString("(")
    out.WriteString(pe.Operator)
    out.WriteString(pe.Right.String())
    out.WriteString(")")

    return out.String()
}

type InfixExpression struct {
    Token token.Token       // operator token e.g. '+', '*'
    Left Expression
    Operator string
    Right Expression
}

func (ie *InfixExpression) expressionNode() {}

func (ie *InfixExpression) TokenLiteral() string {
    return ie.Token.Literal
}

func (ie *InfixExpression) String() string {
    var out bytes.Buffer

    out.WriteString("(")
    out.WriteString(ie.Left.String())
    out.WriteString(" " + ie.Operator + " ")
    out.WriteString(ie.Right.String())
    out.WriteString(")")

    return out.String()
}

type Boolean struct {
    Token token.Token
    Value bool
}

func (b *Boolean) expressionNode() {}

func (b *Boolean) TokenLiteral() string {
    return b.Token.Literal
}

func (b *Boolean) String() string {
    return b.Token.Literal
}

type IfExpression struct {
    Token token.Token       // the IF token
    Condition Expression
    Consequence *BlockStatement
    Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode() {}

func (ie *IfExpression) TokenLiteral() string {
    return ie.Token.Literal
}

func (ie *IfExpression) String() string {
    var out bytes.Buffer

    out.WriteString("if ")
    out.WriteString(ie.Condition.String())
    out.WriteString(" ")
    out.WriteString(ie.Consequence.String())

    if ie.Alternative != nil {
        out.WriteString("else ")
        out.WriteString(ie.Alternative.String())
    }

    return out.String()
}

type FunctionLiteral struct {
    Token token.Token       // the Fn token
    Parameters []*Identifier
    Body *BlockStatement
}

func (fl *FunctionLiteral) expressionNode() {}

func (fl *FunctionLiteral) TokenLiteral() string {
    return fl.Token.Literal
}

func (fl *FunctionLiteral) String() string {
    var out bytes.Buffer

    params := []string{}

    for _, p := range fl.Parameters {
        params = append(params,p.String())
    }

    out.WriteString(fl.TokenLiteral())
    out.WriteString("(")
    out.WriteString(strings.Join(params, ", "))
    out.WriteString(")")
    out.WriteString(fl.Body.String())

    return out.String()
}

type CallExpression struct {
    Token token.Token       // the '(' token
    Function Expression     // Identifier or FunctionLiteral
    Arguments []Expression
}

func (ce *CallExpression) expressionNode() {}

func (ce *CallExpression) TokenLiteral() string {
    return ce.Token.Literal
}

func (ce *CallExpression) String() string {
    var out bytes.Buffer

    args := []string{}
    for _, arg := range ce.Arguments {
        args = append(args, arg.String())
    }

    out.WriteString(ce.Function.String())
    out.WriteString("(")
    out.WriteString(strings.Join(args, ", "))
    out.WriteString(")")

    return out.String()
}

type ArrayLiteral struct {
    Token token.Token       // the [ token
    Elements []Expression
}

func (al *ArrayLiteral) expressionNode() {}

func (al *ArrayLiteral) TokenLiteral() string {
    return al.Token.Literal
}

func (al *ArrayLiteral) String() string {
    var out bytes.Buffer

    elements := []string{}

    for _, el := range al.Elements {
        elements = append(elements, el.String())
    }

    out.WriteString("[")
    out.WriteString(strings.Join(elements, ", "))
    out.WriteString("]")

    return out.String()
}

type IndexExpression struct {
    Token token.Token       // the [ token
    Left Expression
    Index Expression
}

func (ie *IndexExpression) expressionNode() {}

func (ie *IndexExpression) TokenLiteral() string {
    return ie.Token.Literal
}

func (ie *IndexExpression) String() string {
    var out bytes.Buffer

    out.WriteString("(")
    out.WriteString(ie.Left.String())
    out.WriteString("[")
    out.WriteString(ie.Index.String())
    out.WriteString("])")

    return out.String()
}
