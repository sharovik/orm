package clients

import (
	"fmt"
	"github.com/sharovik/orm/dto"
	"github.com/sharovik/orm/query"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

const testSQLiteDatabasePath = "testing.sqlite"

type expectation struct {
	Expected interface{}
	Original interface{}
}

var (
	columns           = []interface{}{"col1", "col2"}
	model             = initTestModel("test_table_name")
	model2            = initTestModel("test_table_name2")
	SqliteSelectCases = [...]expectation{
		{
			Expected: "SELECT col1, col2 FROM test_table_name",
			Original: SQLiteClient{}.ToSql(new(Query).Select(columns).From(&model)),
		},
		{
			Expected: "SELECT * FROM test_table_name",
			Original: SQLiteClient{}.ToSql(new(Query).Select([]interface{}{}).From(&model)),
		},
		{
			Expected: "SELECT * FROM test_table_name LEFT JOIN test_table_name2 ON (test_table_name2.id = test_table_name.relation_id)",
			Original: SQLiteClient{}.ToSql(new(Query).Select(nil).
				From(&model).
				Join(query.Join{
					Target: query.Reference{
						Table: model2.GetTableName(),
						Key:   model2.GetField("id").Name,
					},
					With: query.Reference{
						Table: model.GetTableName(),
						Key:   model.GetField("relation_id").Name,
					},
					Condition: "=",
					Type:      query.LeftJoinType,
				})),
		},
		{
			Expected: "SELECT * FROM test_table_name LEFT JOIN test_table_name2 ON (test_table_name2.id = test_table_name.relation_id) ORDER BY id DESC",
			Original: SQLiteClient{}.ToSql(new(Query).Select([]interface{}{}).
				From(&model).
				Join(query.Join{
					Target: query.Reference{
						Table: model2.GetTableName(),
						Key:   model2.GetField("id").Name,
					},
					With: query.Reference{
						Table: model.GetTableName(),
						Key:   model.GetField("relation_id").Name,
					},
					Condition: "=",
					Type:      query.LeftJoinType,
				}).OrderBy(model.GetPrimaryKey().Name, query.OrderDirectionDesc)),
		},
		{
			Expected: "SELECT * FROM test_table_name LEFT JOIN test_table_name2 ON (test_table_name2.id = test_table_name.relation_id) ORDER BY id DESC",
			Original: SQLiteClient{}.ToSql(new(Query).Select([]interface{}{}).
				From(&model).
				Join(query.Join{
					Target: query.Reference{
						Table: model2.GetTableName(),
						Key:   model2.GetField("id").Name,
					},
					With: query.Reference{
						Table: model.GetTableName(),
						Key:   model.GetField("relation_id").Name,
					},
					Condition: "=",
					Type:      query.LeftJoinType,
				}).OrderBy(model.GetPrimaryKey().Name, query.OrderDirectionDesc)),
		},
		{
			Expected: "SELECT * FROM test_table_name LEFT JOIN test_table_name2 ON (test_table_name2.id = test_table_name.relation_id) GROUP BY test_table_name.id ORDER BY id DESC",
			Original: SQLiteClient{}.ToSql(new(Query).Select([]interface{}{}).
				From(&model).
				Join(query.Join{
					Target: query.Reference{
						Table: model2.GetTableName(),
						Key:   model2.GetField("id").Name,
					},
					With: query.Reference{
						Table: model.GetTableName(),
						Key:   model.GetField("relation_id").Name,
					},
					Condition: "=",
					Type:      query.LeftJoinType,
				}).
				OrderBy(model.GetPrimaryKey().Name, query.OrderDirectionDesc).
				GroupBy("test_table_name.id")),
		},
		{
			Expected: "SELECT * FROM test_table_name LEFT JOIN test_table_name2 ON (test_table_name2.id = test_table_name.relation_id) WHERE test_table_name2.relation_id = 2 GROUP BY test_table_name.id ORDER BY id DESC",
			Original: SQLiteClient{}.ToSql(new(Query).Select([]interface{}{}).
				From(&model).
				Join(query.Join{
					Target: query.Reference{
						Table: model2.GetTableName(),
						Key:   model2.GetField("id").Name,
					},
					With: query.Reference{
						Table: model.GetTableName(),
						Key:   model.GetField("relation_id").Name,
					},
					Condition: "=",
					Type:      query.LeftJoinType,
				}).
				Where(query.Where{
					First:    "test_table_name2.relation_id",
					Operator: "=",
					Second:   "2",
				}).
				OrderBy(model.GetPrimaryKey().Name, query.OrderDirectionDesc).
				GroupBy("test_table_name.id")),
		},
		{
			Expected: `SELECT * FROM test_table_name WHERE test_table_name2.relation_id = 2 AND col1 = "test" LIMIT 11`,
			Original: SQLiteClient{}.ToSql(new(Query).Select([]interface{}{}).
				From(&model).
				Where(query.Where{
					First:    "test_table_name2.relation_id",
					Operator: "=",
					Second:   "2",
				}).
				Where(query.Where{
					First:    "col1",
					Operator: "=",
					Second:   `"test"`,
				}).
				Limit(query.Limit{
					From: 0,
					To:   11,
				})),
		},
		{
			Expected: `SELECT * FROM test_table_name WHERE test_table_name2.relation_id = 2 AND col1 = "test" AND ? = ? LIMIT 11`,
			Original: SQLiteClient{}.ToSql(new(Query).Select([]interface{}{}).
				From(&model).
				Where(query.Where{
					First:    "test_table_name2.relation_id",
					Operator: "=",
					Second:   "2",
				}).
				Where(query.Where{
					First:    "col1",
					Operator: "=",
					Second:   `"test"`,
				}).
				Where(query.Where{
					First: query.Bind{
						Field: "",
						Value: 1,
					},
					Operator: "=",
					Second: query.Bind{
						Field: "",
						Value: 1,
					},
				}).
				Limit(query.Limit{
					From: 0,
					To:   11,
				})),
		},
	}
)

func TestSQLiteClient_SelectToSql(t *testing.T) {
	for _, testCase := range SqliteSelectCases {
		assert.Equal(t, testCase.Expected, testCase.Original)
	}
}

func TestSQLiteClient_InsertToSql(t *testing.T) {
	var (
		model     = initTestModel("test_table_name")
		testCases = [...]expectation{
			{
				Expected: "INSERT INTO test_table_name (relation_id, col1, col2, col3) VALUES (?, ?, ?, ?)",
				Original: SQLiteClient{}.ToSql(new(Query).Insert(&model)),
			},
		}
	)

	for _, testCase := range testCases {
		assert.Equal(t, testCase.Expected, testCase.Original)
	}
}

func TestSQLiteClient_DropToSql(t *testing.T) {
	var (
		model     = initTestModel("test_table_name")
		testCases = [...]expectation{
			{
				Expected: "DROP TABLE test_table_name",
				Original: SQLiteClient{}.ToSql(new(Query).Drop(&model)),
			},
		}
	)

	for _, testCase := range testCases {
		assert.Equal(t, testCase.Expected, testCase.Original)
	}
}

func TestSQLiteClient_AlterToSql(t *testing.T) {
	var (
		model     = initTestModel("test_table_name")
		testCases = [...]expectation{
			{
				Expected: "ALTER TABLE test_table_name\nADD COLUMN new_field integer DEFAULT 1 NOT NULL",
				Original: SQLiteClient{}.ToSql(new(Query).Alter(&model).AddColumn(dto.ModelField{
					Name:          "new_field",
					Type:          "integer",
					Value:         nil,
					Default:       1,
					Length:        10,
					IsNullable:    false,
					AutoIncrement: false,
				})),
			},
			{
				Expected: "ALTER TABLE test_table_name\nADD COLUMN new_field integer DEFAULT 1 NOT NULL",
				Original: SQLiteClient{}.ToSql(new(Query).Alter(&model).AddColumn(dto.ModelField{
					Name:          "new_field",
					Type:          "integer",
					Value:         nil,
					Default:       1,
					Length:        10,
					IsNullable:    false,
					AutoIncrement: false,
				})),
			},
		}
	)

	for _, testCase := range testCases {
		assert.Equal(t, testCase.Expected, testCase.Original)
	}
}

func TestSQLiteClient_CreateToSql(t *testing.T) {
	var (
		model = dto.BaseModel{
			TableName: "test_table_name",
			Fields: []interface{}{
				dto.ModelField{
					Name: "relation_id",
					Type: dto.IntegerColumnType,
				},
				dto.ModelField{
					Name: "relation_id2",
					Type: dto.IntegerColumnType,
				},
				dto.ModelField{
					Name:    "title",
					Type:    dto.VarcharColumnType,
					Default: "test",
				},
				dto.ModelField{
					Name:       "description",
					Type:       dto.VarcharColumnType,
					IsNullable: true,
				},
			},
			PrimaryKey: dto.ModelField{
				Name:          "id",
				Type:          dto.IntegerColumnType,
				AutoIncrement: true,
			},
		}
		testCases = [...]expectation{
			{
				Expected: "CREATE TABLE test_table_name (id INTEGER CONSTRAINT test_table_name_pk primary key autoincrement, relation_id INTEGER NOT NULL, relation_id2 INTEGER NOT NULL, title VARCHAR DEFAULT \"test\" NOT NULL, description VARCHAR NULL);",
				Original: SQLiteClient{}.ToSql(new(Query).Create(&model)),
			},
			{
				Expected: "CREATE TABLE test_table_name (id INTEGER CONSTRAINT test_table_name_pk primary key autoincrement, relation_id INTEGER NOT NULL, relation_id2 INTEGER NOT NULL, title VARCHAR DEFAULT \"test\" NOT NULL, description VARCHAR NULL,\nCONSTRAINT fk_test\nFOREIGN KEY (relation_id)\n REFERENCES test_table_name2 (id)\nON DELETE NO ACTION\nON UPDATE NO ACTION);",
				Original: SQLiteClient{}.ToSql(new(Query).Create(&model).
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
					})),
			},
			{
				Expected: "CREATE TABLE test_table_name (id INTEGER CONSTRAINT test_table_name_pk primary key autoincrement, relation_id INTEGER NOT NULL, relation_id2 INTEGER NOT NULL, title VARCHAR DEFAULT \"test\" NOT NULL, description VARCHAR NULL,\nCONSTRAINT fk_test\nFOREIGN KEY (relation_id)\n REFERENCES test_table_name2 (id)\nON DELETE NO ACTION\nON UPDATE NO ACTION,\nCONSTRAINT fk_test2\nFOREIGN KEY (relation_id2)\n REFERENCES test_table_name3 (id)\nON DELETE CASCADE\nON UPDATE NO ACTION);",
				Original: SQLiteClient{}.ToSql(new(Query).Create(&model).
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
					}).AddForeignKey(dto.ForeignKey{
					Name: "fk_test2",
					Target: query.Reference{
						Table: "test_table_name3",
						Key:   "id",
					},
					With: query.Reference{
						Table: "test_table_name",
						Key:   "relation_id2",
					},
					OnDelete: dto.CascadeAction,
					OnUpdate: "",
				})),
			},
			{
				Expected: "CREATE TABLE test_table_name (id INTEGER CONSTRAINT test_table_name_pk primary key autoincrement, relation_id INTEGER NOT NULL, relation_id2 INTEGER NOT NULL, title VARCHAR DEFAULT \"test\" NOT NULL, description VARCHAR NULL); CREATE INDEX the_index_name \nON test_table_name (relation_id);",
				Original: SQLiteClient{}.ToSql(new(Query).Create(&model).
					AddIndex(dto.Index{
						Name:   "the_index_name",
						Target: model.GetTableName(),
						Key:    "relation_id",
						Unique: false,
					})),
			},
			{
				Expected: "CREATE TABLE test_table_name (id INTEGER CONSTRAINT test_table_name_pk primary key autoincrement, relation_id INTEGER NOT NULL, relation_id2 INTEGER NOT NULL, title VARCHAR DEFAULT \"test\" NOT NULL, description VARCHAR NULL); CREATE UNIQUE INDEX the_index_name \nON test_table_name (relation_id);",
				Original: SQLiteClient{}.ToSql(new(Query).Create(&model).
					AddIndex(dto.Index{
						Name:   "the_index_name",
						Target: model.GetTableName(),
						Key:    "relation_id",
						Unique: true,
					})),
			},
		}
	)

	for _, testCase := range testCases {
		assert.Equal(t, testCase.Expected, testCase.Original)
	}
}

