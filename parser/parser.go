package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
	"strconv"
)

const (
    _ int = iota
    LOWEST
    EQUALS          // ==
    LESSGREATER     // > or <
    SUM             // +
    PRODUCT         // *
    PREFIX          // -x or !x
    CALL            // function(x)
)

var precedences = map[token.TokenType]int {
    token.EQ: EQUALS,
    token.NOT_EQ: EQUALS,
    token.LT: LESSGREATER,
    token.GT: LESSGREATER,
    token.PLUS: SUM,
    token.MINUS: SUM,
    token.SLASH: PRODUCT,
    token.ASTERISK: PRODUCT,
}

type (
    prefixParseFn func() ast.Expression
    infixParseFn func(ast.Expression) ast.Expression
)
 
type Parser struct {
    l *lexer.Lexer
    currentToken token.Token
    peekToken token.Token
    errors []string
    // maps to check whether a token has any appropriate parsing function associated with it
    prefixParseFns map[token.TokenType]prefixParseFn
    infixParseFns map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
    p := &Parser {
        l: l,
        errors: []string {},
    }

    p.prefixParseFns = make(map [token.TokenType]prefixParseFn)
    p.registerPrefix(token.IDENT, p.parseIdentifier)
    p.registerPrefix(token.INT, p.parserIntegerLiteral)
    p.registerPrefix(token.BANG, p.parsePrefixExpression)
    p.registerPrefix(token.MINUS, p.parsePrefixExpression)

    p.infixParseFns = make(map [token.TokenType]infixParseFn)
    p.registerInfix(token.PLUS, p.parseInfixExpression)
    p.registerInfix(token.MINUS, p.parseInfixExpression)
    p.registerInfix(token.SLASH, p.parseInfixExpression)
    p.registerInfix(token.ASTERISK, p.parseInfixExpression)
    p.registerInfix(token.EQ, p.parseInfixExpression)
    p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
    p.registerInfix(token.LT, p.parseInfixExpression)
    p.registerInfix(token.GT, p.parseInfixExpression)

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

func (p *Parser) peekPrecedence() int {
    if p, ok := precedences[p.peekToken.Type]; ok {
        return p
    }
    return LOWEST
}

func (p *Parser) currentPrecedence() int {
    if p, ok := precedences[p.currentToken.Type]; ok {
        return p
    }
    return LOWEST
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
    p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
    p.infixParseFns[tokenType] = fn
}

func (p *Parser) parseIdentifier() ast.Expression {
    return &ast.Identifier{
        Token: p.currentToken,
        Value: p.currentToken.Literal,
    }
}

func (p *Parser) parserIntegerLiteral() ast.Expression {
    literal := &ast.IntegerLiteral{Token: p.currentToken}

    value, err := strconv.ParseInt(p.currentToken.Literal, 0, 64)
    if err != nil {
        msg := fmt.Sprintf("could not parse %q as integer...", p.currentToken.Literal)
        p.errors = append(p.errors, msg)
        return nil
    }

    literal.Value = value

    return literal
}

// ------------------------------------------------------
//                  STATEMENT PARSERS
// ------------------------------------------------------

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

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
    stmt := &ast.ExpressionStatement {Token: p.currentToken}

    stmt.Expression = p.parseExpression(LOWEST)

    // semicolon is optional
    if p.peekTokenIs(token.SEMICOLON) {
        p.nextToken()
    }

    return stmt
}

func (p *Parser) parsePrefixExpression() ast.Expression {
    expression := &ast.PrefixExpression{
        Token: p.currentToken,
        Operator: p.currentToken.Literal,
    }

    p.nextToken()

    expression.Right = p.parseExpression(PREFIX)

    return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
    expression := &ast.InfixExpression{
        Token: p.currentToken,
        Operator: p.currentToken.Literal,
        Left: left,
    }

    precedence := p.currentPrecedence()
    p.nextToken()
    expression.Right = p.parseExpression(precedence)

    return expression
}

func (p *Parser) noPrefixParseFunctionError(t token.TokenType) {
    msg := fmt.Sprintf("no prefix parse function found for %s.", t)
    p.errors = append(p.errors, msg)
}

// this function lies at the heart of Pratt parsing
func (p *Parser) parseExpression(precedence int) ast.Expression {
    prefixFn := p.prefixParseFns[p.currentToken.Type]
    if prefixFn == nil {
        p.noPrefixParseFunctionError(p.currentToken.Type)
        return nil
    }
    leftExp := prefixFn()

    for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
        infixFn := p.infixParseFns[p.peekToken.Type]
        if infixFn == nil {
            return leftExp
        }

        p.nextToken()

        leftExp = infixFn(leftExp)
    }

    return leftExp
}

func (p *Parser) parseStatement() ast.Statement {
    switch p.currentToken.Type {
        case token.LET:
            return p.parseLetStatement()
        case token.RETURN:
            return p.parseReturnStatement()
    default:
        return p.parseExpressionStatement()
    }
}

// ------------------------------------------------------
//              DRIVER FUNCTION
// ------------------------------------------------------

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
