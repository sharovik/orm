# Create tables
With the current version of query builder you are able to create tables in your database using "create" query statement.

## How to
To create a table you are require to use the model. You can define the model object in your project or use just temporary object in your query.

### Using model object
The example with model object:
```go
type MyModel struct {
    dto.BaseModel
}

another := new(MyModel)
another.SetTableName("another")
another.SetPrimaryKey(dto.ModelField{
    Name:          "id",
    Type:          "integer",
    Value:         nil,
    Default:       nil,
    Length:        0,
    IsNullable:    false,
    IsUnsigned:    true,
    AutoIncrement: true,
})

//Let's create a table for that model
q := new(clients.Query).Create(another)
res, err := client.Execute(q)
if err != nil {
    panic(err)
}
```
The output will look like:
```sql
CREATE TABLE another (
    id integer unsigned NOT NULL AUTO_INCREMENT,
    PRIMARY KEY (id)
);
```

You can also use the `IfNotExists` flag to prevent the duplicate create statement execution in case if table already exists.
```go
q := new(clients.Query).Create(another).IfNotExists()
res, err := client.Execute(q)
if err != nil {
    panic(err)
}
```
The output will look like:
```sql
CREATE TABLE IF NOT EXISTS another (
    id integer unsigned NOT NULL AUTO_INCREMENT,
    PRIMARY KEY (id)
);
```