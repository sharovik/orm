package clients

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sharovik/orm/dto"
	"github.com/sharovik/orm/query"
	"os"
	"strings"
)

//SQLiteClient the SQLite client
type SQLiteClient struct {
	Client *sql.DB
	Config DatabaseConfig
}

func (c SQLiteClient) Connect(config DatabaseConfig) (client SQLiteClient, err error) {
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
	return nil
}

func (c SQLiteClient) GetClient() *sql.DB {
	return c.Client
}

//prepareSelectQuery method prepares the select query statement
func prepareSelectQuery(q QueryInterface) string {
	var queryStr = "SELECT "

	//Target we need to prepare the select columns list
	queryStr += generateSelectColumnsStr(q.GetColumns())

	//Now we need to create FROM string
	queryStr += fmt.Sprintf(" FROM %s", q.GetDestination().GetTableName())

	//Next step is appending the join statements if there are joins specified
	if len(q.GetJoins()) > 0 {
		queryStr += fmt.Sprintf(" %s", generateJoinsStr(q.GetJoins()))
	}

	if len(q.GetWheres()) > 0 {
		queryStr += fmt.Sprintf(" %s", generateWhereStr(q.GetWheres()))
	}

	if len(q.GetGroupBy()) > 0 {
		queryStr += fmt.Sprintf(" %s", generateGroupByStr(q.GetGroupBy()))
	}

	if len(q.GetOrderBy()) > 0 {
		queryStr += fmt.Sprintf(" %s", generateOrderByStr(q.GetOrderBy()))
	}

	if q.GetLimit() != *new(query.Limit) {
		queryStr += fmt.Sprintf(" %s", generateLimitStr(q.GetLimit()))
	}

	return queryStr
}

func generateJoinsStr(joins []query.Join) string {
	var joinsStr string
	for _, join := range joins {
		joinsStr += fmt.Sprintf("%s JOIN %s ON (%s %s %s)",
			strings.ToUpper(join.Type),
			join.Target.Table,
			fmt.Sprintf("%s.%s", join.Target.Table, join.Target.Key),
			join.Condition,
			fmt.Sprintf("%s.%s", join.With.Table, join.With.Key),
		)
	}

	return joinsStr
}

func generateWhereStr(wheres []query.Where) string {
	var preparedWheres []string

	for _, where := range wheres {
		whereStr := fmt.Sprintf("%s %s %s", where.First, where.Operator, where.Second)
		preparedWheres = append(preparedWheres, whereStr)
	}

	return fmt.Sprintf("WHERE %s", strings.Join(preparedWheres, " AND "))
}

func generateSelectColumnsStr(columns []interface{}) string {
	if len(columns) == 0 {
		return "*"
	}

	var preparedColumns []string
	//Target we need to prepare the select columns list
	for _, column := range columns {
		switch v := column.(type) {
		case string:
			preparedColumns = append(preparedColumns, v)
		case dto.ModelField:
			preparedColumns = append(preparedColumns, v.Name)
		}
	}

	return strings.Join(preparedColumns, ", ")
}

func generateGroupByStr(groupBys []string) string {
	var preparedColumns []string
	//Target we need to prepare the select columns list
	for _, column := range groupBys {
		preparedColumns = append(preparedColumns, column)
	}

	return fmt.Sprintf("GROUP BY %s", strings.Join(preparedColumns, ", "))
}

func generateOrderByStr(orderBys []query.OrderByColumn) string {
	var preparedColumns []string
	//Target we need to prepare the select columns list
	for _, column := range orderBys {
		preparedColumns = append(preparedColumns, fmt.Sprintf("%s %s", column.Column, column.Direction))
	}

	return fmt.Sprintf("ORDER BY %s", strings.Join(preparedColumns, ", "))
}

func generateLimitStr(limit query.Limit) string {
	if limit.From == 0 && limit.To == 0 {
		return ""
	}

	if limit.From == 0 && limit.To > 0 {
		return fmt.Sprintf("LIMIT %d", limit.To)
	}

	return fmt.Sprintf("LIMIT %d, %d", limit.From, limit.To)
}

func generateBindingsStr(bindings []query.Bind) string {
	var items []string
	for _, bind := range bindings {
		items = append(items, bind.Field)
	}

	return strings.Join(items, ", ")
}

//prepareInsertQuery method prepares the insert query statement
func prepareInsertQuery(q QueryInterface) string {
	var queryStr = fmt.Sprintf("INSERT INTO %s", q.GetDestination().GetTableName())

	var schema []string
	for _, column := range q.GetColumns() {
		switch v := column.(type) {
		case dto.ModelField:
			if v.AutoIncrement {
				break
			}

			schema = append(schema, v.Name)
		}
	}

	queryStr += fmt.Sprintf(" (%s)", strings.Join(schema, ", "))

	queryStr += fmt.Sprintf(" VALUES (%s)", generateBindingsStr(q.GetBindings()))

	return queryStr
}

