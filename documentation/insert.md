# Insert queries
Here you can find information about the insert queries, the way how to use the Insert query statement, etc.

## Insert a model
Like in all other ORM's you can use the model as the source of data for insert queries.
```go
type TestTableModel struct {
    dto.BaseModel
}

model := new(TestTableModel)
model.SetTableName("test_table_name")
model.SetPrimaryKey(dto.ModelField{
    Name:          "id",
    Type:          "integer",
    Value:         nil,
    Default:       nil,
    Length:        0,
    IsNullable:    false,
    IsUnsigned:    true,
    AutoIncrement: true,
})
model.AddModelField(dto.ModelField{
    Name:          "col1",
    Type:          dto.VarcharColumnType,
    Value:         "something",
    Default:       "",
    Length:        255,
    IsNullable:    true,
    IsPrimaryKey:  false,
    AutoIncrement: false,
})
model.AddModelField(dto.ModelField{
    Name:          "col2",
    Type:          dto.VarcharColumnType,
    Value:         "something",
    Default:       "",
    Length:        255,
    IsNullable:    true,
    IsPrimaryKey:  false,
    AutoIncrement: false,
})
model.AddModelField(dto.ModelField{
    Name:          "col3",
    Type:          dto.VarcharColumnType,
    Value:         "something",
    Default:       "",
    Length:        255,
    IsNullable:    true,
    IsPrimaryKey:  false,
    AutoIncrement: false,
})

//We insert new item into our table
q := new(clients.Query).Insert(model)
res, err := client.Execute(q)
if err != nil {
    panic(err)
}
```
That structure will generate the next query
```sql
INSERT INTO test_table_name (col1, col2, col3) VALUES (?, ?, ?)
```
## Insert-select
You can use as the values for your insert queries the output of other select statement.
```go
q := new(Query).Insert(&model).Values(new(Query).Select([]interface{}{}).From(&dto.BaseModel{
    TableName: "other_table_name",
}))
res, err := client.Execute(q)
fmt.Println(err)
fmt.Println(res)
if err != nil {
    panic(err)
}
```
As you can see, here as VALUES we are using another query statement. In example, you see the basic example of insert-select query, the output will be like:
```sql
INSERT INTO test_table_name (relation_id, col1, col2, col3) SELECT * FROM other_table_name
```