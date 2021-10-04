# SQLite warnings
Here you will find the information about the tricks we used to handle scenarios for SQLite client to update, drop or add the columns, foreign keys.

## Drop foreign keys, columns
Currently, SQLite does not support that. To trigger that action you need to do the next:
1. create temp table with the new schema and all indexes, all foreign keys
2. import data from the old schema into the temp table
3. drop old table
4. rename temp table to the old table name

To make your life easier, using the SQLite client you would need just to trigger simple alter query builder with the definition of the foreign keys and indexes which should exist.

```go
SQLiteClient{}.ToSql(new(Query).Alter(&model).
//You must define again all the foreign keys and indexes
AddForeignKey(dto.ForeignKey{
    Name: "fk_test",
    Target: query.Reference{
        Table: "test_table_name2",
        Key:   "id",
    },
    With: query.Reference{
        Table: "test_table_name",
        Key:   "relation_id",
    },
    OnDelete: "",
    OnUpdate: "",
}).
DropColumn(dto.ModelField{
    Name: "col3",
}))
```
That structure will generate the next SQLite query snippet example:
```sql
CREATE TABLE temp_test_table_name (id INTEGER CONSTRAINT temp_test_table_name_pk primary key autoincrement, relation_id INTEGER NOT NULL, col1 INTEGER NOT NULL, col2 INTEGER NOT NULL, col3 VARCHAR NOT NULL,
CONSTRAINT fk_test
FOREIGN KEY (relation_id)
 REFERENCES test_table_name2 (id)
ON DELETE NO ACTION
ON UPDATE NO ACTION);
INSERT INTO temp_test_table_name (relation_id, col1, col2, col3) SELECT relation_id, col1, col2, col3 FROM test_table_name;
ALTER TABLE `test_table_name` RENAME TO `old_test_table_name`;
ALTER TABLE `temp_test_table_name` RENAME TO `test_table_name`;
DROP TABLE old_test_table_name;
```