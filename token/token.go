package token

const (
    ILLEGAL = "ILLEGAL"
    EOF = "EOF"

    // identifiers and literals
    IDENT = "IDENT"
    INT = "INT"

    // operators
    ASSIGN = "="
    PLUS = "+"

    // delimiters
    COMMA = ","
    SEMICOLON = ";"

    LPAREN = "("
    RPAREN = ")"
    LBRACE = "{"
    RBRACE = "}"

    // keywords
    FUNCTION = "FUNCTION"
    LET = "LET"
)

// can this be an enum???
type TokenType string

type Token struct {
    Type TokenType
    Literal string // value of the token
}
