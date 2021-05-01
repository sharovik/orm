package dto

//ModelField interface for the model fields
type ModelField struct {
	Name string
	Type string
	Value interface{}
	Default interface{}
	Length int64
	IsNullable bool
	IsPrimaryKey bool
	AutoIncrement bool
}
