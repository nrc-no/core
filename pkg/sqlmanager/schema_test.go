package sqlmanager

import (
	"fmt"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/mocks"
	"github.com/nrc-no/core/pkg/sql/schema"
	"github.com/nrc-no/core/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"strings"
	"testing"
)

type Suite struct {
	suite.Suite
	db   *gorm.DB
	done func()
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func TestMain(m *testing.M) {
	// embedded-postgres runs as a separate process altogether
	// We must setup/teardown here, otherwise the process does
	// not get properly cleaned up
	done, err := testutils.TryGetPostgres()
	defer done()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	exitVal := m.Run()
	done()
	os.Exit(exitVal)
}

func (s *Suite) SetupSuite() {

	sqlDb, err := gorm.Open(postgres.Open("host=localhost port=15432 user=postgres password=postgres dbname=postgres sslmode=disable"))
	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}
	dbFactory := mocks.NewMockFactory(sqlDb)

	db, err := dbFactory.Get()
	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}

	s.db = db

}

func (s *Suite) TearDownSuite() {
	if s.done != nil {
		s.done()
	}
}

func TestSchema(t *testing.T) {
	s := sqlManager{}
	var err error

	s, err = s.handleCreateTable(sqlActionCreateTable{
		sqlTable: schema.SQLTable{
			Schema:      "schema",
			Name:        "tableName",
			Columns:     schema.SQLColumns{},
			Constraints: []schema.SQLTableConstraint{},
		},
	})

	s, err = s.handleCreateColumn(sqlActionCreateColumn{
		tableName:  "tableName",
		schemaName: "schemaName",
		sqlColumn: schema.SQLColumn{
			Name: "field",
			DataType: schema.SQLDataType{
				VarChar: &schema.SQLDataTypeVarChar{
					Length: 10,
				},
			},
		},
	})

	s, err = s.handleCreateConstraint(sqlActionCreateConstraint{
		tableName:  "tableName",
		schemaName: "schemaName",
		sqlConstraint: schema.SQLTableConstraint{
			Name: "uk_constraint",
			Unique: &schema.SQLTableConstraintUnique{
				ColumnNames: []string{"column1"},
			},
		},
	})

	assert.Equal(t, []schema.DDL{
		schema.NewDDL(`create table "schema"."tableName"();`),
		schema.NewDDL(`alter table "schemaName"."tableName" add "field" varchar(10);`),
		schema.NewDDL(`alter table "schemaName"."tableName" add constraint "uk_constraint" unique ("column1");`),
	}, s.Statements)

	assert.NoError(t, err)
}

