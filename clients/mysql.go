package clients

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sharovik/orm/dto"
	"strings"
)

//MySQLClient the SQLite client
type MySQLClient struct {
	Client *sql.DB
	Config DatabaseConfig
}

func (c MySQLClient) Connect(config DatabaseConfig) (client BaseClientInterface, err error) {
	c.Config = config
	c.Client, err = sql.Open("mysql", c.generateDSN())
	if err != nil {
		return c, err
	}

	return c, nil
}

func (c MySQLClient) generateDSN() string {
	var (
		//We init DSN with the username only
		dsn = fmt.Sprintf("%s", c.Config.Username)
	)

	if c.Config.Password != "" {
		dsn += fmt.Sprintf(":%s", c.Config.Password)
	}

	dsn += "@"

	var host = c.generateHost()
	dsn += fmt.Sprintf("%s", host)

	if c.Config.Database == "" {
		return dsn
	}

	dsn += fmt.Sprintf("/%s", c.Config.Database)
	return dsn
}

func (c MySQLClient) generateHost() string {
	if c.Config.Host == "" {
		return ""
	}

	var host = c.Config.Host
	if c.Config.Port != 0 {
		host += fmt.Sprintf(":%d", c.Config.Port)
	}

	return fmt.Sprintf("tcp(%s)", host)
}

func (c MySQLClient) Disconnect() error {
	return c.Client.Close()
}

func (c MySQLClient) GetClient() *sql.DB {
	return c.Client
}

func (c MySQLClient) ToSql(q QueryInterface) string {
	switch q.GetQueryType() {
	case SelectType:
		return prepareSelectQuery(q)
	case InsertType:
		return prepareInsertQuery(q)
	case DeleteType:
		return prepareDeleteQuery(q)
	case AlterType:
		return prepareAlterSQLStr(q)
	case UpdateType:
		return prepareUpdateQuery(q)
	case DropType:
		return prepareDropQuery(q)
	case CreateType:
		return c.prepareCreateSQLQuery(q)
	}

	return ""
}

func (c MySQLClient) Execute(q QueryInterface) (result dto.BaseResult, err error) {
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

func (c MySQLClient) executeSelect(queryStr string, bindings []interface{}) (result dto.BaseResult, err error) {
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
				Type:          getColumnTypeByValue(value),
				Value:         value,
			})
		}

		result.AddItem(model)
	}

	return result, nil
}

