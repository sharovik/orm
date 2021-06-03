package dto

//ModelInterface the main interface for the object model
type ModelInterface interface {
	GetTableName() string
	SetTableName(string)
	GetColumns() []interface{}
	GetField(name string) ModelField
	AddModelField(ModelField)
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

func (m BaseModel) GetTableName() string {
	return m.TableName
}

func (m BaseModel) GetColumns() []interface{} {
	var columns []interface{}

	if m.GetPrimaryKey() != (ModelField{IsPrimaryKey: true}) {
		columns = append(columns, m.GetPrimaryKey())
	}

	if len(m.Fields) == 0 {
		return nil
	}

	for _, field := range m.Fields {
		switch v := field.(type) {
		case ModelField:
			if isFieldExists(columns, v) {
				continue
			}
		}
		columns = append(columns, field)
	}
	return columns
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

func (m *BaseModel) AddModelField(field ModelField) {
	var isExistingModelField bool
	for _, modelField := range m.GetColumns() {
		switch v := modelField.(type) {
		case ModelField:
			if v.Name == field.Name {
				v.Value = field.Value
				isExistingModelField = true
			}
		}
	}

	if !isExistingModelField {
		m.Fields = append(m.GetColumns(), field)
	}
}

func (m BaseModel) GetField(name string) ModelField {
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

func (m *BaseModel) SetField(name string, value interface{}) {
	var columns []interface{}
	for _, field := range m.GetColumns() {
		switch v := field.(type) {
		case ModelField:
			if m.GetPrimaryKey() == v {
				continue
			}

			if v.Name == name {
				v.Value = value
			}
			columns = append(columns, v)
		}
	}

	m.Fields = columns
}

func (m BaseModel) GetPrimaryKey() ModelField {
	m.PrimaryKey.IsPrimaryKey = true
	return m.PrimaryKey
}

func (m *BaseModel) SetPrimaryKey(field ModelField) {
	field.IsPrimaryKey = true
	m.PrimaryKey = field
}