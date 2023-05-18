package dto

import "github.com/sharovik/orm/query"

const (
	CascadeAction  = "CASCADE"
	NoActionAction = "NO ACTION"
	SetNullAction  = "SET NULL"
)

type Index struct {
	Name   string
	Target string
	Key    string
	Unique bool
}

type ForeignKey struct {
	Name     string
	Target   query.Reference
	With     query.Reference
	OnDelete string
	OnUpdate string
}

func (f ForeignKey) GetOnDelete() string {
	if f.OnDelete == "" {
		return NoActionAction
	}

	return f.OnDelete
}

func (f ForeignKey) GetOnUpdate() string {
	if f.OnUpdate == "" {
		return NoActionAction
	}

	return f.OnUpdate
}
