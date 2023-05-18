package clients

import (
	"database/sql"
	"errors"

	"github.com/sharovik/orm/dto"
	"github.com/sharovik/orm/query"
)

const (
	CreateType = "CREATE"
	AlterType  = "ALTER"
	RenameType = "RENAME"
	DropType   = "DROP"
	SelectType = "SELECT"
	InsertType = "INSERT"
	UpdateType = "UPDATE"
	DeleteType = "DELETE"

	DatabaseTypeMySQL   = "mysql"
	DatabaseTypeSqlite  = "sqlite"
	DefaultDatabaseType = DatabaseTypeSqlite
)

// DatabaseConfig the config which will be used by the client
type DatabaseConfig struct {
	Host     string
	Database string
	Username string
	Password string
	Port     int64
	Engine   string
	Charset  string
	Collate  string
	Type     string
}

func (c DatabaseConfig) GetType() string {
	var result = DefaultDatabaseType
	if c.Type != "" {
		result = c.Type
	}

	return result
}

func (c DatabaseConfig) GetEngine() string {
	return c.Engine
}

func (c DatabaseConfig) GetCharset() string {
	return c.Charset
}

func (c DatabaseConfig) GetCollate() string {
	return c.Collate
}

// BaseClientInterface the main interface for the client
type BaseClientInterface interface {
	Connect(config DatabaseConfig) (client BaseClientInterface, err error)
	Disconnect() error
	GetClient() *sql.DB
	ToSql(query QueryInterface) string
	Execute(query QueryInterface) (result dto.BaseResult, err error)
}

// QueryInterface the interface for the query builder of the client
type QueryInterface interface {
	//Create method will return the query object for table creation
	Create(dto.ModelInterface) QueryInterface

	//Drop method will return the query object for table drop
	Drop(dto.ModelInterface) QueryInterface

	//Select using that method you can set the attributes for selection. This method should be used from the beginning of your query, to specify the initial query string.
	//This method returns the updated Query object.
	Select(columns interface{}) QueryInterface

	//Insert method should be used when you need to insert something into selected table.
	//It should be used from the beginning of your query, to specify the initial query string.
	//This method receives the dto.ModelInterface object and returns the last inserted ID and error(if it exists)
	Insert(dto.ModelInterface) QueryInterface

	//Update method should be used when you need to update something in your selected table.
	//It should be used from the beginning of your query, to specify the initial query string.
	//This method receives the dto.ModelInterface object and returns the updated QueryInterface object.
	Update(dto.ModelInterface) QueryInterface

	//Alter method should be used when you need to alter selected table.
	//It should be used from the beginning of your query, to specify the initial query string and with combination one of the methods AddColumn, DropColumn, AddIndex, DropIndex, AddForeignKey, DropForeignKey
	//This method receives the dto.ModelInterface object and returns the updated QueryInterface object.
	Alter(dto.ModelInterface) QueryInterface

	//Rename method renames the table to a new table name
	Rename(original string, newTableName string) QueryInterface

	//Delete method deletes the row.
	Delete() QueryInterface

	//GetQueryType Returns the type of the Query. Eg: INSERT, ALTER, SELECT, DELETE, CREATE, UPDATE, DROP
	GetQueryType() string
	GetNewTableName() string
	GetDestination() dto.ModelInterface
	GetColumns() []interface{}
	GetColumnsToDrop() []interface{}
	GetForeignKeysToAdd() []dto.ForeignKey
	GetForeignKeysToDrop() []dto.ForeignKey
	GetIndexesToAdd() []dto.Index
	GetIndexesToDrop() []dto.Index
	GetWheres() []query.Where
	GetJoins() []query.Join
	GetOrderBy() []query.OrderByColumn
	GetGroupBy() []string
	GetLimit() query.Limit

	Values(interface{}) QueryInterface
	GetValues() interface{}

	//From using this method you can specify the ORDER BY fields with the right direction to order.
	From(model interface{}) QueryInterface

	//IfNotExists Sets the IfNotExists flag. Method can be used in the combination with CREATE TABLE statement to have condition CREATE TABLE IF NOT EXISTS
	IfNotExists() QueryInterface

	//GetIfNotExists method can be used in the combination with CREATE TABLE statement to have condition CREATE TABLE IF NOT EXISTS
	GetIfNotExists() bool

	//OrderBy using this method you can specify the ORDER BY fields with the right direction to order.
	OrderBy(field string, direction string) QueryInterface

	Limit(limit query.Limit) QueryInterface

	//GroupBy using this method you can specify the fields for the GROUP BY clause.
	GroupBy(field string) QueryInterface

	//Where method needed for WHERE clause configuration.
	Where(where query.Where) QueryInterface

	//Join method can be used for specification of JOIN clause.
	Join(join query.Join) QueryInterface

	//AddColumn the method which identifies which field we need to add for selected model in Alter method
	AddColumn(column dto.ModelField) QueryInterface

	AddBinding(query.Bind) QueryInterface

	//DropColumn the method which identifies which field we need to drop for selected model in Alter method
	DropColumn(dto.ModelField) QueryInterface

	//AddForeignKey the method which identifies which foreign key we need to add for selected model in Alter method
	AddForeignKey(field dto.ForeignKey) QueryInterface

	//DropForeignKey the method which identifies which foreign key we need to drop for selected model in Alter method
	DropForeignKey(field dto.ForeignKey) QueryInterface

	//AddIndex the method which identifies which index key we need to add for selected model in Alter method
	AddIndex(index dto.Index) QueryInterface

	//DropIndex the method which identifies which index key we need to add for selected model in Alter method
	DropIndex(index dto.Index) QueryInterface

	//GetBindings method returns the binding collected during the query building
	GetBindings() []query.Bind
}

