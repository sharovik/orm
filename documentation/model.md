# Models
This ORM does support models usage in the query requests

## How to init
Below you can find the example of how you can init your model, which will represent the table in your database
```go
model := dto.BaseModel{
    TableName: "my_table_name",
    Fields: []interface{}{
        dto.ModelField{
            Name:  "relation_id",
            Type:  dto.IntegerColumnType,
            Value: int64(1),
        },
        dto.ModelField{
            Name:  "col3",
            Type:  dto.VarcharColumnType,
            Value: "Test",
        },
    },
    PrimaryKey: dto.ModelField{
        Name:          "id",
        Type:          dto.IntegerColumnType,
        AutoIncrement: true,
    },
}
```

Each model should content the fields, the table name and the primary key.

Each field should be described as type of `dto.ModelField`. Same type should be for primary key field. The table name, should be type of string.