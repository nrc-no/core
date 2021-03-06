package tests

import (
	"context"

	"github.com/nrc-no/core/pkg/api/types"
	tu "github.com/nrc-no/core/pkg/testutils"
	"github.com/stretchr/testify/assert"
)

// TestFormCreateGetList tests that we can create forms with different field kinds
func (s *Suite) TestFormCreateGetList() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var db types.Database
	if err := s.createDatabase(ctx, &db); !assert.NoError(s.T(), err) {
		return
	}

	var otherForm types.FormDefinition
	if err := s.cli.CreateForm(ctx, &types.FormDefinition{
		Name:       "Other Form",
		DatabaseID: db.ID,
		Fields: []*types.FieldDefinition{
			{
				Name: "Some Text Field",
				FieldType: types.FieldType{
					Text: &types.FieldTypeText{},
				},
			},
		},
	}, &otherForm); !assert.NoError(s.T(), err) {
		return
	}

	selectOptions := []*types.SelectOption{
		{Name: "Option 1"},
		{Name: "Option 2"},
	}

	tcs := []struct {
		name     string
		formType types.FormType
		fields   []*types.FieldDefinition
	}{
		{
			name: "With Text Field",
			fields: tu.Fields(
				tu.ATextField(tu.FieldName("My Field")),
			),
		}, {
			name: "With Required Text Field",
			fields: tu.Fields(
				tu.ATextField(tu.FieldName("My Field"), tu.FieldRequired(true)),
			),
		}, {
			name: "With Key Text Field",
			fields: tu.Fields(
				tu.ATextField(tu.FieldName("My Field"), tu.FieldRequired(true), tu.FieldKey(true)),
			),
		}, {
			name: "With Multiline Field",
			fields: tu.Fields(
				tu.AMultilineTextField(tu.FieldName("My Field")),
			),
		}, {
			name: "With Required Multiline Field",
			fields: tu.Fields(
				tu.AMultilineTextField(tu.FieldName("My Field"), tu.FieldRequired(true)),
			),
		}, {
			name: "With Quantity Field",
			fields: tu.Fields(
				tu.AQuantityField(tu.FieldName("My Field")),
			),
		}, {
			name: "With Required Quantity Field",
			fields: tu.Fields(
				tu.AQuantityField(tu.FieldName("My Field"), tu.FieldRequired(true)),
			),
		}, {
			name: "With Key Quantity Field",
			fields: tu.Fields(
				tu.AQuantityField(tu.FieldName("My Field"), tu.FieldRequired(true), tu.FieldKey(true)),
			),
		}, {
			name: "With Month Field",
			fields: tu.Fields(
				tu.AMonthField(tu.FieldName("My Field")),
			),
		}, {
			name: "With Required Month Field",
			fields: tu.Fields(
				tu.AMonthField(tu.FieldName("My Field"), tu.FieldRequired(true)),
			),
		}, {
			name: "With Key Month Field",
			fields: tu.Fields(
				tu.AMonthField(tu.FieldName("My Field"), tu.FieldRequired(true), tu.FieldKey(true)),
			),
		}, {
			name: "With Date Field",
			fields: tu.Fields(
				tu.ADateField(tu.FieldName("My Field")),
			),
		}, {
			name: "With Required Date Field",
			fields: tu.Fields(
				tu.ADateField(tu.FieldName("My Field"), tu.FieldRequired(true)),
			),
		}, {
			name: "With Key Date Field",
			fields: tu.Fields(
				tu.ADateField(tu.FieldName("My Field"), tu.FieldRequired(true), tu.FieldKey(true)),
			),
		}, {
			name: "With Week Field",
			fields: tu.Fields(
				tu.AWeekField(tu.FieldName("My Field")),
			),
		}, {
			name: "With Required Week Field",
			fields: tu.Fields(
				tu.AWeekField(tu.FieldName("My Field"), tu.FieldRequired(true)),
			),
		}, {
			name: "With Key Week Field",
			fields: tu.Fields(
				tu.AWeekField(tu.FieldName("My Field"), tu.FieldRequired(true), tu.FieldKey(true)),
			),
		}, {
			name: "With Reference Field",
			fields: tu.Fields(
				tu.AReferenceField(otherForm, tu.FieldName("My Field")),
			),
		}, {
			name: "With Required Reference Field",
			fields: tu.Fields(
				tu.AReferenceField(otherForm, tu.FieldName("My Field"), tu.FieldRequired(true)),
			),
		}, {
			name: "With Key Reference Field",
			fields: tu.Fields(
				tu.AReferenceField(otherForm, tu.FieldName("My Field"), tu.FieldRequired(true), tu.FieldKey(true)),
			),
		}, {
			name: "With Single Select Field",
			fields: tu.Fields(
				tu.ASingleSelectField(
					selectOptions,
					tu.FieldName("My Field"),
				),
			),
		}, {
			name: "With Required Single Select Field",
			fields: tu.Fields(
				tu.ASingleSelectField(
					selectOptions,
					tu.FieldName("My Field"),
					tu.FieldRequired(true),
				),
			),
		}, {
			name: "With Key Single Select Field",
			fields: tu.Fields(
				tu.ASingleSelectField(
					selectOptions,
					tu.FieldName("My Field"),
					tu.FieldKey(true),
					tu.FieldRequired(true),
				),
			),
		}, {
			name: "With Sub Form Field",
			fields: []*types.FieldDefinition{
				{
					Name: "Sub Form Field",
					FieldType: types.FieldType{
						SubForm: &types.FieldTypeSubForm{
							Fields: []*types.FieldDefinition{
								{
									Name: "Sub Text Field",
									FieldType: types.FieldType{
										Text: &types.FieldTypeText{},
									},
								},
							},
						},
					},
				},
			},
		}, {
			name: "Recipient Form",
			fields: tu.Fields(
				tu.ATextField(tu.FieldName("My Field")),
			),
			formType: types.RecipientFormType,
		},
	}

	for _, tc := range tcs {
		testCase := tc
		s.Run(testCase.name, func() {
			var form types.FormDefinition
			in := &types.FormDefinition{
				Name:       testCase.name,
				DatabaseID: db.ID,
				Type:       testCase.formType,
				Fields:     testCase.fields,
			}
			if err := s.cli.CreateForm(ctx, in, &form); !assert.NoError(s.T(), err) {
				return
			}
			var got types.FormDefinition
			if err := s.cli.GetForm(ctx, form.ID, &got); !assert.NoError(s.T(), err) {
				return
			}
			assert.Equal(s.T(), form, got)
		})
	}

	s.Run("can list forms", func() {
		var list types.FormDefinitionList
		if err := s.cli.ListForms(ctx, &list); !assert.NoError(s.T(), err) {
			return
		}
	})

}
