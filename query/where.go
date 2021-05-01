package query

//Where is an object which will be used for WHERE clause generation
type Where struct {
	First interface{}
	Operator string
	Second interface{}
}
