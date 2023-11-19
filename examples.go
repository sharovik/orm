package main

import (
	"errors"
	"fmt"

	"github.com/sharovik/orm/clients"
	"github.com/sharovik/orm/dto"
	"github.com/sharovik/orm/query"
)

type AnotherModel struct {
	dto.BaseModel
}

type TestTableModel struct {
	dto.BaseModel
}

var (
	err     error
	client  clients.BaseClientInterface
	model   *TestTableModel
	another *AnotherModel
)

func main() {
	//We create database configuration for MySQL database
	configuration := clients.DatabaseConfig{
		Host:     "localhost",
		Username: "root",
		Password: "secret",
		Database: "test",
		Type:     clients.DatabaseTypeMySQL,
	}

	client, err = clients.InitClient(configuration)
	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to connect to the database. Reason: %s", err))
		return
	}

	if err = createTables(); err != nil {
		panic(err)
	}

	fmt.Println("Tables created")

	if err = triggerSelectQueries(); err != nil {
		panic(err)
	}

	fmt.Println("Select queries executed")
	if err = triggerInsert(); err != nil {
		panic(err)
	}

	fmt.Println("Insert queries executed")

	if err = triggerUpdate(); err != nil {
		panic(err)
	}

	fmt.Println("Update queries executed")

	if err = triggerTransactions(); err != nil {
		panic(err)
	}

	fmt.Println("Transaction queries executed")

	if err = triggerAlter(); err != nil {
		panic(err)
	}

	fmt.Println("Alter queries executed")

	if err = triggerDelete(); err != nil {
		panic(err)
	}

	fmt.Println("Delete queries executed")

	if err = triggerTableDrop(); err != nil {
		panic(err)
	}

	fmt.Println("Drop queries executed")

	if err = triggerTableRenaming(); err != nil {
		panic(err)
	}

	fmt.Println("Tables renamed and dropped queries executed")
	fmt.Println("All good")
}

func createTables() error {
	another = new(AnotherModel)
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
	q := new(clients.Query).Create(another).IfNotExists()
	_, err := client.Execute(q)
	if err != nil {
		return err
	}

	model = new(TestTableModel)
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
	q = new(clients.Query).Create(model).
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
	}).IfNotExists()
	_, err = client.Execute(q)
	return err
}

