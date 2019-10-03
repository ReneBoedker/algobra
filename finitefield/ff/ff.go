package ff

type Field interface {
	Card() uint
	Char() uint
	Elements() []Element
	MultGenerator() Element
	One() Element
	String() string
	Zero() Element
}

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