// Query the query object of the SQLite client
type Query struct {
	destination     dto.ModelInterface
	bindings        []query.Bind
	queryType       string
	newTableName    string
	columns         []interface{}
	columnsDrop     []interface{}
	ifNotExists     bool
	indexAdd        []dto.Index
	indexDrop       []dto.Index
	foreignKeysAdd  []dto.ForeignKey
	foreignKeysDrop []dto.ForeignKey
	wheres          []query.Where
	joins           []query.Join
	orderBys        []query.OrderByColumn
	groupBys        []string
	values          interface{}
	limit           query.Limit
}

func (q *Query) GetQueryType() string {
	return q.queryType
}

func (q *Query) GetNewTableName() string {
	return q.newTableName
}

func (q *Query) GetDestination() dto.ModelInterface {
	return q.destination
}

func (q *Query) GetColumns() []interface{} {
	return q.columns
}

func (q *Query) GetColumnsToDrop() []interface{} {
	return q.columnsDrop
}

func (q *Query) GetForeignKeysToAdd() []dto.ForeignKey {
	return q.foreignKeysAdd
}

func (q *Query) GetForeignKeysToDrop() []dto.ForeignKey {
	return q.foreignKeysDrop
}

func (q *Query) GetIndexesToAdd() []dto.Index {
	return q.indexAdd
}

func (q *Query) GetIndexesToDrop() []dto.Index {
	return q.indexDrop
}

func (q *Query) GetWheres() []query.Where {
	return q.wheres
}

func (q *Query) GetJoins() []query.Join {
	return q.joins
}

func (q *Query) GetOrderBy() []query.OrderByColumn {
	return q.orderBys
}

func (q *Query) GetGroupBy() []string {
	return q.groupBys
}

func (q *Query) GetBindings() []query.Bind {
	return q.bindings
}

func (q *Query) GetLimit() query.Limit {
	return q.limit
}

// From using this method you can specify the ORDER BY fields with the right direction to order.
func (q *Query) From(model interface{}) QueryInterface {
	switch v := model.(type) {
	case dto.ModelInterface:
		q.destination = v
	case string:
		q.destination = &dto.BaseModel{
			TableName: v,
		}
	}

	return q
}

// IfNotExists sets the ifNotExists flag. method can be used in the combination with CREATE TABLE statement to have condition CREATE TABLE IF NOT EXISTS
func (q *Query) IfNotExists() QueryInterface {
	q.ifNotExists = true
	return q
}

// GetIfNotExists method can be used in the combination with CREATE TABLE statement to have condition CREATE TABLE IF NOT EXISTS
func (q *Query) GetIfNotExists() bool {
	return q.ifNotExists
}

// OrderBy using this method you can specify the ORDER BY fields with the right direction to order.
func (q *Query) OrderBy(field string, direction string) QueryInterface {
	q.orderBys = append(q.orderBys, query.OrderByColumn{
		Direction: direction,
		Column:    field,
	})
	return q
}

// Limit using this method you can set the limitation for result of your query
func (q *Query) Limit(limit query.Limit) QueryInterface {
	q.limit = limit
	return q
}

// GroupBy using this method you can specify the fields for the GROUP BY clause.
func (q *Query) GroupBy(field string) QueryInterface {
	q.groupBys = append(q.groupBys, field)
	return q
}

// Where method needed for WHERE clause configuration.
func (q *Query) Where(where query.Where) QueryInterface {
	switch v := where.First.(type) {
	case query.Bind:
		q.AddBinding(v)
		where.First = "?"
	}

	switch v := where.Second.(type) {
	case query.Bind:
		q.AddBinding(v)
		where.Second = "?"
	}

	q.wheres = append(q.wheres, where)
	return q
}

// Join method can be used for specification of JOIN clause.
func (q *Query) Join(join query.Join) QueryInterface {
	q.joins = append(q.joins, join)
	return q
}

// AddColumn the method which identifies which field we need to add for selected model in Alter method
func (q *Query) AddColumn(column dto.ModelField) QueryInterface {
	q.columns = append(q.columns, column)
	return q
}

// DropColumn the method which identifies which field we need to drop for selected model in Alter method
func (q *Query) DropColumn(column dto.ModelField) QueryInterface {
	q.columnsDrop = append(q.columnsDrop, column)
	return q
}

