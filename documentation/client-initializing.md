# How to initialize the client
There are several steps which you need to follow, for proper sql client initialisation.

## Database
Before all, you need to make sure you configured the database configuration.
To do so, you need to initialize first the configuration for the database connection.

```go
var dCfg = clients.DatabaseConfig{
    Host:     testSQLiteDatabasePath,
    Username: "__USERNAME__",//in case if you are using sqlite, please leave that empty
    Password: "__PASSWORD__",//in case if you are using sqlite, please leave that empty
    Port:     0,
}
```
## Create a client
Once you have a database configuration, you can initialise the client
```go
client, err := SQLiteClient{}.Connect(dCfg)
```

Now the client is ready for your first sql query!
