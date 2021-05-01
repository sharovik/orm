# Select queries
To do a select query you can use the following instructions. Before we start playing with the query builder, please [make sure you prepare the client](client-initializing.md).

## How it works
First you initialize the model object which will represent the table in your database. [You can find how to do it here](model.md).

After you created the model, you can start using the query builder to build your query.

## The examples
Below you will find the examples of how to use the query builder and what will be the output.

### Base query
```go
imports (
 "github.com/sharovik/orm/clients"
)

var model = dto.BaseModel{
    TableName: "test_table_name"
}

results, err := clients.SQLiteClient{}.Execute(new(Query).Select([]interface{}{"col1", "col2"}).From(&model))
```
This part of code will generate the sql query
```go
SELECT col1, col2 FROM test_table_name
```
and execute it. Then it will return the output results and the error, if there was an error during the query execution.