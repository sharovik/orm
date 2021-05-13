# GoLang query builder
Simple query builder, for SQL-like databases

## How to use?
1. Init the client
```go
client, err := SQLiteClient{}.Connect(clients.DatabaseConfig{
  Host:     "mysqlite-database.sqlite",
})

//Create model for needle table
var model = dto.BaseModel{
    TableName: "test_table_name"
}

//Get results
results, err := clients.SQLiteClient{}.Execute(new(Query).Select([]interface{}{"col1", "col2"}).From(&model))
```

### More examples
Please see the `example_mysql.go`