// AddForeignKey the method which identifies which foreign key we need to add for selected model in Alter method
func (q *Query) AddForeignKey(field dto.ForeignKey) QueryInterface {
	q.foreignKeysAdd = append(q.foreignKeysAdd, field)
	return q
}

// AddBinding the method which identifies which foreign key we need to add for selected model in Alter method
func (q *Query) AddBinding(field query.Bind) QueryInterface {
	q.bindings = append(q.bindings, field)
	return q
}

// DropForeignKey the method which identifies which foreign key we need to drop for selected model in Alter method
func (q *Query) DropForeignKey(field dto.ForeignKey) QueryInterface {
	q.foreignKeysDrop = append(q.foreignKeysDrop, field)
	return q
}

// AddIndex the method which identifies which index key we need to add for selected model in Alter method
func (q *Query) AddIndex(field dto.Index) QueryInterface {
	q.indexAdd = append(q.indexAdd, field)
	return q
}

// DropIndex the method which identifies which index key we need to add for selected model in Alter method
func (q *Query) DropIndex(field dto.Index) QueryInterface {
	q.indexDrop = append(q.indexDrop, field)
	return q
}

// Values sets the values, which will be used for insert queries
func (q *Query) Values(values interface{}) QueryInterface {
	q.values = values
	return q
}

// GetValues retrieves the values added by Values method
func (q *Query) GetValues() interface{} {
	return q.values
}

// Create method will return the query object for table creation
func (q *Query) Create(model dto.ModelInterface) QueryInterface {
	q.queryType = CreateType
	q.From(model)
	return q
}

// Drop method will return the query object for table drop
func (q *Query) Drop(model dto.ModelInterface) QueryInterface {
	q.queryType = DropType
	q.From(model)
	return q
}

// Rename will rename the table to the new table name
func (q *Query) Rename(table string, newTableName string) QueryInterface {
	q.queryType = RenameType
	q.From(&dto.BaseModel{
		TableName: table,
	})
	q.newTableName = newTableName
	return q
}

// Select using that method you can set the attributes for selection. This method should be used from the beginning of your query, to specify the initial query string.
// This method returns the updated Query object.
func (q *Query) Select(columns interface{}) QueryInterface {
	q.queryType = SelectType
	switch c := columns.(type) {
	case []interface{}:
		for _, field := range c {
			switch v := field.(type) {
			case string:
				q.AddColumn(dto.ModelField{
					Name:  v,
					Type:  "",
					Value: nil,
				})
			case dto.ModelField:
				q.AddColumn(v)
			}
		}
	case string:
		q.AddColumn(dto.ModelField{
			Name:  c,
			Type:  "",
			Value: nil,
		})
	}

	return q
}

// Insert method should be used when you need to insert something into selected table.
// It should be used from the beginning of your query, to specify the initial query string.
// This method receives the dto.ModelInterface object and returns the last inserted ID and error(if it exists)
func (q *Query) Insert(model dto.ModelInterface) QueryInterface {
	q.queryType = InsertType
	for _, field := range model.GetColumns() {
		switch v := field.(type) {
		case dto.ModelField:
			q.AddColumn(v)
			//For primary keys we don't
			if v.Value == nil && v.IsPrimaryKey {
				break
			}

			q.AddBinding(query.Bind{
				Field: "?",
				Value: v.Value,
			})
		}
	}

	q.destination = model
	return q
}

// Update method should be used when you need to update something in your selected table.
// It should be used from the beginning of your query, to specify the initial query string.
// This method receives the dto.ModelInterface object and returns the updated query.UpdateQuery object.
func (q *Query) Update(model dto.ModelInterface) QueryInterface {
	q.queryType = UpdateType
	for _, field := range model.GetColumns() {
		switch v := field.(type) {
		case dto.ModelField:
			if v.IsPrimaryKey {
				continue
			}

			q.AddColumn(v)
			q.AddBinding(query.Bind{
				Field: "?",
				Value: v.Value,
			})
		}
	}

	q.destination = model

	return q
}

// Alter method should be used when you need to alter selected table.
// It should be used from the beginning of your query, to specify the initial query string and with combination one of the methods AddColumn, DropColumn, AddIndex, DropIndex, AddForeignKey, DropForeignKey
// This method receives the dto.ModelInterface object and returns the updated query.AlterQuery object.
func (q *Query) Alter(model dto.ModelInterface) QueryInterface {
	q.queryType = AlterType
	q.From(model)
	return q
}

// Delete method deletes the row.
func (q *Query) Delete() QueryInterface {
	q.queryType = DeleteType
	return q
}

// InitClient method can be used for the database client init
func InitClient(config DatabaseConfig) (BaseClientInterface, error) {
	switch config.GetType() {
	case DatabaseTypeSqlite:
		return SQLiteClient{}.Connect(config)
	case DatabaseTypeMySQL:
		return MySQLClient{}.Connect(config)
	}

	return nil, errors.New("failed to init the database client. ")
}
