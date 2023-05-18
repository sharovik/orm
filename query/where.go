package query

// You can use this types for building of WHERE clause
const (
	WhereAndType = "AND"
	WhereOrType  = "OR"
	WhereNotType = "NOT"
)

// Where is an object which will be used for WHERE clause generation
type Where struct {
	First    interface{}
	Operator string
	Second   interface{}
	Type     string
}

func (w Where) GetType() string {
	if w.Type == "" {
		return WhereAndType
	}
	return w.Type
}
