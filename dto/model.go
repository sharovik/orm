package dto

//ModelInterface the main interface for the object model
type ModelInterface interface {
	GetTableName() string
	SetTableName(string)
	GetColumns() []interface{}
	GetField(name string) ModelField
	UpdateFieldValue(name string, value interface{})
	AddModelField(ModelField) ModelInterface
	RemoveModelField(fieldName string) ModelInterface
	GetPrimaryKey() ModelField
	SetPrimaryKey(ModelField)
}

type BaseModel struct {
	TableName  string
	PrimaryKey ModelField
	Fields     []interface{}
}

func (m *BaseModel) SetTableName(name string) {
	m.TableName = name
}

func (m *BaseModel) GetTableName() string {
	return m.TableName
}

func (m *BaseModel) GetColumns() []interface{} {
	return m.Fields
}

func isFieldExists(columns []interface{}, field ModelField) bool {
	for _, column := range columns {
		switch v := column.(type) {
		case ModelField:
			if v.Name == field.Name {
				return true
			}
		}
	}

	return false
}

func (m *BaseModel) AddModelField(field ModelField) ModelInterface {
	var (
		columns []interface{}
		exists  = false
	)
	for _, modelField := range m.GetColumns() {
		switch v := modelField.(type) {
		case ModelField:
			if v.Name == field.Name {
				v.Value = field.Value
				if v.Name == m.PrimaryKey.Name {
					m.PrimaryKey.Value = field.Value
				}

				exists = true
			}

			columns = append(columns, v)
		}
	}

	if !exists {
		columns = append(columns, field)
	}

	m.Fields = columns

	return m
}

func (m *BaseModel) GetField(name string) ModelField {
	for _, field := range m.GetColumns() {
		switch v := field.(type) {
		case ModelField:
			if v.Name == name {
				return v
			}
		}
	}

	return ModelField{}
}

func (m *BaseModel) UpdateFieldValue(name string, value interface{}) {
	var columns []interface{}
	for _, field := range m.Fields {
		switch v := field.(type) {
		case ModelField:
			if v.Name == name {
				v.Value = value
			}

			columns = append(columns, v)
		}
	}

	m.Fields = columns
}

func (m *BaseModel) GetPrimaryKey() ModelField {
	return m.PrimaryKey
}

func (m *BaseModel) SetPrimaryKey(field ModelField) {
	field.IsPrimaryKey = true
	m.PrimaryKey = field
	m.AddModelField(field)
}

func (m *BaseModel) RemoveModelField(field string) ModelInterface {
	var columns []interface{}
	for _, f := range m.Fields {
		switch v := f.(type) {
		case ModelField:
			if field == v.Name {
				continue
			}
		}

		columns = append(columns, f)
	}

	m.Fields = columns

	return m
}
