package clients

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sharovik/orm/dto"
)

// SQLiteClient the SQLite client
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

func (c SQLiteClient) ToSql(q QueryInterface) string {
	return toSql(c, q)
}

func (c SQLiteClient) prepareTransactionBegin() string {
	return "BEGIN TRANSACTION;"
}

func (c SQLiteClient) prepareTransactionCommit() string {
	return "COMMIT;"
}

func (c SQLiteClient) prepareTransactionRollback() string {
	return "ROLLBACK;"
}

// prepareCreateSQLQuery method prepares the create query statement
func (c SQLiteClient) prepareCreateQuery(q QueryInterface) string {
	ifNotExists := ""
	if q.GetIfNotExists() {
		ifNotExists = "IF NOT EXISTS "
	}

	queryStr := fmt.Sprintf("CREATE TABLE %s%s (", ifNotExists, q.GetDestination().GetTableName())

	if q.GetDestination().GetPrimaryKey() != *(new(dto.ModelField)) {
		queryStr += fmt.Sprintf("%s %s CONSTRAINT %s_pk primary key", q.GetDestination().GetPrimaryKey().Name, q.GetDestination().GetPrimaryKey().Type, q.GetDestination().GetTableName())
		if q.GetDestination().GetPrimaryKey().AutoIncrement {
			queryStr += " autoincrement"
		}
	}

	if len(q.GetDestination().GetColumns()) > 0 {
		queryStr += fmt.Sprintf(", %s", generateColumnsWithTypesStr(q.GetDestination().GetColumns()))
	}

	if len(q.GetForeignKeysToAdd()) > 0 {
		queryStr += fmt.Sprintf(",\n%s", generateForeignKeysStr(q.GetForeignKeysToAdd()))
	}

	queryStr += ");"

	if len(q.GetIndexesToAdd()) > 0 {
		queryStr += fmt.Sprintf(" %s", generateIndexesStr(q.GetIndexesToAdd()))
	}

	return queryStr
}

// prepareAlterSQLStr method prepares the alter query statement
func (c SQLiteClient) prepareAlterQuery(q QueryInterface) string {
	var queryStr = ""

	if isNewSchemaShouldBeGenerated(q) {
		//We first generate the "create" statement for the new table
		qb := buildTempTableSQLiteQuery(q)
		queryStr = fmt.Sprintf("%s\n", c.prepareCreateQuery(qb))

		//Then we insert the data from the old table into the new table
		var selectColumns []interface{}
		for _, column := range qb.GetDestination().GetColumns() {
			switch v := column.(type) {
			case dto.ModelField:
				if v.AutoIncrement {
					break
				}

				selectColumns = append(selectColumns, v.Name)
			}
		}
		selQb := (new(Query)).Select(selectColumns).From(q.GetDestination())
		inQb := (new(Query)).Insert(qb.GetDestination()).Values(selQb)
		queryStr += fmt.Sprintf("%s;\n", prepareInsertQuery(inQb))

		//Now we need to switch the names of the new and the old tables
		queryStr += fmt.Sprintf("%s;\n", prepareRenameTableQuery(new(Query).
			Rename(q.GetDestination().GetTableName(), fmt.Sprintf("%s%s", OldTablePrefix, q.GetDestination().GetTableName())),
		))
		queryStr += fmt.Sprintf("%s;\n", prepareRenameTableQuery(new(Query).
			Rename(qb.GetDestination().GetTableName(), q.GetDestination().GetTableName()),
		))

		//We drop the old table
		queryStr += fmt.Sprintf("%s;", prepareDropQuery(new(Query).Drop(&dto.BaseModel{
			TableName: fmt.Sprintf("%s%s", OldTablePrefix, q.GetDestination().GetTableName()),
		})))

		return queryStr
	}

	var result []string
	if len(q.GetColumns()) > 0 {
		queryStr = fmt.Sprintf("ALTER TABLE %s ", q.GetDestination().GetTableName())
		for _, column := range q.GetColumns() {
			switch v := column.(type) {
			case dto.ModelField:
				result = append(result, fmt.Sprintf("ADD COLUMN %s", generateColumnStr(v)))
			}
		}
	}

	//Generate indexes to add
	if len(q.GetIndexesToAdd()) > 0 {
		for _, column := range q.GetIndexesToAdd() {
			str := "CREATE"
			if column.Unique {
				str += " UNIQUE"
			}

			str += " INDEX"
			if column.Name != "" {
				str += fmt.Sprintf(" %s", column.Name)
			}

			str += fmt.Sprintf(" on %s (%s)", q.GetDestination().GetTableName(), column.Key)
			result = append(result, str)
		}
	}

	//Generate indexes to add
	if len(q.GetIndexesToDrop()) > 0 {
		for _, column := range q.GetIndexesToDrop() {
			key := column.Name
			if key == "" {
				key = column.Key
			}

			result = append(result, fmt.Sprintf("DROP INDEX %s", key))
		}
	}

	if len(result) > 0 {
		queryStr += strings.Join(result, ";\n")
	}

	return queryStr
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
		result, err = c.executeSelect(queryStr, bindings)
	case CreateType:
		result, err = c.executeQuery(queryStr, bindings)
	case AlterType:
		result, err = c.executeQuery(queryStr, bindings)
	case RenameType:
		result, err = c.executeQuery(queryStr, bindings)
	case DeleteType:
		result, err = c.executeQuery(queryStr, bindings)
	case DropType:
		result, err = c.executeQuery(queryStr, bindings)
	case InsertType:
		result, err = c.executeQuery(queryStr, bindings)
	case UpdateType:
		result, err = c.executeQuery(queryStr, bindings)
	}

	if err != nil {
		return result, err
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
	for i := range values {
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
				Name:  name,
				Type:  columnTypes[i],
				Value: normalizeValue(value, columnTypes[i]),
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
