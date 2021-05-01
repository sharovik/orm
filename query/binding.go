package query

//Bind should be used when we need to bind something in the query. Eg: we need to make sure the data we put in the query is secure
type Bind struct {
	Field string
	Value interface{}
}
