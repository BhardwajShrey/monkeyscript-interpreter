package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
)
 
type Parser struct {
    l *lexer.Lexer
    currentToken token.Token
    peekToken token.Token
    errors []string
}

func New(l *lexer.Lexer) *Parser {
    p := &Parser {
        l: l,
        errors: []string {},
    }

    // read two tokens so that both currentToken and peekToken get set
    p.nextToken()
    p.nextToken()

    return p
}

func (p *Parser) nextToken() {
    p.currentToken = p.peekToken
    p.peekToken = p.l.NextToken()
}

func (p *Parser) Errors() []string {
    return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
    msg := fmt.Sprintf("Expected next token to be %s, got %s instead...", t, p.peekToken.Type)
    p.errors = append(p.errors, msg)
}

func (p *Parser) currentTokenIs(t token.TokenType) bool {
    return p.currentToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
    return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
    if p.peekTokenIs(t) {
        // notice that we sneakily move on to the next token here
        p.nextToken()
        return true
    } else {
        p.peekError(t)
        return false
    }
}

func (p *Parser) parseStatement() ast.Statement {
    switch p.currentToken.Type {
        case token.LET:
            return p.parseLetStatement()
        case token.RETURN:
            return p.parseReturnStatement()
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

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
    statement := &ast.ReturnStatement{Token: p.currentToken}

    p.nextToken()

    // skipping expressions for now, till semicolon is found
    for !p.currentTokenIs(token.SEMICOLON) {
        p.nextToken()
    }

    return statement
}

func (p *Parser) ParseProgram() *ast.Program {
    program := &ast.Program{}
    program.Statements = []ast.Statement{}

    for !p.currentTokenIs(token.EOF) {
        statement := p.parseStatement()

        if statement != nil {
            program.Statements = append(program.Statements, statement)
        }
        p.nextToken()
    }

    return program
}