//prepareUpdateQuery method prepares the update query statement
func prepareUpdateQuery(q QueryInterface) string {
	queryStr := fmt.Sprintf("UPDATE %s SET", q.GetDestination().GetTableName())

	var toUpdate []string
	for i, column := range q.GetColumns() {
		switch v := column.(type) {
		case dto.ModelField:
			if v.IsPrimaryKey {
				continue
			}

			if q.GetBindings()[i] == *(new(query.Bind)) {
				break
			}

			toUpdate = append(toUpdate, fmt.Sprintf("%s = %s", v.Name, q.GetBindings()[i].Field))
		}
	}

	queryStr += fmt.Sprintf(" %s", strings.Join(toUpdate, ", "))

	if len(q.GetJoins()) > 0 {
		queryStr += fmt.Sprintf(" %s", generateJoinsStr(q.GetJoins()))
	}

	if len(q.GetWheres()) > 0 {
		queryStr += fmt.Sprintf(" %s", generateWhereStr(q.GetWheres()))
	}

	return queryStr
}

//prepareDeleteQuery method prepares the delete query statement
func prepareDeleteQuery(q QueryInterface) string {
	queryStr := fmt.Sprintf("DELETE FROM %s", q.GetDestination().GetTableName())

	if len(q.GetJoins()) > 0 {
		queryStr += fmt.Sprintf(" %s", generateJoinsStr(q.GetJoins()))
	}

	if len(q.GetWheres()) > 0 {
		queryStr += fmt.Sprintf(" %s", generateWhereStr(q.GetWheres()))
	}

	if len(q.GetGroupBy()) > 0 {
		queryStr += fmt.Sprintf(" %s", generateGroupByStr(q.GetGroupBy()))
	}

	if len(q.GetOrderBy()) > 0 {
		queryStr += fmt.Sprintf(" %s", generateOrderByStr(q.GetOrderBy()))
	}

	if q.GetLimit() != *new(query.Limit) {
		queryStr += fmt.Sprintf(" %s", generateLimitStr(q.GetLimit()))
	}

	return queryStr
}

//prepareCreateQuery method prepares the create query statement
func prepareCreateQuery(q QueryInterface) string {
	queryStr := fmt.Sprintf("CREATE TABLE %s (", q.GetDestination().GetTableName())

	if q.GetDestination().GetPrimaryKey() != *(new(dto.ModelField)) {
		queryStr += fmt.Sprintf("%s %s CONSTRAINT %s_pk primary key", q.GetDestination().GetPrimaryKey().Name, q.GetDestination().GetPrimaryKey().Type, q.GetDestination().GetTableName())
		if q.GetDestination().GetPrimaryKey().AutoIncrement {
			queryStr += " autoincrement"
		}

		queryStr += ","
	}

	if len(q.GetDestination().GetColumns()) > 0 {
		queryStr += fmt.Sprintf(" %s", generateColumnsWithTypesStr(q.GetDestination().GetColumns()))
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

func generateColumnsWithTypesStr(columns []interface{}) string {
	var result []string
	for _, column := range columns {
		switch v := column.(type) {
		case dto.ModelField:
			if v.IsPrimaryKey {
				continue
			}

			result = append(result, generateColumnStr(v))
		}
	}

	return strings.Join(result, ", ")
}

func generateColumnStr(column dto.ModelField) string {
	var resultStr string

	//column_1 varchar default "test" not null
	resultStr += fmt.Sprintf("%s %s", column.Name, column.Type)

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

	return resultStr
}

func generateForeignKeysStr(columns []dto.ForeignKey) string {
	var result []string
	for _, column := range columns {
		result = append(result, generateForeignKey(column))
	}

	return strings.Join(result, ",\n")
}

func generateIndexesStr(columns []dto.Index) string {
	var result []string
	for _, column := range columns {
		result = append(result, generateIndexStr(column))
	}

	return strings.Join(result, ",\n")
}

func generateForeignKey(column dto.ForeignKey) string {
	str := ""
	if column.Name != "" {
		str = fmt.Sprintf("CONSTRAINT %s\n", column.Name)
	}

	str += fmt.Sprintf("FOREIGN KEY (%s)\n REFERENCES %s (%s)\n", column.With.Key, column.Target.Table, column.Target.Key)
	str += fmt.Sprintf("ON DELETE %s\nON UPDATE %s", column.GetOnDelete(), column.GetOnUpdate())
	return str
}

func generateIndexStr(column dto.Index) string {
	resultStr := "CREATE "
	if column.Unique {
		resultStr += "UNIQUE "
	}

	resultStr += fmt.Sprintf("INDEX %s \nON %s (%s);", column.Name, column.Target, column.Key)
	return resultStr
}

//prepareAlterQuery method prepares the alter query statement
func prepareAlterQuery(q QueryInterface) string {
	var queryStr = fmt.Sprintf("ALTER TABLE %s", q.GetDestination().GetTableName())

	if len(q.GetColumns()) > 0 {
		var result []string
		for _, column := range q.GetColumns() {
			switch v := column.(type) {
			case dto.ModelField:
				result = append(result, fmt.Sprintf("ADD COLUMN %s", generateColumnStr(v)))
			}
		}

		if len(result) > 0 {
			queryStr += fmt.Sprintf("\n%s", strings.Join(result, ","))
		}
	}

	return queryStr
}

//prepareDropQuery method prepares the drop query statement
func prepareDropQuery(q QueryInterface) string {
	return fmt.Sprintf("DROP TABLE %s", q.GetDestination().GetTableName())
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

func getColumnTypeByValue(value interface{}) string {
	switch value.(type)  {
	case int:
		return "INTEGER"
	case int64:
		return "INTEGER"
	case string:
		return "VARCHAR"
	case bool:
		return "BOOLEAN"
	}

	return "VARCHAR"
}