func TestSQLiteClient_UpdateToSql(t *testing.T) {
	var (
		model     = initTestModel("test_table_name")
		testCases = [...]expectation{
			{
				Expected: "UPDATE test_table_name SET relation_id = ?, col1 = ?, col2 = ?, col3 = ?",
				Original: SQLiteClient{}.ToSql(new(Query).Update(&model)),
			},
			{
				Expected: "UPDATE test_table_name SET relation_id = ?, col1 = ?, col2 = ?, col3 = ? LEFT JOIN test ON (test.ref_id = test_table_name.id)",
				Original: SQLiteClient{}.ToSql(new(Query).Update(&model).Join(query.Join{
					Target: query.Reference{
						Table: "test",
						Key:   "ref_id",
					},
					With: query.Reference{
						Table: model.GetTableName(),
						Key:   model.GetPrimaryKey().Name,
					},
					Condition: "=",
					Type:      query.LeftJoinType,
				})),
			},
			{
				Expected: "UPDATE test_table_name SET relation_id = ?, col1 = ?, col2 = ?, col3 = ? WHERE relation_id = test",
				Original: SQLiteClient{}.ToSql(new(Query).Update(&model).Where(query.Where{
					First:    "relation_id",
					Operator: "=",
					Second:   "test",
				})),
			},
		}
	)

	for _, testCase := range testCases {
		assert.Equal(t, testCase.Expected, testCase.Original)
	}
}

