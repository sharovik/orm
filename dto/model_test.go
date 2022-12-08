package dto

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBaseModel_AddModelField(t *testing.T) {
	model := new(BaseModel)
	model.SetPrimaryKey(ModelField{
		Name:          "id",
		Type:          IntegerColumnType,
		Value:         nil,
		Default:       nil,
		Length:        11,
		IsNullable:    false,
		IsPrimaryKey:  true,
		IsUnsigned:    true,
		AutoIncrement: true,
	})
	model.AddModelField(ModelField{
		Name:          "type",
		Type:          IntegerColumnType,
		Value:         2222,
		Default:       nil,
		Length:        11,
		IsNullable:    true,
		IsPrimaryKey:  false,
		IsUnsigned:    true,
		AutoIncrement: false,
	})
	model.AddModelField(ModelField{
		Name:          "name",
		Type:          VarcharColumnType,
		Value:         "test",
		Default:       "def",
		Length:        11,
		IsNullable:    false,
		IsPrimaryKey:  false,
		IsUnsigned:    false,
		AutoIncrement: false,
	})

	assert.Equal(t, ModelField{
		Name:          "id",
		Type:          IntegerColumnType,
		Value:         nil,
		Default:       nil,
		Length:        11,
		IsNullable:    false,
		IsPrimaryKey:  true,
		IsUnsigned:    true,
		AutoIncrement: true,
	}, model.GetPrimaryKey())

	assert.Equal(t, ModelField{
		Name:          "type",
		Type:          IntegerColumnType,
		Value:         2222,
		Default:       nil,
		Length:        11,
		IsNullable:    true,
		IsPrimaryKey:  false,
		IsUnsigned:    true,
		AutoIncrement: false,
	}, model.GetField("type"))

	assert.Equal(t, ModelField{
		Name:          "name",
		Type:          VarcharColumnType,
		Value:         "test",
		Default:       "def",
		Length:        11,
		IsNullable:    false,
		IsPrimaryKey:  false,
		IsUnsigned:    false,
		AutoIncrement: false,
	}, model.GetField("name"))

	assert.Equal(t, ModelField{
		Name:          "name",
		Type:          VarcharColumnType,
		Value:         "test",
		Default:       "def",
		Length:        11,
		IsNullable:    false,
		IsPrimaryKey:  false,
		IsUnsigned:    false,
		AutoIncrement: false,
	}, model.GetField("name"))

	//Now we add one more time the id field
	model.AddModelField(ModelField{
		Name:  "id",
		Value: 2222,
	})

	assert.Equal(t, ModelField{
		Name:          "id",
		Type:          IntegerColumnType,
		Value:         2222,
		Default:       nil,
		Length:        11,
		IsNullable:    false,
		IsPrimaryKey:  true,
		IsUnsigned:    true,
		AutoIncrement: true,
	}, model.GetPrimaryKey())
}
