package parser

import (
    "monkey/ast"
    "monkey/token"
    "monkey/lexer"
)
 
type Parser struct {
    l *lexer.Lexer
    currentToken token.Token
    peekToken token.Token
}

func New(l *lexer.Lexer) *Parser {
    p := &Parser{l: l}

    // read two tokens so that both currentToken and peekToken get set
    p.nextToken()
    p.nextToken()

    return p
}

func (p *Parser) nextToken() {
    p.currentToken = p.peekToken
    p.peekToken = p.l.NextToken()
}

func (p *Parser) parseStatement() ast.Statement {
    switch p.currentToken.Type {
        case token.LET:
        return p.parseLetStatement()
    default:
        return nil
    }
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
    statement := &ast.LetStatement{Token: p.currentToken}

    if !p.expectPeek(token.IDENT) {
        return nil
    }

    statement.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}

    if !p.expectPeek(token.ASSIGN) {
        return nil
    }

    // skip till semicolon is found. We're skipping the expression for now
    for !p.currentTokenIs(token.SEMICOLON) {
        p.nextToken()
    }

    return statement
}

func (p *Parser) currentTokenIs(t token.TokenType) bool {
    return p.currentToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
    return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
    if p.peekTokenIs(t) {
        p.nextToken()
        return true
    } else {
        return false
    }
}

func (p *Parser) ParseProgram() *ast.Program {
    program := &ast.Program{}
    program.Statements = []ast.Statement{}

    for p.currentToken.Type != token.EOF {
        statement := p.parseStatement()

        if statement != nil {
            program.Statements = append(program.Statements, statement)
        }
        p.nextToken()
    }

    return program
}
