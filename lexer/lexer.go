package lexer

import (
	"monkey/token"
)

type Lexer struct {
    input string
    position int     // current position in input (pointer to current char)
    readPosition int // points to next character to be parsed
    ch byte          // current character being read
}

func (lex *Lexer) readChar() {
    if lex.readPosition >= len(lex.input) {
        lex.ch = 0
    } else {
        lex.ch = lex.input[lex.readPosition]
    }

    lex.position = lex.readPosition
    lex.readPosition++
}

func (l *Lexer) NextToken() token.Token {
    var tok token.Token

    switch l.ch {
    case '=':
        tok = newToken(token.ASSIGN, l.ch)
    case ';':
        tok = newToken(token.SEMICOLON, l.ch)
    case '(':
        tok = newToken(token.LPAREN, l.ch)
    case ')':
        tok = newToken(token.RPAREN, l.ch)
    case ',':
        tok = newToken(token.COMMA, l.ch)
    case '+':
        tok = newToken(token.PLUS, l.ch)
    case '{':
        tok = newToken(token.LBRACE, l.ch)
    case '}':
        tok = newToken(token.RBRACE, l.ch)
    case 0:
        tok.Literal = ""
        tok.Type = token.EOF
    }

    l.readChar()
    return tok
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
    return token.Token {Type: tokenType, Literal: string(ch)}
}

func New(input string) *Lexer {
    l := &Lexer {input: input}
    l.readChar()
    return l
}