func TestSQLiteClient_DeleteToSql(t *testing.T) {
	var (
		model     = initTestModel("test_table_name")
		model2    = initTestModel("test_table_name2")
		testCases = [...]expectation{
			{
				Expected: "DELETE FROM test_table_name",
				Original: SQLiteClient{}.ToSql(new(Query).Delete().From(&model)),
			},
			{
				Expected: "DELETE FROM test_table_name LEFT JOIN test_table_name2 ON (test_table_name2.id = test_table_name.relation_id)",
				Original: SQLiteClient{}.ToSql(new(Query).Delete().
					From(&model).
					Join(query.Join{
						Target: query.Reference{
							Table: model2.GetTableName(),
							Key:   model2.GetField("id").Name,
						},
						With: query.Reference{
							Table: model.GetTableName(),
							Key:   model.GetField("relation_id").Name,
						},
						Condition: "=",
						Type:      query.LeftJoinType,
					})),
			},
			{
				Expected: "DELETE FROM test_table_name ORDER BY id DESC",
				Original: SQLiteClient{}.ToSql(new(Query).
					Delete().
					From(&model).
					OrderBy(model.GetPrimaryKey().Name, query.OrderDirectionDesc)),
			},
			{
				Expected: "DELETE FROM test_table_name GROUP BY test_table_name.id",
				Original: SQLiteClient{}.ToSql(new(Query).
					Delete().
					From(&model).
					GroupBy("test_table_name.id")),
			},
			{
				Expected: "DELETE FROM test_table_name WHERE test_table_name.relation_id = 2",
				Original: SQLiteClient{}.ToSql(new(Query).
					Delete().
					From(&model).
					Where(query.Where{
						First:    "test_table_name.relation_id",
						Operator: "=",
						Second:   "2",
					})),
			},
			{
				Expected: `DELETE FROM test_table_name LIMIT 11`,
				Original: SQLiteClient{}.ToSql(new(Query).
					Delete().
					From(&model).
					Limit(query.Limit{
						From: 0,
						To:   11,
					})),
			},
		}
	)

	for _, testCase := range testCases {
		assert.Equal(t, testCase.Expected, testCase.Original)
	}
}

