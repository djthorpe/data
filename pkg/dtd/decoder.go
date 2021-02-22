package dtd

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
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

	// AddChild adds a child element to a parent
	AddChild(*Document, State) error
}

/////////////////////////////////////////////////////////////////////
// TYPES

type Decoder struct {
	*bufio.Scanner
	stack []State
}

type root struct {
	buf string
}

type element struct {
	// Parse <!ELEMENT ...>
	name string
}

type comment struct {
	// Parse <!--#PCDATA-->
	*root
	comment string
}

type keyword struct {
	tok string
}

type group struct {
	tok string
}

/////////////////////////////////////////////////////////////////////
// CONSTANTS

var (
	// reKeyword defines how an element name must appear
	reKeyword = regexp.MustCompile("^([A-Za-z][A-Za-z0-9-]*)$")
)

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
			parent := this.stack[len(this.stack)-2]
			// Perform action on parent when removing the child element
			if err := parent.AddChild(doc, state); err != nil {
				return err
			}
			this.stack = this.stack[:len(this.stack)-1]
		} else if newstate != state {
			// Push new state
			this.stack = append(this.stack, newstate)
		}
	}

	// The stack should have a root element on it only
	if len(this.stack) != 1 {
		return data.ErrBadParameter.WithPrefix("Decode")
	}

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
// ROOT METHODS

func (this *root) Append(r rune) {
	// Return if whitespace at beginning
	if unicode.IsSpace(r) == false || len(this.buf) > 0 {
		this.buf += string(r)
	}
}

func (this *root) Token(r rune) (State, error) {
	this.Append(r)

	switch this.buf {
	case "<!--":
		return &comment{this, ""}, nil
	case "<!ELEMENT":
		return &element{""}, nil
	default:
		return this, nil
	}
}

func (this *root) AddChild(doc *Document, state State) error {
	// Append elements to the document
	if _, ok := state.(*element); ok {
		if err := doc.append(state); err != nil {
			return err
		}
	}

	// Clear token buffer
	this.buf = ""

	// Return success
	return nil
}

/////////////////////////////////////////////////////////////////////
// COMMENT METHODS

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

/////////////////////////////////////////////////////////////////////
// ELEMENT METHODS

func (this *element) Token(r rune) (State, error) {
	// Ingest the name of the element, or check to make sure it's
	// valid element name
	if this.name == "" {
		return &keyword{""}, nil
	}

	// Consume whitespace
	if unicode.IsSpace(r) {
		return this, nil
	}

	// Next can be open bracket or keyword
	if r == '(' {
		return &group{"("}, nil
	} else {
		return &keyword{string(r)}, nil
	}
}

func (this *element) AddChild(_ *Document, state State) error {
	// Append keyword to the element
	if name, ok := state.(*keyword); ok {
		if name.Valid() == false {
			return data.ErrBadParameter.WithPrefix(strconv.Quote(name.Value()))
		} else {
			this.name = name.Value()
		}
	}

	fmt.Println("element add child=", state)

	// Return success
	return nil
}

/////////////////////////////////////////////////////////////////////
// GROUP METHODS

func (this *group) Token(r rune) (State, error) {
	if r == '(' {
		// Start new child group
		return &group{"("}, nil
	} else if r == ')' {
		return nil, nil
	} else {
		this.tok += string(r)
	}
	return this, nil
}

func (this *group) AddChild(_ *Document, state State) error {
	fmt.Println("group add child=", state)
	return nil
}

/////////////////////////////////////////////////////////////////////
// KEYWORD METHODS

func (this *keyword) Token(r rune) (State, error) {
	if unicode.IsNumber(r) || unicode.IsLetter(r) || r == '-' {
		this.tok += string(r)
		return this, nil
	} else if unicode.IsSpace(r) {
		return nil, nil
	} else {
		return nil, data.ErrBadParameter.WithPrefix(strconv.Quote(this.tok))
	}
}

func (this *keyword) AddChild(_ *Document, _ State) error {
	return nil
}

func (this *keyword) Valid() bool {
	if this.tok == "" {
		return false
	} else {
		return reKeyword.MatchString(this.tok)
	}
}

func (this *keyword) Value() string {
	return this.tok
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *root) String() string {
	str := "<root"
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
	if this.name != "" {
		str += fmt.Sprintf(" name=%q", this.name)
	}
	return str + ">"
}

func (this *keyword) String() string {
	str := "<keyword"
	if this.tok != "" {
		str += " " + strconv.Quote(this.tok)
	}
	return str + ">"
}

func (this *group) String() string {
	str := "<group"
	if this.tok != "" {
		str += " " + strconv.Quote(this.tok)
	}
	return str + ">"
}
