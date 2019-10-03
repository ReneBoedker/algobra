package ff

// Field defines the methods that a finite field must support
type Field interface {
	Card() uint
	Char() uint
	ElementFromSigned(int) Element
	ElementFromUnsigned(uint) Element
	Elements() []Element
	MultGenerator() Element
	One() Element
	String() string
	Zero() Element
}

// Element defines the methods that an element of a finite field must support
type Element interface {
	Add(Element) Element
	Copy() Element
	Equal(Element) bool
	Err() error
	Inv() Element
	IsNonzero() bool
	IsOne() bool
	IsZero() bool
	Minus(Element) Element
	Mult(Element) Element
	Neg() Element
	Plus(Element) Element
	Pow(uint) Element
	Prod(Element, Element) Element
	SetNeg() Element
	SetUnsigned(uint)
	String() string
	Sub(Element) Element
	Times(Element) Element
}
