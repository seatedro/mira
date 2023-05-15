// lexer/lexer.go

package lexer

type Lexer struct {
	input        string
	position     int  // Current position in input
	readPosition int  // Current reading position (after current char)
	ch           byte // Current character
}

func New(input string) *Lexer {
  l := &Lexer{input: input}
  return l
}
