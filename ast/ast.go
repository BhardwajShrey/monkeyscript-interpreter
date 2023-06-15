package ast

import (
    "monkey/token"
    "bytes"
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

// ---------------------------------------------------------------
//              IDENTIFIER
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