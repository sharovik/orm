package query

const (
	LeftJoinType  = "LEFT"
	RightJoinType = "RIGHT"
	InnerJoinType = "INNER"
)

// Join the object which will be used in JOIN clause generation
type Join struct {
	Target    Reference
	With      Reference
	Condition string
	Type      string
}

// Reference the reference table struct. It can be used for definition of the related table in the join clause
type Reference struct {
	Table string
	Key   string
}
