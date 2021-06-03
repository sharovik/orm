# GoLang query builder
Simple query builder, for SQL-like databases.

## How to use?
### Import the ORM package
```go
import (
    "github.com/sharovik/orm/clients"
    cdto "github.com/sharovik/orm/dto" //If you don't have the `dto` package name in your project, then you can remove custom `cdto` alias
)

```
### Init the client
```go

//For sqlite
databaseClient, err := clients.InitClient(clients.DatabaseConfig{
  Host:     "mysqlite-database.sqlite",
})

//For mysql
databaseClient, err := clients.InitClient(clients.DatabaseConfig{
    Host:     "localhost",
    Username: "root",
    Password: "secret",
    Database: "test",
    Type:     clients.DatabaseTypeMySQL,
    Port:     0,
})

```
### Start using the query builder
```go
//Create model for needle table
var model = dto.BaseModel{
    TableName: "test_table_name"
}

//Get results sqlite
results, err := databaseClient.Execute(new(Query).Select([]interface{}{"col1", "col2"}).From(&model))

//Get results mysql
results, err := databaseClient.Execute(new(Query).Select([]interface{}{"col1", "col2"}).From(&model))
```
Which is equal to
```sql
SELECT col1, col2 FROM test_table_name
```

OR you can also build more complex queries, like:
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
    From(model).
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
q = new(clients.Query).Select(model.GetColumns()).
    From(model).
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

### More examples
Please see the [examples.go](examples.go) file. 