package main

import (
	"fmt"
	"github.com/sharovik/orm/clients"
	"github.com/sharovik/orm/dto"
)

func main() {
	sqliteClient, err := clients.SQLiteClient{}.Connect(clients.DatabaseConfig{
		Host:     "test.sqlite",
		Username: "",
		Password: "",
		Port:     0,
	})
	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to connect to the database. Reason: %s", err))
		return
	}

	model := new(dto.BaseModel)
	model.SetTableName("test_table_name")
	model.SetPrimaryKey(dto.ModelField{
		Name:          "id",
		Type:          "integer",
		Value:         nil,
		Default:       nil,
		Length:        0,
		IsNullable:    false,
		AutoIncrement: true,
	})

	var columns = []interface{}{"id", "relation_id"}
	query := new(clients.Query).Select(columns).From(model)
	res, err := sqliteClient.Execute(query)
	fmt.Println(err)
	fmt.Println(res)
}