func TestFormConversion(t *testing.T) {

	const createTableDDL = `create table "databaseId"."formId"( "id" varchar(36) primary key, "created_at" timestamp with time zone not null default NOW());`
	const formId = "formId"
	const databaseId = "databaseId"

	tests := []struct {
		name    string
		args    []*types.FormDefinition
		want    []schema.DDL
		wantErr bool
	}{
		{
			name: "empty form",
			args: []*types.FormDefinition{{
				ID:         formId,
				DatabaseID: databaseId,
			}},
			want: []schema.DDL{
				{Query: createTableDDL},
			},
		},
		{
			name: "form with text field",
			args: []*types.FormDefinition{{
				ID:         formId,
				DatabaseID: databaseId,
				Fields: []*types.FieldDefinition{
					{
						ID: "textField",
						FieldType: types.FieldType{
							Text: &types.FieldTypeText{},
						},
					},
				},
			}},
			want: []schema.DDL{
				{Query: createTableDDL},
				{Query: `alter table "databaseId"."formId" add "textField" varchar(1024);`},
			},
		},
		{
			name: "form with multiLine text field",
			args: []*types.FormDefinition{{
				ID:         formId,
				DatabaseID: databaseId,
				Fields: []*types.FieldDefinition{
					{
						ID: "multiLineTextField",
						FieldType: types.FieldType{
							MultilineText: &types.FieldTypeMultilineText{},
						},
					},
				},
			}},
			want: []schema.DDL{
				{Query: createTableDDL},
				{Query: `alter table "databaseId"."formId" add "multiLineTextField" text;`},
			},
		},
		{
			name: "form with date field",
			args: []*types.FormDefinition{{
				ID:         formId,
				DatabaseID: databaseId,
				Fields: []*types.FieldDefinition{
					{
						ID: "dateField",
						FieldType: types.FieldType{
							Date: &types.FieldTypeDate{},
						},
					},
				},
			}},
			want: []schema.DDL{
				{Query: createTableDDL},
				{Query: `alter table "databaseId"."formId" add "dateField" date;`},
			},
		},
		{
			name: "form with week field",
			args: []*types.FormDefinition{{
				ID:         formId,
				DatabaseID: databaseId,
				Fields: []*types.FieldDefinition{
					{
						ID: "weekField",
						FieldType: types.FieldType{
							Week: &types.FieldTypeWeek{},
						},
					},
				},
			}},
			want: []schema.DDL{
				{Query: createTableDDL},
				{Query: `alter table "databaseId"."formId" add "weekField" date;`},
			},
		},
		{
			name: "form with month field",
			args: []*types.FormDefinition{{
				ID:         formId,
				DatabaseID: databaseId,
				Fields: []*types.FieldDefinition{
					{
						ID: "monthField",
						FieldType: types.FieldType{
							Month: &types.FieldTypeMonth{},
						},
					},
				},
			}},
			want: []schema.DDL{
				{Query: createTableDDL},
				{Query: `alter table "databaseId"."formId" add "monthField" date;`},
			},
		},
		{
			name: "form with quantity field",
			args: []*types.FormDefinition{{
				ID:         formId,
				DatabaseID: databaseId,
				Fields: []*types.FieldDefinition{
					{
						ID: "quantityField",
						FieldType: types.FieldType{
							Quantity: &types.FieldTypeQuantity{},
						},
					},
				},
			}},
			want: []schema.DDL{
				{Query: createTableDDL},
				{Query: `alter table "databaseId"."formId" add "quantityField" int;`},
			},
		},
		{
			name: "form with subForm and key fields",
			args: []*types.FormDefinition{{
				ID:         formId,
				DatabaseID: databaseId,
				Fields: []*types.FieldDefinition{
					{
						ID: "subFormField",
						FieldType: types.FieldType{
							SubForm: &types.FieldTypeSubForm{
								Fields: []*types.FieldDefinition{
									{
										ID:  "subTextField",
										Key: true,
										FieldType: types.FieldType{
											Text: &types.FieldTypeText{},
										},
									},
								},
							},
						},
					},
				},
			}},
			want: []schema.DDL{
				{Query: createTableDDL},
				{Query: `
create table "databaseId"."subFormField"(
  "id" varchar(36) primary key,
  "created_at" timestamp with time zone not null default NOW(),
  "owner_id" varchar(36) not null references "databaseId"."formId" ("id")
);`},
				{Query: `alter table "databaseId"."subFormField" add "subTextField" varchar(1024) not null;`},
				{Query: `alter table "databaseId"."subFormField" add constraint "uk_key_subFormField" unique ("subTextField");`},
			},
		},
		{
			name: "form with reference field",
			args: []*types.FormDefinition{{
				ID:         formId,
				DatabaseID: databaseId,
				Fields: []*types.FieldDefinition{
					{
						ID: "referenceField",
						FieldType: types.FieldType{
							Reference: &types.FieldTypeReference{
								DatabaseID: "otherDatabaseId",
								FormID:     "otherFormId",
							},
						},
					},
				},
			}},
			want: []schema.DDL{
				{Query: createTableDDL},
				{Query: `alter table "databaseId"."formId" add "referenceField" varchar(36) references "otherDatabaseId"."otherFormId" ("id");`},
			},
		},
		{
			name: "form with single select field",
			args: []*types.FormDefinition{{
				ID:         formId,
				DatabaseID: databaseId,
				Fields: []*types.FieldDefinition{
					{
						ID: "singleSelectField",
						FieldType: types.FieldType{
							SingleSelect: &types.FieldTypeSingleSelect{
								Options: []*types.SelectOption{
									{
										ID:   "option1",
										Name: "Option 1",
									}, {
										ID:   "option2",
										Name: "Option 2",
									},
								},
							},
						},
					},
				},
			}},
			want: []schema.DDL{
				{Query: createTableDDL},
				{Query: `create table "databaseId"."singleSelectField"( "id" varchar(36) primary key, "name" varchar(128) not null unique);`},
				{Query: `insert into "databaseId"."singleSelectField" ("id","name") values ($1,$2);`, Args: []interface{}{"option1", "Option 1"}},
				{Query: `insert into "databaseId"."singleSelectField" ("id","name") values ($1,$2);`, Args: []interface{}{"option2", "Option 2"}},
				{Query: `alter table "databaseId"."formId" add "singleSelectField" varchar(36) references "databaseId"."formId" ("id");`},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			manager := New()
			manager, err := manager.PutForms(&types.FormDefinitionList{Items: test.args})
			if test.wantErr {
				assert.Error(t, err)
				return
			}
			if !assert.NoError(t, err) {
				return
			}
			statements := manager.GetStatements()
			for i := 0; i < len(statements); i++ {
				ddl := statements[i]
				ddl.Query = strings.ReplaceAll(ddl.Query, "\n", "")
				ddl.Query = strings.ReplaceAll(ddl.Query, "  ", " ")
				statements[i] = ddl
			}
			for i := 0; i < len(test.want); i++ {
				ddl := test.want[i]
				ddl.Query = strings.ReplaceAll(ddl.Query, "\n", "")
				ddl.Query = strings.ReplaceAll(ddl.Query, "  ", " ")
				test.want[i] = ddl
			}
			assert.Equal(t, test.want, statements)
		})
	}

}

