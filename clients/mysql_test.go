package clients

import (
	"github.com/sharovik/orm/dto"
	"github.com/sharovik/orm/query"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	MySqlSelectCases = [...]expectation{
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
		{
			Expected: `SELECT * FROM test_table_name WHERE test_table_name2.relation_id = 2 OR col1 = "test"`,
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
					Type:     query.WhereOrType,
				})),
		},
		{
			Expected: `SELECT * FROM test_table_name WHERE test_table_name2.relation_id = 2 OR col1 = "test" NOT col2 = "test"`,
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
					Type:     query.WhereOrType,
				}).
				Where(query.Where{
					First:    "col2",
					Operator: "=",
					Second:   `"test"`,
					Type:     query.WhereNotType,
				})),
		},
		{
			Expected: `SELECT * FROM test_table_name WHERE (test_table_name2.relation_id = 2 OR col1 = "test") AND col2 = "test"`,
			Original: SQLiteClient{}.ToSql(new(Query).Select([]interface{}{}).
				From(&model).
				Where(query.Where{
					First: query.Where{
						First:    "test_table_name2.relation_id",
						Operator: "=",
						Second:   "2",
					},
					Operator: "",
					Second: query.Where{
						First:    "col1",
						Operator: "=",
						Second:   `"test"`,
						Type:     query.WhereOrType,
					},
					Type: query.WhereAndType,
				}).
				Where(query.Where{
					First:    "col2",
					Operator: "=",
					Second:   `"test"`,
					Type:     query.WhereAndType,
				})),
		},
		{
			Expected: `SELECT * FROM test_table_name WHERE ((test_table_name2.relation_id = 2 AND col1 = "test") OR col1 = "test") AND col2 = "test"`,
			Original: SQLiteClient{}.ToSql(new(Query).Select([]interface{}{}).
				From(&model).
				Where(query.Where{
					First: query.Where{
						First:    query.Where{
							First:    "test_table_name2.relation_id",
							Operator: "=",
							Second:   "2",
						},
						Operator: "",
						Second:   query.Where{
							First:    "col1",
							Operator: "=",
							Second:   `"test"`,
							Type:     query.WhereAndType,
						},
					},
					Operator: "",
					Second: query.Where{
						First:    "col1",
						Operator: "=",
						Second:   `"test"`,
						Type:     query.WhereOrType,
					},
					Type: query.WhereAndType,
				}).
				Where(query.Where{
					First:    "col2",
					Operator: "=",
					Second:   `"test"`,
					Type:     query.WhereAndType,
				})),
		},
	}
)

func TestMySQLClient_SelectToSql(t *testing.T) {
	for _, testCase := range MySqlSelectCases {
		assert.Equal(t, testCase.Expected, testCase.Original)
	}
}

func TestMySQLClient_InsertToSql(t *testing.T) {
	var (
		model     = initTestModel("test_table_name")
		testCases = [...]expectation{
			{
				Expected: "INSERT INTO test_table_name (relation_id, col1, col2, col3) VALUES (?, ?, ?, ?)",
				Original: MySQLClient{}.ToSql(new(Query).Insert(&model)),
			},
			{
				Expected: "INSERT INTO test_table_name (relation_id, col1, col2, col3) SELECT * FROM test_table_name1",
				Original: MySQLClient{}.ToSql(new(Query).Insert(&model).Values(new(Query).Select([]interface{}{}).From(&dto.BaseModel{
					TableName: "test_table_name1",
				}))),
			},
		}
	)

	for _, testCase := range testCases {
		assert.Equal(t, testCase.Expected, testCase.Original)
	}
}

