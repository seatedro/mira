// token/token.go
package token

type TokenType string

type Token struct {
  Type TokenType
  Literal string
}

const (
  ILLEGAL = "ILLEGAL"
  EOF = "EOF"

  // Identifiers + literals
  IDENTIFIER = "IDENTIFIER" // add, foobar, x, y, ...
  INT = "INT" // 1343456

  // Operators
  ASSIGN = "="
  PLUS = "+"

  // Delimiters
  COMMA = ","
  SEMICOLON = ";"
  LPAREN = "("
  RPAREN = ")"
  LBRACE = "{"
  RBRACE = "}"

  // Keywords
  LET = "LET"
  FUNCTION = "FUNCTION"
  RETURN = "RETURN"
)