func (s *Suite) TestSchemaActions() {

	const publicSchema = "public"
	const formId = "formId"
	const otherFormId = "other-form"
	const field1Id = "field-1"

	otherForm := &types.FormDefinition{
		ID:         otherFormId,
		DatabaseID: publicSchema,
		Fields: []*types.FieldDefinition{
			{
				ID: field1Id,
				FieldType: types.FieldType{
					Text: &types.FieldTypeText{},
				},
			},
		},
	}

	formDef := &types.FormDefinition{
		ID:         formId,
		DatabaseID: publicSchema,
		Fields: []*types.FieldDefinition{
			{
				ID:  "textField",
				Key: true,
				FieldType: types.FieldType{
					Text: &types.FieldTypeText{},
				},
			},
			{
				ID:  "dateField",
				Key: true,
				FieldType: types.FieldType{
					Date: &types.FieldTypeDate{},
				},
			},
			{
				ID: "referenceField",
				FieldType: types.FieldType{
					Reference: &types.FieldTypeReference{
						DatabaseID: publicSchema,
						FormID:     otherFormId,
					},
				},
			},
			{
				ID: "subFormField",
				FieldType: types.FieldType{
					SubForm: &types.FieldTypeSubForm{
						Fields: []*types.FieldDefinition{
							{
								ID: "sub-field",
								FieldType: types.FieldType{
									Text: &types.FieldTypeText{},
								},
							},
						},
					},
				},
			},
		},
	}

	manager := New()
	var err error
	manager, err = manager.PutForms(&types.FormDefinitionList{
		Items: []*types.FormDefinition{
			otherForm,
			formDef,
		},
	})
	if !assert.NoError(s.T(), err) {
		return
	}

	statements := manager.GetStatements()

	if err := s.db.Transaction(func(tx *gorm.DB) error {
		db, err := tx.DB()
		if err != nil {
			return err
		}
		var ddl string
		var args []interface{}
		for _, statement := range statements {
			ddl = ddl + statement.Query + "\n"
			args = append(args, statement.Args...)
		}

		s.T().Logf("executing statement: \n%s", ddl)

		if _, err := db.Exec(ddl, args...); err != nil {
			s.T().Logf("error with statement %s", ddl)
			return err
		}

		return nil

	}); !assert.NoError(s.T(), err) {
		return
	}

}