func (c MySQLClient) executeQuery(queryStr string, bindings []interface{}) (result dto.BaseResult, err error) {
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

//prepareCreateQuery method prepares the create query statement
func (c MySQLClient) prepareCreateSQLQuery(q QueryInterface) string {
	queryStr := fmt.Sprintf("CREATE TABLE %s (", q.GetDestination().GetTableName())

	if q.GetDestination().GetPrimaryKey() != *(new(dto.ModelField)) {
		queryStr += generateColumnSQLStr(q.GetDestination().GetPrimaryKey())

		if  len(q.GetDestination().GetColumns()) > 0 {
			queryStr += ", "
		}
	}

	if len(q.GetDestination().GetColumns()) > 0 {
		queryStr += fmt.Sprintf("%s", generateColumnsWithTypesSQLStr(q.GetDestination().GetColumns()))
	}

	if q.GetDestination().GetPrimaryKey() != *(new(dto.ModelField)) {
		queryStr += fmt.Sprintf(",\nPRIMARY KEY (%s)", q.GetDestination().GetPrimaryKey().Name)
	}

	if len(q.GetForeignKeysToAdd()) > 0 {
		queryStr += fmt.Sprintf(",\n%s", generateForeignKeysSQLStr(q.GetForeignKeysToAdd()))
	}

	if len(q.GetIndexesToAdd()) > 0 {
		queryStr += fmt.Sprintf(",\n%s", generateIndexesSQLStr(q.GetIndexesToAdd()))
	}

	queryStr += fmt.Sprintf(") ENGINE=%s DEFAULT CHARSET=%s COLLATE=%s;", c.Config.GetEngine(), c.Config.GetCharset(), c.Config.GetCollate())

	return queryStr
}

func generateColumnsWithTypesSQLStr(columns []interface{}) string {
	var result []string
	for _, column := range columns {
		switch v := column.(type) {
		case dto.ModelField:
			if v.IsPrimaryKey {
				continue
			}

			result = append(result, generateColumnSQLStr(v))
		}
	}

	return strings.Join(result, ", ")
}

func generateColumnSQLStr(column dto.ModelField) string {
	var resultStr string

	//column_1 varchar default "test" not null
	resultStr += fmt.Sprintf("%s %s", column.Name, column.Type)

	if column.Length > 0 {
		resultStr += fmt.Sprintf("(%d)", column.Length)
	}

	if column.IsUnsigned {
		resultStr += fmt.Sprintf(" unsigned")
	}

	if column.Default != nil {
		resultStr += " DEFAULT"
		switch v := column.Default.(type) {
		case int:
			resultStr += fmt.Sprintf(" %d", v)
			break
		case int64:
			resultStr += fmt.Sprintf(" %d", v)
			break
		case string:
			resultStr += fmt.Sprintf(` "%s"`, v)
			break
		case bool:
			resultStr += fmt.Sprintf(" %t", v)
			break
		}
	}

	if column.IsNullable {
		resultStr += fmt.Sprintf(" %s", "NULL")
	} else {
		resultStr += fmt.Sprintf(" %s", "NOT NULL")
	}

	if column.AutoIncrement {
		resultStr += " AUTO_INCREMENT"
	}

	return resultStr
}

func generateForeignKeysSQLStr(columns []dto.ForeignKey) string {
	var result []string
	for _, column := range columns {
		result = append(result, generateForeignKeySQL(column))
	}

	return strings.Join(result, ",\n")
}

func generateIndexesSQLStr(columns []dto.Index) string {
	var result []string
	for _, column := range columns {
		result = append(result, generateIndexSQLStr(column))
	}

	return strings.Join(result, ",\n")
}

func generateForeignKeySQL(column dto.ForeignKey) string {
	str := ""
	if column.Name != "" {
		str = fmt.Sprintf("CONSTRAINT %s ", column.Name)
	}

	str += fmt.Sprintf("FOREIGN KEY (%s) REFERENCES %s (%s)", column.With.Key, column.Target.Table, column.Target.Key)
	str += fmt.Sprintf(" ON DELETE %s ON UPDATE %s", column.GetOnDelete(), column.GetOnUpdate())
	return str
}

func generateIndexSQLStr(column dto.Index) string {
	resultStr := ""
	if column.Unique {
		resultStr = "UNIQUE "
	}

	resultStr += fmt.Sprintf("KEY %s (%s)", column.Name, column.Key)
	return resultStr
}

//prepareAlterQuery method prepares the alter query statement
func prepareAlterSQLStr(q QueryInterface) string {
	var queryStr = fmt.Sprintf("ALTER TABLE %s", q.GetDestination().GetTableName())

	var result []string
	//Generate Add columns
	if len(q.GetColumns()) > 0 {
		for _, column := range q.GetColumns() {
			switch v := column.(type) {
			case dto.ModelField:
				result = append(result, generateAlterColumnAddSQLStr(v))
			}
		}
	}

	//Generate columns to drop
	if len(q.GetColumnsToDrop()) > 0 {
		for _, column := range q.GetColumnsToDrop() {
			switch v := column.(type) {
			case dto.ModelField:
				result = append(result, fmt.Sprintf("DROP %s", v.Name))
			}
		}
	}

	//Generate indexes to add
	if len(q.GetIndexesToAdd()) > 0 {
		for _, column := range q.GetIndexesToAdd() {
			str := "ADD INDEX"
			if column.Name != "" {
				str += fmt.Sprintf(" %s", column.Name)
			}

			str += fmt.Sprintf(" %s", column.Key)
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

	//Generate foreign keys to add
	if len(q.GetForeignKeysToAdd()) > 0 {
		for _, column := range q.GetForeignKeysToAdd() {
			result = append(result, fmt.Sprintf("ADD %s", generateForeignKeySQL(column)))
		}
	}

	//Generate foreign keys to drop
	if len(q.GetForeignKeysToDrop()) > 0 {
		for _, column := range q.GetForeignKeysToDrop() {
			result = append(result, fmt.Sprintf("DROP FOREIGN KEY %s", column.Name))
		}
	}

	if len(result) > 0 {
		queryStr += fmt.Sprintf("\n%s", strings.Join(result, ","))
	}

	return queryStr
}

func generateAlterColumnAddSQLStr(column dto.ModelField) string {
	var result = "ADD "

	result += fmt.Sprintf("%s %s", column.Name, column.Type)
	if column.Length > 0 {
		result += fmt.Sprintf("(%d)", column.Length)
	}
	result += fmt.Sprintf(" %s", toSQLValue(column.Value))
	result += fmt.Sprintf(" DEFAULT %s", toSQLValue(column.Default))

	return result
}
