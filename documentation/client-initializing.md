# How to initialize the client
There are several steps which you need to follow, for proper sql client initialisation.

## Database configuration
In example below you can see the configuration for MySQL database connection
```go
var configuration := clients.DatabaseConfig{
    Host:     "localhost",
    Username: "root",
    Password: "secret",
    Database: "test",
    Type:     clients.DatabaseTypeMySQL,
}
```
And here is for sqlite
```go
var configuration := clients.DatabaseConfig{
    Host:     "testing.sqlite",
}
```
## Create a client
Once you have a database configuration, you can initialise the client
### Sqlite
```go
client, err := SQLiteClient{}.Connect(configuration)
```
### MySQL
```go
client, err := MySQLClient{}.Connect(configuration)
```

Now the client is ready for your first sql query!