func TestMySQLClient_DropToSql(t *testing.T) {
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

func TestMySQLClient_AlterToSql(t *testing.T) {
	var (
		model     = initTestModel("test_table_name")
		testCases = [...]expectation{
			{
				Expected: "ALTER TABLE test_table_name\nADD new_field integer(10) NULL DEFAULT 1",
				Original: MySQLClient{}.ToSql(new(Query).Alter(&model).AddColumn(dto.ModelField{
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
				Expected: "ALTER TABLE test_table_name\nADD INDEX my_brand_new_index (key_id)",
				Original: MySQLClient{}.ToSql(new(Query).Alter(&model).AddIndex(dto.Index{
					Name:   "my_brand_new_index",
					Target: "test_table_name",
					Key:    "key_id",
					Unique: false,
				})),
			},
			{
				Expected: "ALTER TABLE test_table_name\nADD UNIQUE INDEX my_brand_unique_new_index (key_id)",
				Original: MySQLClient{}.ToSql(new(Query).Alter(&model).AddIndex(dto.Index{
					Name:   "my_brand_unique_new_index",
					Target: "test_table_name",
					Key:    "key_id",
					Unique: true,
				})),
			},
			{
				Expected: "ALTER TABLE test_table_name\nADD UNIQUE INDEX my_brand_unique_new_index (key_id)",
				Original: MySQLClient{}.ToSql(new(Query).Alter(&model).AddIndex(dto.Index{
					Name:   "my_brand_unique_new_index",
					Target: "test_table_name",
					Key:    "key_id",
					Unique: true,
				})),
			},
		}
	)

	for _, testCase := range testCases {
		assert.Equal(t, testCase.Expected, testCase.Original)
	}
}

func TestMySQLClient_RenameToSql(t *testing.T) {
	var (
		model     = initTestModel("test_table_name")
		testCases = [...]expectation{
			{
				Expected: "ALTER TABLE `test_table_name` RENAME TO `new_test_table`",
				Original: MySQLClient{}.ToSql(new(Query).Rename(model.GetTableName(), "new_test_table")),
			},
			{
				Expected: "ALTER TABLE `test_table` RENAME TO `new_test_table`",
				Original: MySQLClient{}.ToSql(new(Query).Rename("test_table", "new_test_table")),
			},
		}
	)

	for _, testCase := range testCases {
		assert.Equal(t, testCase.Expected, testCase.Original)
	}
}

func TestMySQLClient_CreateToSql(t *testing.T) {
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
		otherModel = dto.BaseModel{
			TableName:  "some_other_table",
			PrimaryKey: dto.ModelField{},
			Fields:     nil,
		}

		otherModel2 = dto.BaseModel{
			TableName:  "some_other_table2",
			PrimaryKey: dto.ModelField{},
			Fields:     nil,
		}
		testCases = [...]expectation{
			{
				Expected: "CREATE TABLE test_table_name (id INTEGER CONSTRAINT test_table_name_pk primary key autoincrement, relation_id INTEGER NOT NULL, relation_id2 INTEGER NOT NULL, title VARCHAR DEFAULT \"test\" NOT NULL, description VARCHAR NULL,\nCONSTRAINT event_id\nFOREIGN KEY (event_id)\n REFERENCES some_other_table (id)\nON DELETE CASCADE\nON UPDATE NO ACTION,\nCONSTRAINT scenario_id\nFOREIGN KEY (scenario_id)\n REFERENCES some_other_table2 (id)\nON DELETE CASCADE\nON UPDATE NO ACTION); CREATE INDEX user_id_index \nON test_table_name (user);\nCREATE INDEX channel_index \nON test_table_name (channel);\nCREATE INDEX created_index \nON test_table_name (created);",
				Original: SQLiteClient{}.ToSql(new(Query).
					Create(&model).
					AddIndex(dto.Index{
						Name:   "user_id_index",
						Target: model.GetTableName(),
						Key:    "user",
						Unique: false,
					}).
					AddIndex(dto.Index{
						Name:   "channel_index",
						Target: model.GetTableName(),
						Key:    "channel",
						Unique: false,
					}).
					AddIndex(dto.Index{
						Name:   "created_index",
						Target: model.GetTableName(),
						Key:    "created",
						Unique: false,
					}).
					AddForeignKey(dto.ForeignKey{
						Name: "event_id",
						Target: query.Reference{
							Table: otherModel.GetTableName(),
							Key:   "id",
						},
						With: query.Reference{
							Table: model.GetTableName(),
							Key:   "event_id",
						},
						OnDelete: dto.CascadeAction,
						OnUpdate: dto.NoActionAction,
					}).
					AddForeignKey(dto.ForeignKey{
						Name: "scenario_id",
						Target: query.Reference{
							Table: otherModel2.GetTableName(),
							Key:   "id",
						},
						With: query.Reference{
							Table: model.GetTableName(),
							Key:   "scenario_id",
						},
						OnDelete: dto.CascadeAction,
						OnUpdate: dto.NoActionAction,
					})),
			},
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
			{
				Expected: "CREATE TABLE IF NOT EXISTS test_table_name (id INTEGER CONSTRAINT test_table_name_pk primary key autoincrement, relation_id INTEGER NOT NULL, relation_id2 INTEGER NOT NULL, title VARCHAR DEFAULT \"test\" NOT NULL, description VARCHAR NULL); CREATE UNIQUE INDEX the_index_name \nON test_table_name (relation_id);",
				Original: SQLiteClient{}.ToSql(new(Query).Create(&model).
					IfNotExists().
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

func TestMySQLClient_UpdateToSql(t *testing.T) {
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

func TestMySQLClient_DeleteToSql(t *testing.T) {
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
