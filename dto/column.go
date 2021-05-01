package dto

const (
	VarcharColumnType = "VARCHAR"
	IntegerColumnType = "INTEGER"
	BooleanColumnType = "BOOL"
)

//Column the object which can be used as main type for select queries or queries where we can specify the aliases for the queried object fields.
type Column struct {
	Field ModelField
	Alias string
}
