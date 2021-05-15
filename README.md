# GoLang query builder
Simple query builder, for SQL-like databases

## How to use?
1. Init the client
```go

//For sqlite
client, err := SQLiteClient{}.Connect(clients.DatabaseConfig{
  Host:     "mysqlite-database.sqlite",
})

//For mysql
client, err := clients.MySQLClient{}.Connect(clients.DatabaseConfig{
    Host:     "localhost",
    Username: "root",
    Password: "secret",
    Database: "test",
    Port:     0,
})

//Create model for needle table
var model = dto.BaseModel{
    TableName: "test_table_name"
}

//Get results sqlite
results, err := clients.SQLiteClient{}.Execute(new(Query).Select([]interface{}{"col1", "col2"}).From(&model))

//Get results mysql
results, err := clients.MySQLClient{}.Execute(new(Query).Select([]interface{}{"col1", "col2"}).From(&model))
```

### More examples
Please see the `examples.go`