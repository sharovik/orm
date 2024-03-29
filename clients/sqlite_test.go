package clients

import (
	"fmt"
	"os"
	"testing"

	"github.com/sharovik/orm/dto"
	"github.com/sharovik/orm/query"
	"github.com/stretchr/testify/assert"
)

const testSQLiteDatabasePath = "testing.sqlite"

type expectation struct {
	Expected interface{}
	Original interface{}
}

var (
	columns           = []interface{}{"col1", "col2"}
	m                 = initTestModel("test_table_name")
	model2            = initTestModel("test_table_name2")
	SqliteSelectCases = [...]expectation{
		{
			Expected: "SELECT col1, col2 FROM test_table_name",
			Original: SQLiteClient{}.ToSql(new(Query).Select(columns).From(&m)),
		},
		{
			Expected: "SELECT * FROM test_table_name",
			Original: SQLiteClient{}.ToSql(new(Query).Select([]interface{}{}).From(&m)),
		},
		{
			Expected: "SELECT * FROM test_table_name LEFT JOIN test_table_name2 ON (test_table_name2.id = test_table_name.relation_id)",
			Original: SQLiteClient{}.ToSql(new(Query).Select(nil).
				From(&m).
				Join(query.Join{
					Target: query.Reference{
						Table: model2.GetTableName(),
						Key:   model2.GetField("id").Name,
					},
					With: query.Reference{
						Table: m.GetTableName(),
						Key:   m.GetField("relation_id").Name,
					},
					Condition: "=",
					Type:      query.LeftJoinType,
				})),
		},
		{
			Expected: "SELECT * FROM test_table_name LEFT JOIN test_table_name2 ON (test_table_name2.id = test_table_name.relation_id) ORDER BY id DESC",
			Original: SQLiteClient{}.ToSql(new(Query).Select([]interface{}{}).
				From(&m).
				Join(query.Join{
					Target: query.Reference{
						Table: model2.GetTableName(),
						Key:   model2.GetField("id").Name,
					},
					With: query.Reference{
						Table: m.GetTableName(),
						Key:   m.GetField("relation_id").Name,
					},
					Condition: "=",
					Type:      query.LeftJoinType,
				}).OrderBy(m.GetPrimaryKey().Name, query.OrderDirectionDesc)),
		},
		{
			Expected: "SELECT * FROM test_table_name LEFT JOIN test_table_name2 ON (test_table_name2.id = test_table_name.relation_id) ORDER BY id DESC",
			Original: SQLiteClient{}.ToSql(new(Query).Select([]interface{}{}).
				From(&m).
				Join(query.Join{
					Target: query.Reference{
						Table: model2.GetTableName(),
						Key:   model2.GetField("id").Name,
					},
					With: query.Reference{
						Table: m.GetTableName(),
						Key:   m.GetField("relation_id").Name,
					},
					Condition: "=",
					Type:      query.LeftJoinType,
				}).OrderBy(m.GetPrimaryKey().Name, query.OrderDirectionDesc)),
		},
		{
			Expected: "SELECT * FROM test_table_name LEFT JOIN test_table_name2 ON (test_table_name2.id = test_table_name.relation_id) GROUP BY test_table_name.id ORDER BY id DESC",
			Original: SQLiteClient{}.ToSql(new(Query).Select([]interface{}{}).
				From(&m).
				Join(query.Join{
					Target: query.Reference{
						Table: model2.GetTableName(),
						Key:   model2.GetField("id").Name,
					},
					With: query.Reference{
						Table: m.GetTableName(),
						Key:   m.GetField("relation_id").Name,
					},
					Condition: "=",
					Type:      query.LeftJoinType,
				}).
				OrderBy(m.GetPrimaryKey().Name, query.OrderDirectionDesc).
				GroupBy("test_table_name.id")),
		},
		{
			Expected: "SELECT * FROM test_table_name LEFT JOIN test_table_name2 ON (test_table_name2.id = test_table_name.relation_id) WHERE test_table_name2.relation_id = 2 GROUP BY test_table_name.id ORDER BY id DESC",
			Original: SQLiteClient{}.ToSql(new(Query).Select([]interface{}{}).
				From(&m).
				Join(query.Join{
					Target: query.Reference{
						Table: model2.GetTableName(),
						Key:   model2.GetField("id").Name,
					},
					With: query.Reference{
						Table: m.GetTableName(),
						Key:   m.GetField("relation_id").Name,
					},
					Condition: "=",
					Type:      query.LeftJoinType,
				}).
				Where(query.Where{
					First:    "test_table_name2.relation_id",
					Operator: "=",
					Second:   "2",
				}).
				OrderBy(m.GetPrimaryKey().Name, query.OrderDirectionDesc).
				GroupBy("test_table_name.id")),
		},
		{
			Expected: `SELECT * FROM test_table_name WHERE test_table_name2.relation_id = 2 AND col1 = "test" LIMIT 11`,
			Original: SQLiteClient{}.ToSql(new(Query).Select([]interface{}{}).
				From(&m).
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
				From(&m).
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
				From(&m).
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
				From(&m).
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
				From(&m).
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
				From(&m).
				Where(query.Where{
					First: query.Where{
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
			{
				Expected: "INSERT INTO test_table_name (relation_id, col1, col2, col3) SELECT * FROM test_table_name1",
				Original: SQLiteClient{}.ToSql(new(Query).Insert(&model).Values(new(Query).Select([]interface{}{}).From(&dto.BaseModel{
					TableName: "test_table_name1",
				}))),
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

func TestSQLiteClient_RenameToSql(t *testing.T) {
	var (
		model     = initTestModel("test_table_name")
		testCases = [...]expectation{
			{
				Expected: "ALTER TABLE `test_table_name` RENAME TO `new_test_table`",
				Original: SQLiteClient{}.ToSql(new(Query).Rename(model.GetTableName(), "new_test_table")),
			},
			{
				Expected: "ALTER TABLE `test_table` RENAME TO `new_test_table`",
				Original: SQLiteClient{}.ToSql(new(Query).Rename("test_table", "new_test_table")),
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
				Expected: "ALTER TABLE test_table_name ADD COLUMN new_field integer DEFAULT 1 NOT NULL",
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
				Expected: "ALTER TABLE test_table_name ADD COLUMN new_field integer DEFAULT 1 NOT NULL",
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
				Expected: "CREATE INDEX my_brand_new_index on test_table_name (request_id)",
				Original: SQLiteClient{}.ToSql(new(Query).Alter(&model).AddIndex(dto.Index{
					Name:   "my_brand_new_index",
					Target: "test_table_name",
					Key:    "request_id",
					Unique: false,
				})),
			},
			{
				Expected: "CREATE UNIQUE INDEX my_brand_unique_new_index on test_table_name (request_id);\nCREATE INDEX my_brand_non_unique_new_index on test_table_name (name)",
				Original: SQLiteClient{}.ToSql(new(Query).Alter(&model).AddIndex(dto.Index{
					Name:   "my_brand_unique_new_index",
					Target: "test_table_name",
					Key:    "request_id",
					Unique: true,
				}).AddIndex(dto.Index{
					Name:   "my_brand_non_unique_new_index",
					Target: "test_table_name",
					Key:    "name",
				})),
			},
			{
				Expected: "CREATE INDEX my_brand_non_unique_new_index on test_table_name (name);\nDROP INDEX my_brand_unique_new_index",
				Original: SQLiteClient{}.ToSql(new(Query).Alter(&model).DropIndex(dto.Index{
					Name: "my_brand_unique_new_index",
				}).AddIndex(dto.Index{
					Name:   "my_brand_non_unique_new_index",
					Target: "test_table_name",
					Key:    "name",
				})),
			},
			{
				Expected: "CREATE TABLE temp_test_table_name (id INTEGER CONSTRAINT temp_test_table_name_pk primary key autoincrement, relation_id INTEGER NOT NULL, col1 INTEGER NOT NULL, col2 INTEGER NOT NULL);\nINSERT INTO temp_test_table_name (relation_id, col1, col2) SELECT relation_id, col1, col2 FROM test_table_name;\nALTER TABLE `test_table_name` RENAME TO `old_test_table_name`;\nALTER TABLE `temp_test_table_name` RENAME TO `test_table_name`;\nDROP TABLE old_test_table_name;",
				Original: SQLiteClient{}.ToSql(new(Query).Alter(&model).
					DropColumn(dto.ModelField{
						Name: "col3",
					})),
			},
			{
				Expected: "CREATE TABLE temp_test_table_name (id INTEGER CONSTRAINT temp_test_table_name_pk primary key autoincrement, relation_id INTEGER NOT NULL, col1 INTEGER NOT NULL, col2 INTEGER NOT NULL, col3 VARCHAR NOT NULL);\nINSERT INTO temp_test_table_name (relation_id, col1, col2, col3) SELECT relation_id, col1, col2, col3 FROM test_table_name;\nALTER TABLE `test_table_name` RENAME TO `old_test_table_name`;\nALTER TABLE `temp_test_table_name` RENAME TO `test_table_name`;\nDROP TABLE old_test_table_name;",
				Original: SQLiteClient{}.ToSql(new(Query).Alter(&model).
					DropForeignKey(dto.ForeignKey{
						Name: "test_foreign_key",
					})),
			},
			{
				Expected: "CREATE TABLE temp_test_table_name (id INTEGER CONSTRAINT temp_test_table_name_pk primary key autoincrement, relation_id INTEGER NOT NULL, col1 INTEGER NOT NULL, col2 INTEGER NOT NULL, col3 VARCHAR NOT NULL,\nCONSTRAINT fk_test\nFOREIGN KEY (relation_id)\n REFERENCES test_table_name2 (id)\nON DELETE NO ACTION\nON UPDATE NO ACTION);\nINSERT INTO temp_test_table_name (relation_id, col1, col2, col3) SELECT relation_id, col1, col2, col3 FROM test_table_name;\nALTER TABLE `test_table_name` RENAME TO `old_test_table_name`;\nALTER TABLE `temp_test_table_name` RENAME TO `test_table_name`;\nDROP TABLE old_test_table_name;",
				Original: SQLiteClient{}.ToSql(new(Query).Alter(&model).
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
		}
	)

	for _, testCase := range testCases {
		assert.Equal(t, testCase.Expected, testCase.Original)
	}
}

func TestSQLiteClient_CreateToSql(t *testing.T) {
	var model = dto.BaseModel{
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
	}
	model.SetPrimaryKey(dto.ModelField{
		Name:          "id",
		Type:          dto.IntegerColumnType,
		AutoIncrement: true,
	})
	var (
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
		fmt.Printf("Failed to connect to the database. Reason: %s\n", err)
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
				v.Value = 1
				v.IsPrimaryKey = false
				v.AutoIncrement = false
			}

			expected = append(expected, dto.ModelField{
				Name:          v.Name,
				Type:          v.Type,
				Value:         v.Value,
				Default:       v.Default,
				Length:        v.Length,
				IsNullable:    v.IsNullable,
				IsPrimaryKey:  v.IsPrimaryKey,
				AutoIncrement: v.AutoIncrement,
			})
		}
	}
	assert.Equal(t, expected, res.Items()[0].GetColumns())

	model.AddModelField(dto.ModelField{
		Name:  "relation_id",
		Value: 2,
	})

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

	exp := []dto.ModelField{
		{
			Name:          "id",
			Type:          "INTEGER",
			Value:         1,
			IsPrimaryKey:  false,
			AutoIncrement: false,
		},
		{
			Name:  "relation_id",
			Type:  "INTEGER",
			Value: 2,
		},
	}

	for _, field := range exp {
		assert.Equal(t, field, res.Items()[0].GetField(field.Name))
	}

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
		Second: query.Bind{
			Field: "id",
			Value: 1,
		},
	}))
	assert.NoError(t, err)
	assert.NoError(t, res.Error())
	assert.Len(t, res.Items(), 0)

	//And we make sure we really delete the row
	res, err = sqliteClient.Execute(new(Query).
		Alter(&model).
		AddColumn(dto.ModelField{
			Name:       "new_column",
			Type:       dto.IntegerColumnType,
			Default:    1,
			Length:     10,
			IsNullable: true,
		}).
		AddIndex(dto.Index{
			Name:   "the_index_name",
			Target: model.GetTableName(),
			Key:    "new_column",
			Unique: false,
		}),
	)
	assert.NoError(t, err)
	assert.NoError(t, res.Error())

	res, err = sqliteClient.Execute(new(Query).Select([]interface{}{"new_column"}).From(&model))
	assert.NoError(t, err)
	assert.NoError(t, res.Error())

	//We test drop foreign keys
	testNewTable := initTestModel("test_table")
	res, err = sqliteClient.Execute(new(Query).
		Create(&testNewTable).
		AddForeignKey(dto.ForeignKey{
			Name: "fk_test",
			Target: query.Reference{
				Table: model.GetTableName(),
				Key:   "id",
			},
			With: query.Reference{
				Table: "test_table",
				Key:   "col2",
			},
			OnDelete: "",
			OnUpdate: "",
		}).
		AddForeignKey(dto.ForeignKey{
			Name: "fk_test2",
			Target: query.Reference{
				Table: model.GetTableName(),
				Key:   "id",
			},
			With: query.Reference{
				Table: "test_table",
				Key:   "col2",
			},
			OnDelete: "",
			OnUpdate: "",
		}))
	assert.NoError(t, err)

	res, err = sqliteClient.Execute(new(Query).Insert(&testNewTable))
	assert.NoError(t, err)
	assert.NoError(t, res.Error())
	assert.Equal(t, int64(1), res.LastInsertID())
	assert.Len(t, res.Items(), 0)

	res, err = sqliteClient.Execute(new(Query).Insert(&testNewTable))
	assert.NoError(t, err)
	assert.NoError(t, res.Error())
	assert.Equal(t, int64(2), res.LastInsertID())
	assert.Len(t, res.Items(), 0)

	res, err = sqliteClient.Execute(new(Query).
		Alter(&testNewTable).
		AddForeignKey(dto.ForeignKey{
			Name: "fk_test",
			Target: query.Reference{
				Table: model.GetTableName(),
				Key:   "id",
			},
			With: query.Reference{
				Table: "test_table",
				Key:   "col2",
			},
			OnDelete: "",
			OnUpdate: "",
		}).
		DropForeignKey(dto.ForeignKey{
			Name: "fk_test2",
		}),
	)
	assert.NoError(t, err)

	removeDatabase()
}

func initTestModel(table string) dto.BaseModel {
	model := dto.BaseModel{
		TableName: table,
		Fields: []interface{}{
			dto.ModelField{
				Name:  "relation_id",
				Type:  "INTEGER",
				Value: 1,
			},
			dto.ModelField{
				Name:  "col1",
				Type:  "INTEGER",
				Value: 2,
			},
			dto.ModelField{
				Name:  "col2",
				Type:  "INTEGER",
				Value: 2,
			},
			dto.ModelField{
				Name:  "col3",
				Type:  "VARCHAR",
				Value: "Test",
			},
		},
	}
	model.SetPrimaryKey(dto.ModelField{
		Name:          "id",
		Type:          dto.IntegerColumnType,
		AutoIncrement: true,
	})
	return model
}

func TestFromInterfaceUsage(t *testing.T) {
	var actual string
	q := new(Query).Select([]interface{}{"id", "name"}).From("test_table")
	actual = SQLiteClient{}.ToSql(q)
	assert.NotEmpty(t, actual)
	assert.Equal(t, "SELECT id, name FROM test_table", actual)

	actual = MySQLClient{}.ToSql(q)
	assert.NotEmpty(t, actual)
	assert.Equal(t, "SELECT id, name FROM test_table", actual)
}

func TestSQLiteClient_Transactions(t *testing.T) {
	assert.Equal(t, "BEGIN TRANSACTION;", SQLiteClient{}.ToSql(new(Query).BeginTransaction()))
	assert.Equal(t, "COMMIT;", SQLiteClient{}.ToSql(new(Query).CommitTransaction()))
	assert.Equal(t, "ROLLBACK;", SQLiteClient{}.ToSql(new(Query).RollbackTransaction()))
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
