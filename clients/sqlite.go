package clients

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sharovik/orm/dto"
	"os"
)

//SQLiteClient the SQLite client
type SQLiteClient struct {
	Client *sql.DB
	Config DatabaseConfig
}

func (c SQLiteClient) Connect(config DatabaseConfig) (client BaseClientInterface, err error) {
	c.Config = config
	_, err = os.Stat(config.Host)
	if err != nil {
		return c, err
	}

	c.Client, err = sql.Open("sqlite3", config.Host)
	if err != nil {
		return c, err
	}

	return c, nil
}

func (c SQLiteClient) Disconnect() error {
	return c.Client.Close()
}

func (c SQLiteClient) GetClient() *sql.DB {
	return c.Client
}

func (SQLiteClient) ToSql(q QueryInterface) string {
	switch q.GetQueryType() {
	case SelectType:
		return prepareSelectQuery(q)
	case InsertType:
		return prepareInsertQuery(q)
	case DeleteType:
		return prepareDeleteQuery(q)
	case AlterType:
		return prepareAlterQuery(q)
	case UpdateType:
		return prepareUpdateQuery(q)
	case DropType:
		return prepareDropQuery(q)
	case CreateType:
		return prepareCreateQuery(q)
	}

	return ""
}

func (c SQLiteClient) Execute(q QueryInterface) (result dto.BaseResult, err error) {
	var queryStr = c.ToSql(q)
	if queryStr == "" {
		return result, errors.New("Query string cannot be empty ")
	}

	var bindings []interface{}
	for _, bind := range q.GetBindings() {
		bindings = append(bindings, bind.Value)
	}

	switch q.GetQueryType() {
	case SelectType:
		return c.executeSelect(queryStr, bindings)
	case CreateType:
		return c.executeQuery(queryStr, bindings)
	case AlterType:
		return c.executeQuery(queryStr, bindings)
	case DeleteType:
		return c.executeQuery(queryStr, bindings)
	case DropType:
		return c.executeQuery(queryStr, bindings)
	case InsertType:
		return c.executeQuery(queryStr, bindings)
	case UpdateType:
		return c.executeQuery(queryStr, bindings)
	}

	return result, nil
}

func (c SQLiteClient) executeSelect(queryStr string, bindings []interface{}) (result dto.BaseResult, err error) {
	rows, err := c.GetClient().Query(queryStr, bindings...)
	if err != nil {
		result.SetError(err)
		return result, err
	}

	columns, err := rows.Columns()
	if err != nil {
		result.SetError(err)
		return result, err
	}

	var values = make([]interface{}, len(columns))
	for i, _ := range values {
		var f interface{}
		values[i] = &f
	}

	columnTypes, err := prepareColumnTypes(rows)
	if err != nil {
		result.SetError(err)
		return result, err
	}

	for rows.Next() {
		model := new(dto.BaseModel)
		err = rows.Scan(values...)
		if err != nil {
			result.SetError(err)
			return result, err
		}

		for i, name := range columns {
			value := *(values[i].(*interface{}))
			model.AddModelField(dto.ModelField{
				Name:          name,
				Type:          columnTypes[i],
				Value:         normalizeValue(value, columnTypes[i]),
			})
		}

		result.AddItem(model)
	}

	return result, nil
}

func (c SQLiteClient) executeQuery(queryStr string, bindings []interface{}) (result dto.BaseResult, err error) {
	rows, err := c.GetClient().Exec(queryStr, bindings...)
	if err != nil {
		result.SetError(err)
		return result, err
	}

	result.InsertID, err = rows.LastInsertId()
	if err != nil {
		result.SetError(err)
		return result, err
	}

	return result, nil
}