package token

const (
    ILLEGAL = "ILLEGAL"
    EOF = "EOF"

    // identifiers and literals
    IDENT = "IDENT"
    INT = "INT"
    STRING = "STRING"

    // operators
    ASSIGN = "="
    PLUS = "+"
    MINUS = "-"
    BANG = "!"
    ASTERISK = "*"
    SLASH = "/"
    LT = "<"
    GT = ">"
    EQ = "=="
    NOT_EQ = "!="

    // delimiters
    COMMA = ","
    SEMICOLON = ";"
    COLON = ":"

    LPAREN = "("
    RPAREN = ")"
    LBRACE = "{"
    RBRACE = "}"
    LBRACKET = "["
    RBRACKET = "]"

    // keywords
    FUNCTION = "FUNCTION"
    LET = "LET"
    TRUE = "TRUE"
    FALSE = "FALSE"
    IF = "IF"
    ELSE = "ELSE"
    RETURN = "RETURN"
)

// can this be an enum???
type TokenType string

type Token struct {
    Type TokenType
    Literal string // value of the token
}

var keywords = map[string]TokenType {
    "fn": FUNCTION,
    "let": LET,
    "true": TRUE,
    "false": FALSE,
    "if": IF,
    "else": ELSE,
    "return": RETURN,
}

func LookupIdentifier(ident string) TokenType {
    if tokenType, ok := keywords[ident]; ok {
        return tokenType
    }
    return IDENT
}
