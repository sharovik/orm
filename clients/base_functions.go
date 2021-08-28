package clients

import (
	"database/sql"
	"fmt"
	"github.com/sharovik/orm/dto"
	"github.com/sharovik/orm/query"
	"strconv"
	"strings"
)

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
	var resultStr = "WHERE "

	for i, where := range wheres {
		//If we have the type of WHERE clause specified and this is not first element, we do set the type.
		if where.GetType() != "" && i != 0 {
			resultStr += fmt.Sprintf(" %s ", where.GetType())
		}

		resultStr += whereToStr(where)
	}

	return strings.TrimSpace(resultStr)
}

func whereToStr(where query.Where) string {
	var (
		resultStr string
		isFirstIsWhere, isSecondIsWhere bool
		)
	switch w := where.First.(type) {
	case query.Where:
		resultStr += fmt.Sprintf("%s", whereToStr(w))
		isFirstIsWhere = true
		break
	default:
		resultStr += fmt.Sprintf("%s", where.First)
		resultStr += fmt.Sprintf(" %s ", where.Operator)
	}

	switch w := where.Second.(type) {
	case query.Where:
		resultStr += fmt.Sprintf(" %s %s", w.GetType(), whereToStr(w))
		isSecondIsWhere = true
		break
	default:
		resultStr += fmt.Sprintf("%s", where.Second)
	}

	if isFirstIsWhere && isSecondIsWhere {
		resultStr = fmt.Sprintf("(%s)", resultStr)
	}

	return resultStr
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

	if column.IsUnsigned {
		resultStr += fmt.Sprintf(" unsigned")
	}

	if column.Default != nil {
		resultStr += fmt.Sprintf(" DEFAULT %s", toSQLValue(column.Default))
	}

	if column.IsNullable {
		resultStr += fmt.Sprintf(" %s", "NULL")
	} else {
		resultStr += fmt.Sprintf(" %s", "NOT NULL")
	}

	if column.AutoIncrement {
		resultStr += " autoincrement"
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
	var queryStr = ""

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

	//@todo: drop columns
	//To properly drop columns in the sqlite, we need to create a new temp table with the all columns, except that column.
	//Import data from the old table into the new one
	//Rename the old one to the new name and rename the new table to the correct table name
	//Drop the old table

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

	//@todo: generate foreign keys to add
	//The same scenario as for columns drop

	//@todo: generate foreign keys to drop
	//The same scenario as for columns drop

	if len(result) > 0 {
		queryStr += fmt.Sprintf("%s", strings.Join(result, ";\n"))
	}

	return queryStr
}

//prepareDropQuery method prepares the drop query statement
func prepareDropQuery(q QueryInterface) string {
	return fmt.Sprintf("DROP TABLE %s", q.GetDestination().GetTableName())
}

func prepareColumnTypes(rows *sql.Rows) (result []string, err error) {
	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}

	result = make([]string, len(columnTypes))
	for i, columnType := range columnTypes {
		result[i] = normalizeColumnType(columnType.DatabaseTypeName())
	}

	return result, nil
}

func normalizeValue(value interface{}, columnType string) interface{} {
	switch v := value.(type) {
	case int64:
		return int(v)
	case []uint8:
		switch columnType {
		case dto.IntegerColumnType:
			res, _ := strconv.Atoi(string(v))
			return res
		case dto.VarcharColumnType:
			return string(v)
		case dto.CharColumnType:
			return string(v)
		case dto.BooleanColumnType:
			return string(v) == "1" || string(v) == "true"
		}
	default:
		return v
	}

	return value
}

func normalizeColumnType(columnType string) string {
	switch columnType  {
	case "INT":
		return dto.IntegerColumnType
	case "INTEGER":
		return dto.IntegerColumnType
	case "VARCHAR":
		return dto.VarcharColumnType
	case "CHAR":
		return dto.CharColumnType
	case "BOOL":
		return dto.BooleanColumnType
	case "BOOLEAN":
		return dto.BooleanColumnType
	}

	return dto.VarcharColumnType
}

func toSQLValue(value interface{}) string {
	var resultStr string
	if value == nil {
		return "NULL"
	}

	switch v := value.(type) {
	case int:
		resultStr += fmt.Sprintf("%d", v)
		break
	case int64:
		resultStr += fmt.Sprintf("%d", v)
		break
	case string:
		resultStr += fmt.Sprintf(`"%s"`, v)
		break
	case bool:
		resultStr += fmt.Sprintf("%t", v)
		break
	}

	return resultStr
}