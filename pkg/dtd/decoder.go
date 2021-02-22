package dtd

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"unicode"
	"unicode/utf8"

	"github.com/djthorpe/data"
)

type State interface {
	// Token parses a token in current state and returns the current state,
	// some new state or nil if state is completed
	Token(rune) (State, error)
}

type Decoder struct {
	*bufio.Scanner
	stack []State
}

type rootstate struct {
	// Parse <!
	*Document
	tok string
}

// NewDecoder returns a new DTD decoder
func NewDecoder(r io.Reader) *Decoder {
	this := new(Decoder)
	this.Scanner = bufio.NewScanner(r)
	this.Scanner.Split(bufio.ScanRunes)
	return this
}

// Decode processes tokens until the scanner has no more tokens
func (this *Decoder) Decode(doc *Document) error {
	for this.Scanner.Scan() {
		tok, n := utf8.DecodeRune(this.Scanner.Bytes())
		if n == 0 {
			return data.ErrInternalAppError.WithPrefix("Decode")
		}
		state := this.CurrentState(doc)
		newstate, err := state.Token(tok)
		if err != nil {
			return err
		}
		if newstate == nil {
			// Pop last state
		} else if newstate != state {
			// Push new state
			fmt.Println("Push new state")
			this.stack = append(this.stack, newstate)
		}
	}
	// Return success
	return nil
}

// CurrentState returns the last state on the stack or creates a
// root state
func (this *Decoder) CurrentState(doc *Document) State {
	i := len(this.stack)
	if i == 0 {
		this.stack = append(this.stack, &rootstate{doc, ""})
		return this.stack[0]
	} else {
		return this.stack[i-1]
	}
}

func (this *rootstate) Token(r rune) (State, error) {
	// Return if whitespace at beginning
	if unicode.IsSpace(r) && len(this.tok) == 0 {
		return this, nil
	}
	// Append rune to token
	this.tok += string(r)
	// Return if not two chars
	if len(this.tok) < 2 {
		return this, nil
	}
	// Return error if tok is not <!
	if this.tok != "<!" {
		return nil, data.ErrBadParameter.WithPrefix(strconv.Quote(this.tok))
	}
	// TODO: Start of either a comment, element or attlist
	return this, nil
}
