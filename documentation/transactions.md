# Transactions
In order to trigger transaction queries you might need to use `BeginTransaction`, `CommitTransaction` or `RollbackTransaction` methods.

```go
//You trigger begin of transaction
q := new(clients.Query).BeginTransaction()
_, err := client.Execute(q)
if err != nil {
    return err
}

//You run your queries
//....

if err != nil {
	//you handle errors and rollback if needed
	q = new(clients.Query).RollbackTransaction()
    _, err = client.Execute(q)
    if err != nil {
        return err
    }
}

//You commit the changes
q = new(clients.Query).CommitTransaction()
_, err = client.Execute(q)
if err != nil {
    return err
}
```