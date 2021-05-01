package query

const (
	OrderDirectionAsc = "ASC"
	OrderDirectionDesc = "DESC"
)

//OrderByColumn the column type for order by clause
type OrderByColumn struct {
	Direction string
	Column string
}
