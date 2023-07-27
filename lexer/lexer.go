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

func (l *Lexer) peekChar() byte {
    if l.readPosition >= len(l.input) {
        return 0
    } else {
        return l.input[l.readPosition]
    }
    // position and readPosition not updated here
}

func (l *Lexer) NextToken() token.Token {
    var tok token.Token

    l.skipWhitespace()

    switch l.ch {
        // operators
    case '=':
        if l.peekChar() == '=' {
            ch := l.ch
            l.readChar()
            tok = token.Token{Type: token.EQ, Literal: string(ch) + string(l.ch)}
        } else {
            tok = newToken(token.ASSIGN, l.ch)
        }
    case '+':
        tok = newToken(token.PLUS, l.ch)
    case '-':
        tok = newToken(token.MINUS, l.ch)
    case '!':
        if l.peekChar() == '=' {
            ch := l.ch
            l.readChar()
            tok = token.Token{Type: token.NOT_EQ, Literal: string(ch) + string(l.ch)}
        } else {
            tok = newToken(token.BANG, l.ch)
        }
    case '/':
        tok = newToken(token.SLASH, l.ch)
    case '*':
        tok = newToken(token.ASTERISK, l.ch)
    case '<':
        tok = newToken(token.LT, l.ch)
    case '>':
        tok = newToken(token.GT, l.ch)
        // delimiters
    case ',':
        tok = newToken(token.COMMA, l.ch)
    case ';':
        tok = newToken(token.SEMICOLON, l.ch)
    case ':':
        tok = newToken(token.COLON, l.ch)
    case '(':
        tok = newToken(token.LPAREN, l.ch)
    case ')':
        tok = newToken(token.RPAREN, l.ch)
    case '{':
        tok = newToken(token.LBRACE, l.ch)
    case '}':
        tok = newToken(token.RBRACE, l.ch)
    case '[':
        tok = newToken(token.LBRACKET, l.ch)
    case ']':
        tok = newToken(token.RBRACKET, l.ch)
    case '"':
        tok.Type = token.STRING
        tok.Literal = l.readString()
        // EOF
    case 0:
        tok.Literal = ""
        tok.Type = token.EOF
    default:
        if isLetter(l.ch) {
            tok.Literal = l.readIdentifier()               // reads rest of the word to identify whether it is a keyword or an identifier
            tok.Type = token.LookupIdentifier(tok.Literal)
            return tok                                     // return tok here incase of readIdentifier and readNumber so that readChar is not called later
        } else if isDigit(l.ch) {
            tok.Literal = l.readNumber()
            tok.Type = token.INT
            return tok
        } else {
            tok = newToken(token.ILLEGAL, l.ch)
        }
    }

    l.readChar()
    return tok
}

func (l *Lexer) readIdentifier() string {
    position := l.position

    for isLetter(l.ch) {
        l.readChar()
    }

    return l.input[position : l.position]
}

func (l *Lexer) readNumber() string {
    position := l.position

    for isDigit(l.ch) {
        l.readChar()
    }

    return l.input[position : l.position]
}

func (l *Lexer) readString() string {
    position := l.position + 1

    for {
        l.readChar()
        if l.ch == '"' || l.ch == 0 {
            break
        }
    }

    return l.input[position : l.position]
}

func (l *Lexer) skipWhitespace() {
    for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
        l.readChar()
    }
}

func isLetter(ch byte) bool {
    return ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' || ch == '_'
}

// support for floats, hex, oct and all needed
func isDigit(ch byte) bool {
    return ch >= '0' && ch <= '9'
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
    return token.Token {Type: tokenType, Literal: string(ch)}
}

func New(input string) *Lexer {
    l := &Lexer {input: input}
    l.readChar()
    return l
}

