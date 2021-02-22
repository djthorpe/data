package dtd

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/djthorpe/data"
)

/////////////////////////////////////////////////////////////////////
// INTERFACES

type State interface {
	// Token parses a token in current state and returns the current state,
	// some new state or nil if state is completed
	Token(rune) (State, error)

	// Clear will clear the token buffer
	Clear()
}

/////////////////////////////////////////////////////////////////////
// TYPES

type Decoder struct {
	*bufio.Scanner
	stack []State
}

type root struct {
	// Parse <!
	buf string
}

type rule struct {
	// Parse <!--, <!ELEMENT or <!ATTLIST
	*root
}

type comment struct {
	// Parse <!--#PCDATA-->
	*root
	comment string
}

type element struct {
	*root
}

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

// NewDecoder returns a new DTD decoder
func NewDecoder(r io.Reader) *Decoder {
	this := new(Decoder)
	this.Scanner = bufio.NewScanner(r)
	this.Scanner.Split(bufio.ScanRunes)
	return this
}

/////////////////////////////////////////////////////////////////////
// DECODER

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
			// Eject instruction, clear token buffer and pop state
			if err := doc.append(state); err != nil {
				return err
			}
			state.Clear()
			this.stack = this.stack[:len(this.stack)-1]
		} else if newstate != state {
			// Push new state
			this.stack = append(this.stack, newstate)
		}
	}

	// Print stack - TODO: return error if stack is not empty
	fmt.Println("stack=", this.stack)

	// Return success
	return nil
}

/////////////////////////////////////////////////////////////////////
// METHODS

// CurrentState returns the last state on the stack or creates a
// root state
func (this *Decoder) CurrentState(doc *Document) State {
	i := len(this.stack)
	if i == 0 {
		this.stack = append(this.stack, &root{""})
		return this.stack[0]
	} else {
		return this.stack[i-1]
	}
}

/////////////////////////////////////////////////////////////////////
// STATE METHODS

func (this *root) Append(r rune) {
	// Return if whitespace at beginning
	if unicode.IsSpace(r) == false || len(this.buf) > 0 {
		this.buf += string(r)
	}
}

func (this *root) Clear() {
	this.buf = ""
}

func (this *root) Token(r rune) (State, error) {
	this.Append(r)

	// Return if not two chars
	if len(this.buf) < 2 {
		return this, nil
	}

	// Return error if tok is not <! else create a rule
	if this.buf != "<!" {
		return nil, data.ErrBadParameter.WithPrefix(strconv.Quote(this.buf))
	}

	return &rule{this}, nil
}

func (this *rule) Token(r rune) (State, error) {
	this.Append(r)

	switch this.buf {
	case "<!--":
		return &comment{this.root, ""}, nil
	case "<!ELEMENT":
		return &element{this.root}, nil
	default:
		return this, nil
	}
}

func (this *comment) Token(r rune) (State, error) {
	this.Append(r)

	// Consume tokens until buffer ends with -->
	if strings.HasSuffix(this.buf, "-->") {
		this.comment = this.buf
		return nil, nil
	} else {
		return this, nil
	}
}

func (this *element) Token(r rune) (State, error) {
	this.Append(r)

	fmt.Printf("element: %q\n", this.buf)

	return this, nil
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *root) String() string {
	str := "<root"
	return str + ">"
}

func (this *rule) String() string {
	str := "<rule"
	return str + ">"
}

func (this *comment) String() string {
	str := "<comment"
	if this.comment != "" {
		str += " " + strconv.Quote(this.comment)
	}
	return str + ">"
}

func (this *element) String() string {
	str := "<element"
	return str + ">"
}