func TestSQLiteClient_Execute(t *testing.T) {
	removeDatabase()
	initDatabase()

	sqliteClient, err := SQLiteClient{}.Connect(DatabaseConfig{
		Host:     testSQLiteDatabasePath,
		Username: "",
		Password: "",
		Port:     0,
	})
	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to connect to the database. Reason: %s", err))
		return
	}

	model := initTestModel("testing")

	//First let's create a test table
	res, err := sqliteClient.Execute(new(Query).Create(&model))
	assert.NoError(t, err)
	assert.NoError(t, res.Err)
	assert.Len(t, res.Items(), 0)

	//Secondary let's insert something in our test table
	res, err = sqliteClient.Execute(new(Query).Insert(&model))
	assert.NoError(t, err)
	assert.NoError(t, res.Error())
	assert.Equal(t, int64(1), res.LastInsertID())
	assert.Len(t, res.Items(), 0)

	//Now let's select the data from our table and check if it is correct
	q := new(Query).Select(model.GetColumns()).From(&model)
	res, err = sqliteClient.Execute(q)
	assert.NoError(t, err)
	assert.NoError(t, res.Error())
	assert.Len(t, res.Items(), 1)

	var expected []interface{}
	for _, field := range model.GetColumns() {
		switch v := field.(type) {
		case dto.ModelField:
			if v.Name == "id" {
				v.Value = int64(1)
			}

			expected = append(expected, dto.ModelField{
				Name:          v.Name,
				Type:          v.Type,
				Value:         v.Value,
				Default:       v.Default,
				Length:        v.Length,
				IsNullable:    v.IsNullable,
				IsPrimaryKey:  false,
				AutoIncrement: false,
			})
		}
	}
	assert.Equal(t, expected, res.Items()[0].GetColumns())

	model.SetField("relation_id", 2)

	//And it's time for updates
	res, err = sqliteClient.Execute(new(Query).Update(&model))
	assert.NoError(t, err)
	assert.NoError(t, res.Error())
	assert.Equal(t, int64(1), res.LastInsertID())
	assert.Len(t, res.Items(), 0)

	//Now let's select the data from our table and check if it is correct
	var columns = []interface{}{"id", "relation_id"}
	q = new(Query).Select(columns).From(&model)
	res, err = sqliteClient.Execute(q)
	assert.NoError(t, err)
	assert.NoError(t, res.Err)
	assert.Len(t, res.Items(), 1)

	expected = []interface{}{
		dto.ModelField{
			Name:  "id",
			Type:  "INTEGER",
			Value: int64(1),
		},
		dto.ModelField{
			Name:  "relation_id",
			Type:  "INTEGER",
			Value: int64(2),
		},
	}
	assert.Equal(t, expected, res.Items()[0].GetColumns())

	//Now we delete that row
	res, err = sqliteClient.Execute(new(Query).Delete().From(&model))
	assert.NoError(t, err)
	assert.NoError(t, res.Error())
	assert.Equal(t, int64(1), res.LastInsertID())
	assert.Len(t, res.Items(), 0)

	//And we make sure we really delete the row
	res, err = sqliteClient.Execute(new(Query).Select(model.GetColumns()).From(&model).Where(query.Where{
		First:    "id",
		Operator: "=",
		Second:   query.Bind{
			Field: "id",
			Value: 1,
		},
	}))
	assert.NoError(t, err)
	assert.NoError(t, res.Error())
	assert.Len(t, res.Items(), 0)

	removeDatabase()
}

func initTestModel(table string) dto.BaseModel {
	return dto.BaseModel{
		TableName: table,
		Fields: []interface{}{
			dto.ModelField{
				Name:  "relation_id",
				Type:  "INTEGER",
				Value: int64(1),
			},
			dto.ModelField{
				Name:  "col1",
				Type:  "INTEGER",
				Value: int64(1),
			},
			dto.ModelField{
				Name:  "col2",
				Type:  "INTEGER",
				Value: int64(2),
			},
			dto.ModelField{
				Name:  "col3",
				Type:  "VARCHAR",
				Value: "Test",
			},
		},
		PrimaryKey: dto.ModelField{
			Name:          "id",
			Type:          dto.IntegerColumnType,
			AutoIncrement: true,
		},
	}
}

func initDatabase() {
	_, err := os.Create(testSQLiteDatabasePath)
	if err != nil {
		fmt.Println("Failed to create database file: " + err.Error())
	}
}

func removeDatabase() {
	err := os.Remove(testSQLiteDatabasePath)
	if err != nil {
		fmt.Println("Failed to remove database file: " + err.Error())
	}
}
