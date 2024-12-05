package lexer

import "monkey/token"

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// Give us the next character and advance our position in the input
// FOR NOW: we only support ASCII characters.
// To support Unicode, we need to change `ch` from `byte` to `rune`, and cannot just advance readPosition by 1 byte.
// TODO: Support Unicode (and emojis?)
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // ASCII for "NUL" char;
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

// Look at the next character without advancing our position in the input.
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	// Whitespace only acts as a separator of tokens and does not have meaning, so we skip it.
	l.skipWhitespace()

	switch l.ch {
	case '=':
		tok = l.makeTwoCharToken('=', token.EQ, token.ASSIGN)
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
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '!':
		tok = l.makeTwoCharToken('=', token.NOT_EQ, token.BANG)
	case '<':
		tok = l.makeTwoCharToken('=', token.LE, token.LT)
	case '>':
		tok = l.makeTwoCharToken('=', token.GE, token.GT)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
		// Handle keywords and identifiers
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
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
	return l.input[position:l.position]
}

// This is a simple implementation that only supports integers.
// TODO: Support floats, hexadecimals, and other number formats.
func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// This little helper function is found in a lot of parsers. Which characters
// these functions actually skip depends on the language being lexed.
// Some language implementations might create tokens for newline characters and throw parsing errors if they are
// not at the correct place in the stream of tokens. We skip over newline characters to make the
// parsing step later on a little easier.
func (l *Lexer) skipWhitespace() {
	// Skip all whitespace characters
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) makeTwoCharToken(expected byte, twoCharType token.TokenType, oneCharType token.TokenType) token.Token {
	ch := l.ch
	if l.peekChar() == expected {
		l.readChar()
		return token.Token{Type: twoCharType, Literal: string(ch) + string(l.ch)}
	} else {
		return newToken(oneCharType, ch)
	}
}

// This sets the subset of characters that can make up an identifiers and keywords.
// The ch == "_" check allows underscores in identifiers and keywords!
// TODO: If we want to support other characters like `!` and `?`, we can add them here.
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}
