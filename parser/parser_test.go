package parser

import (
    "testing"
    "monkey/ast"
    "monkey/lexer"
)

func TestLetStatements(t *testing.T) {
    input := `
    let x = 5;
    let y = 10;
    let foobar = 696969;
    `

    l := lexer.New(input)
    p := New(l)

    program := p.ParseProgram()
    checkParserErrors(t, p)

    if program == nil {
        t.Fatalf("ParseProgram() returned nil")
    }

    if len(program.Statements) != 3 {
        t.Fatalf("program.Statements does not contain 3 statements. Got %d statements.", len(program.Statements))
    }

    tests := []struct {
        expectedIdentifier string
    } {
        {"x"},
        {"y"},
        {"foobar"},
    }

    for i, tt := range tests {
        statement := program.Statements[i]

        if !testLetStatement(t, statement, tt.expectedIdentifier) {
            return
        }
    }
}

func testLetStatement(t *testing.T, s ast.Statement, expectedName string) bool {
    if s.TokenLiteral() != "let" {
        t.Errorf("s.TokenLiteral does not equal to 'let'. Got = '%q' instead.", s.TokenLiteral())
        return false
    }

    letStatement, ok := s.(*ast.LetStatement)
    if !ok {
        t.Errorf("s is not *ast.LetStatement. Got %q.", s)
        return false
    }

    if letStatement.Name.Value != expectedName {
        t.Errorf("letStatement.Name.Value is not equal to '%s'. Got '%s' instead.", expectedName, letStatement.Name.Value)
        return false
    }

    if letStatement.Name.TokenLiteral() != expectedName {
        t.Errorf("s.Name not '%s'. Got '%s' instead.", expectedName, letStatement.Name)
        return false
    }

    return true
}

func TestReturnStatements(t *testing.T) {
    input := `
    return 5;
    return 10;
    return 69420;
    `

    l := lexer.New(input)
    p := New(l)

    program := p.ParseProgram()
    checkParserErrors(t, p)

    if len(program.Statements) != 3 {
        t.Fatalf("program.Statements does not contain 3 statements. Got %d statements.", len(program.Statements))
    }

    for _, stmt := range program.Statements {
        returnStmt, ok := stmt.(*ast.ReturnStatement)

        if !ok {
            t.Errorf("stmt is not of type *ast.ReturnStatement. Got %T instead.", stmt)
            continue
        }

        if returnStmt.Token.Literal != "return" {
            t.Errorf("returnStmt.TokenLiteral() is not 'return'. Got %q instead.", returnStmt.TokenLiteral())
        }
    }
}

func TestIdentifierExpression(t *testing.T) {
    input := "foobar;"

    l := lexer.New(input)
    p := New(l)
    program := p.ParseProgram()
    checkParserErrors(t, p)

    if len(program.Statements) != 1 {
        t.Fatalf("Program does not have enough statements. Got %d statements", len(program.Statements))
    }

    stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("program.Statement[0] is not of type ast.ExpressionStatement. Got %T instead.", program.Statements[0])
    }

    ident, ok := stmt.Expression.(*ast.Identifier)
    if !ok {
        t.Fatalf("exp is not of type *ast.Identifier. Got %T instead", stmt.Expression)
    }

    if ident.Value != "foobar" {
        t.Errorf("ident.Value is not equal to %s, got '%s' instead.", "foobar", ident.Value)
    }

    if ident.TokenLiteral() != "foobar" {
        t.Errorf("ident.TokenLiteral() is not equal to %s, got '%s' instead.", "foobar", ident.TokenLiteral())
    }
}

func checkParserErrors(t *testing.T, p *Parser) {
    errors := p.Errors()

    if len(errors) == 0 {
        return
    }

    t.Errorf("parser has %d errors", len(errors))

    for _, msg := range errors {
        t.Errorf("Parser error: %q", msg)
    }

    t.FailNow()
}
