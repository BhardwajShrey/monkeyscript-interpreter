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
