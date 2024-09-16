# GoLang query builder
Simple query builder, for SQL-like databases.

## Databases supported
- MySQL
- SQLite

## How to use?

### Install the package
```go
go get -u github.com/sharovik/orm
```

### Init the client
```go
import (
    "github.com/sharovik/orm/clients"
)

//For sqlite
databaseClient, err := clients.InitClient(clients.DatabaseConfig{
  Host:     "database.sqlite",
})

//For mysql
databaseClient, err := clients.InitClient(clients.DatabaseConfig{
    Host:     "localhost",
    Username: "root",
    Password: "secret",
    Database: "test",
    Type:     clients.DatabaseTypeMySQL,
})

```
### Start using the query builder
There are several ways, how you can communicate with database using this query builder

**Basic example**
```go
//Get results sqlite
//We select specific columns from the table
var columns = []interface{}{"col1", "col2"}
q := new(clients.Query).Select(columns).From("test_table_name")
_, err := databaseClient.Execute(q)
if err != nil {
    return err
}
```

**Using model**
```go
package main

import "github.com/sharovik/orm/dto"

type TestTableModel struct {
	dto.BaseModel
}

func main() {
	model := new(TestTableModel)
	model.SetTableName("test_table_name")
	model.SetPrimaryKey(dto.ModelField{
		Name:          "col1",
		Type:          "integer",
		Value:         nil,
		Default:       nil,
		Length:        0,
		IsNullable:    false,
		IsUnsigned:    true,
		AutoIncrement: true,
	})
	model.AddModelField(dto.ModelField{
		Name:          "col2",
		Type:          dto.IntegerColumnType,
		Value:         nil,
		Default:       nil,
		Length:        0,
		IsNullable:    true,
		IsPrimaryKey:  false,
		IsUnsigned:    true,
		AutoIncrement: false,
	})

	//Get results sqlite
	results, err := databaseClient.Execute(new(Query).Select([]interface{}{"col1", "col2"}).From(&model))
    
    //do something with the results...
}
```
This code will execute the next SQL query
```sql
SELECT col1, col2 FROM test_table_name
```

### Complex queries
You can also build more complex queries, like:
```sql
SELECT id, another_id FROM test_table_name 
LEFT JOIN another ON another.id = test_table_name.another_id
WHERE another.id is NULL
```
This how it will look in code:
```go
var columns = []interface{}{"id", "another_id"}
q = new(clients.Query).
    Select(columns).
    From(&model).
    Join(query.Join{
        Target:    query.Reference{Table: "another", Key: "id"},
        With:      query.Reference{Table: model.TableName, Key: "another_id"},
        Condition: "=",
        Type:      query.LeftJoinType,
    }).
    Where(query.Where{
        First:    "another.id",
        Operator: "is",
        Second:   "NULL",
    })
result, err := client.Execute(q)
```

OR even queries with the complex WHERE CLAUSE, like this:
```sql
SELECT id, another_id, test_field, test_field2 FROM test_table_name WHERE (id = 1 OR id = 2) OR id = 3
```
In code this will look like:
```go
q = new(clients.Query).Select(&model.GetColumns()).
    From(&model).
    Where(query.Where{
        First:    query.Where{
            First:    "id",
            Operator: "=",
            Second:   "1",
        },
        Operator: "",
        Second:   query.Where{
            First:    "id",
            Operator: "=",
            Second:   "2",
            Type:     query.WhereOrType, //For OR condition, you can use the Type attribute of Where object
        },
    }).
    Where(query.Where{
        First:    "id",
        Operator: "=",
        Second:   "3",
        Type:     query.WhereOrType, //For OR condition, you can use the Type attribute of Where object
    })
res, err = client.Execute(q)
```

And for sure you always can use the query **binding** for the untrusted data
```go
q = new(clients.Query).
    Select(model.GetColumns()).
    From(&model).
    Where(query.Where{
        First:    "name",
        Operator: "=",
        Second:   query.Bind{
            Field: "name",
            Value: "my test name",
        },
    })
result, err := client.Execute(q)
```
This will generate the next prepared query
```sql
SELECT id, name FROM test_table_name WHERE name = ?
```

### Work with the results
As the output of the `Execute` method you will receive the object type of `dto.BaseResult`. This object can contain the list of items for your query, the database error OR the last inserted ID.

Please see the [examples.go](examples.go) file for more queries examples. And also, please check the [documentation files here](documentation).

### Other notes
- [table renaming](documentation/rename-table.md)
- [CREATE TABLE statement](documentation/create-tables.md)
- [Insert queries](documentation/insert.md)
- [Transactions](documentation/transactions.md)
- [Models](documentation/model.md)
- [SQLite warnings](documentation/sqlite-warnings.md)