func triggerTableRenaming() error {
	//We rename the another table
	q := new(clients.Query).Rename(another.GetTableName(), "new_another_table")
	_, err := client.Execute(q)
	if err != nil {
		fmt.Println(err)
		return err
	}

	q = new(clients.Query).Select([]interface{}{}).From(&dto.BaseModel{
		TableName: "new_another_table",
	})
	_, err = client.Execute(q)
	if err != nil {
		fmt.Println(err)
		return err
	}

	//We drop the table
	another.SetTableName("new_another_table")
	q = new(clients.Query).Drop(another)
	_, err = client.Execute(q)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func triggerTableDrop() error {
	//We drop the table
	q := new(clients.Query).Drop(model)
	_, err := client.Execute(q)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func triggerDelete() error {
	//We delete item from the table
	q := new(clients.Query).Delete().
		From(model).
		Where(query.Where{
			First:    "id",
			Operator: "=",
			Second:   "1",
		})
	_, err := client.Execute(q)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func triggerAlter() error {
	//We alter table
	q := new(clients.Query).Alter(model).
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
	_, err := client.Execute(q)
	if err != nil {
		fmt.Println(err)
		return err
	}

	//We add new indexes
	q = new(clients.Query).Alter(model).
		AddIndex(dto.Index{
			Name:   "my_brand_new_index",
			Target: model.GetTableName(),
			Key:    "new_column",
			Unique: false,
		})
	_, err = client.Execute(q)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func triggerSelectQueries() error {
	//We select specific columns from the table
	var columns = []interface{}{"id", "another_id"}
	q := new(clients.Query).Select(columns).From(model)
	_, err := client.Execute(q)
	if err != nil {
		fmt.Println(err)
		return err
	}

	//We select specific columns from the table
	q = new(clients.Query).Select(columns).From(model.TableName)
	_, err = client.Execute(q)
	if err != nil {
		fmt.Println(err)
		return err
	}

	//We do select with join to other table
	q = new(clients.Query).
		Select([]interface{}{model.GetTableName() + ".id", "another_id"}).
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
	_, err = client.Execute(q)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func triggerInsert() error {
	//We insert new item into our table
	q := new(clients.Query).Insert(model)
	_, err := client.Execute(q)
	if err != nil {
		fmt.Println(err)
		return err
	}

	model.AddModelField(dto.ModelField{
		Name:          "test_field",
		Type:          dto.VarcharColumnType,
		Value:         "something2",
		Default:       "",
		Length:        255,
		IsNullable:    true,
		IsPrimaryKey:  false,
		AutoIncrement: false,
	})

	//We insert new item into our table
	q = new(clients.Query).Insert(model)
	_, err = client.Execute(q)
	if err != nil {
		fmt.Println(err)
		return err
	}

	//We select all model columns from our table
	q = new(clients.Query).Select(model.GetColumns()).From(model)
	_, err = client.Execute(q)
	if err != nil {
		fmt.Println(err)
		return err
	}

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
	_, err = client.Execute(q)
	if err != nil {
		fmt.Println(err)
		return err
	}

	//We do select with the more complex WHERE clause
	//The output will be:
	//SELECT id, id, id, another_id, test_field, test_field2 FROM test_table_name WHERE (id = 1 OR id = 2) OR id = 3
	q = new(clients.Query).Select(model.GetColumns()).
		From(model).
		Where(query.Where{
			First: query.Where{
				First:    "id",
				Operator: "=",
				Second:   "1",
			},
			Operator: "",
			Second: query.Where{
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
	_, err = client.Execute(q)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func triggerUpdate() error {
	//We update the table
	model.AddModelField(dto.ModelField{
		Name:  "test_field2",
		Value: "test test test",
	})
	q := new(clients.Query).Update(model).Where(query.Where{
		First:    model.GetPrimaryKey().Name,
		Operator: "=",
		Second: query.Bind{
			Field: model.GetPrimaryKey().Name,
			Value: model.GetPrimaryKey().Value,
		},
	})
	_, err := client.Execute(q)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func triggerTransactions() error {
	q := new(clients.Query).BeginTransaction()
	_, err := client.Execute(q)
	if err != nil {
		return err
	}

	//We update the table
	model.AddModelField(dto.ModelField{
		Name:  "test_field2",
		Value: "another test",
	})
	q = new(clients.Query).Update(model).Where(query.Where{
		First:    model.GetPrimaryKey().Name,
		Operator: "=",
		Second: query.Bind{
			Field: model.GetPrimaryKey().Name,
			Value: model.GetPrimaryKey().Value,
		},
	})
	_, err = client.Execute(q)
	if err != nil {
		return err
	}

	q = new(clients.Query).CommitTransaction()
	_, err = client.Execute(q)
	if err != nil {
		return err
	}

	//Now we need to trigger rollback scenario
	q = new(clients.Query).BeginTransaction()
	_, err = client.Execute(q)
	if err != nil {
		return err
	}

	//We update the table
	model.AddModelField(dto.ModelField{
		Name:  "test_field2",
		Value: "__SHOULD_NOT_BE_UPDATED__",
	})
	q = new(clients.Query).Update(model).Where(query.Where{
		First:    model.GetPrimaryKey().Name,
		Operator: "=",
		Second: query.Bind{
			Field: model.GetPrimaryKey().Name,
			Value: model.GetPrimaryKey().Value,
		},
	})
	_, err = client.Execute(q)
	if err != nil {
		return err
	}

	//Now we need to trigger rollback scenario
	q = new(clients.Query).RollbackTransaction()
	_, err = client.Execute(q)
	if err != nil {
		return err
	}

	//Now we need to trigger rollback scenario
	q = new(clients.Query).Select(model.GetColumns()).From(model).Where(query.Where{
		First:    "test_field2",
		Operator: "=",
		Second: query.Bind{
			Field: "test_field2",
			Value: "__SHOULD_NOT_BE_UPDATED__",
		},
	})
	res, err := client.Execute(q)
	if err != nil {
		return err
	}

	if len(res.Items()) > 0 {
		return errors.New("there should be no items with that value")
	}

	return nil
}
