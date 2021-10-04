# Table renaming
For a table renaming you will need to use the next structure for the query
```go
q := new(clients.Query).Rename("old_table_name", "new_table_name")
res, err = client.Execute(q)
if err != nil {
    panic(err)
    return
}
```

This structure will generate the next SQL query snippet
```sql
ALTER TABLE `old_table_name` RENAME TO `new_table_name`
```