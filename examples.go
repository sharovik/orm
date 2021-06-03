package main

import (
	"fmt"
	"github.com/sharovik/orm/clients"
	"github.com/sharovik/orm/dto"
	"github.com/sharovik/orm/query"
)

func main() {
	//We create database configuration for MySQL database
	configuration := clients.DatabaseConfig{
		Host:     "localhost",
		Username: "root",
		Password: "secret",
		Database: "test",
		Type:     clients.DatabaseTypeMySQL,
		Port:     0,
	}

	client, err := clients.InitClient(configuration)
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
		IsUnsigned:    true,
		AutoIncrement: true,
	})
	model.AddModelField(dto.ModelField{
		Name:          "another_id",
		Type:          dto.IntegerColumnType,
		Value:         nil,
		Default:       nil,
		Length:        0,
		IsNullable:    true,
		IsPrimaryKey:  false,
		IsUnsigned:    true,
		AutoIncrement: false,
	})
	model.AddModelField(dto.ModelField{
		Name:          "test_field",
		Type:          dto.VarcharColumnType,
		Value:         "something",
		Default:       "",
		Length:        255,
		IsNullable:    true,
		IsPrimaryKey:  false,
		AutoIncrement: false,
	})
	model.AddModelField(dto.ModelField{
		Name:          "test_field2",
		Type:          dto.VarcharColumnType,
		Value:         "something",
		Default:       "",
		Length:        255,
		IsNullable:    true,
		IsPrimaryKey:  false,
		AutoIncrement: false,
	})

	//Let's create a table for that model
	q := new(clients.Query).Create(model).
		AddForeignKey(dto.ForeignKey{
			Name: "another_key",
			Target: query.Reference{
				Table: "another",
				Key:   "id",
			},
			With: query.Reference{
				Table: model.TableName,
				Key:   "another_id",
			},
			OnDelete: dto.CascadeAction,
			OnUpdate: dto.NoActionAction,
		}).AddIndex(dto.Index{
		Name:   "some_test_non_unique_index",
		Target: model.TableName,
		Key:    "test_field2",
		Unique: false,
	}).AddIndex(dto.Index{
		Name:   "some_test_unique_index",
		Target: model.TableName,
		Key:    "test_field",
		Unique: true,
	})
	out := client.ToSql(q)
	fmt.Println(out)
	res, err := client.Execute(q)

	//We select specific columns from the table
	var columns = []interface{}{"id", "another_id"}
	q = new(clients.Query).Select(columns).From(model)
	res, err = client.Execute(q)
	fmt.Println(err)
	fmt.Println(res)

	//We do select with join to other table
	q = new(clients.Query).
		Select(columns).
		From(model).
		Join(query.Join{
			Target:    query.Reference{Table: "another", Key: "id"},
			With:      query.Reference{Table: model.TableName, Key: "another_id"},
			Condition: "=",
			Type:      query.LeftJoinType,
		}).
		Where(query.Where{
			First:    "another.id",
			Operator: "is",
			Second:   "NULL",
		})
	res, err = client.Execute(q)
	fmt.Println(err)
	fmt.Println(res)

	//We insert new item into our table
	q = new(clients.Query).Insert(model)
	res, err = client.Execute(q)
	fmt.Println(err)
	fmt.Println(res)

	//We select all model columns from our table
	q = new(clients.Query).Select(model.GetColumns()).From(model)
	res, err = client.Execute(q)
	fmt.Println(err)
	fmt.Println(res)

	//We do select from the table where id = 1 OR id = 2
	q = new(clients.Query).Select(model.GetColumns()).
		From(model).
		Where(query.Where{
			First:    "id",
			Operator: "=",
			Second:   "1",
		}).
		Where(query.Where{
			First:    "id",
			Operator: "=",
			Second:   "2",
			Type:     query.WhereOrType, //For OR condition, you can use the Type attribute of Where object
		})
	res, err = client.Execute(q)
	fmt.Println(err)
	fmt.Println(res)

	//We do select with the more complex WHERE clause
	//The output will be:
	//SELECT id, id, id, another_id, test_field, test_field2 FROM test_table_name WHERE (id = 1 OR id = 2) OR id = 3
	q = new(clients.Query).Select(model.GetColumns()).
		From(model).
		Where(query.Where{
			First:    query.Where{
				First:    "id",
				Operator: "=",
				Second:   "1",
			},
			Operator: "",
			Second:   query.Where{
				First:    "id",
				Operator: "=",
				Second:   "2",
				Type:     query.WhereOrType, //For OR condition, you can use the Type attribute of Where object
			},
		}).
		Where(query.Where{
			First:    "id",
			Operator: "=",
			Second:   "3",
			Type:     query.WhereOrType, //For OR condition, you can use the Type attribute of Where object
		})
	res, err = client.Execute(q)
	fmt.Println(err)
	fmt.Println(res)

	//We update the table
	model.SetField("test_field2", "test test test")
	q = new(clients.Query).Update(model)
	res, err = client.Execute(q)
	fmt.Println(err)
	fmt.Println(res)

	//We alter table
	q = new(clients.Query).Alter(model).
		AddColumn(dto.ModelField{
			Name:          "new_column",
			Type:          dto.VarcharColumnType,
			Value:         nil,
			Default:       "",
			Length:        244,
			IsNullable:    true,
			IsPrimaryKey:  false,
			IsUnsigned:    false,
			AutoIncrement: false,
		}).DropColumn(dto.ModelField{
		Name: "test_field",
	}).DropForeignKey(dto.ForeignKey{
		Name: "another_key",
	}).DropIndex(dto.Index{
		Name: "some_test_unique_index",
	})
	res, err = client.Execute(q)
	fmt.Println(err)
	fmt.Println(res)

	//We delete item from the table
	q = new(clients.Query).Delete().
		From(model).
		Where(query.Where{
			First:    "id",
			Operator: "=",
			Second:   "1",
		})
	res, err = client.Execute(q)
	fmt.Println(err)
	fmt.Println(res)

	//We drop the table
	q = new(clients.Query).Drop(model)
	res, err = client.Execute(q)
	fmt.Println(err)
	fmt.Println(res)
}
