package tests

import (
	"context"
	"github.com/nrc-no/core/pkg/api/types"
	tu "github.com/nrc-no/core/pkg/testutils"
	"github.com/stretchr/testify/assert"
)

func (s *Suite) TestRecordCreateReadList() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var db types.Database
	if err := s.createDatabase(ctx, &db); !assert.NoError(s.T(), err) {
		return
	}

	var form types.FormDefinition
	const (
		textFieldName                 = "Text Field"
		monthFieldName                = "Month Field"
		dateFieldName                 = "Date Field"
		weekFieldName                 = "Week Field"
		multilineTextFieldName        = "Multiline Text Field"
		quantityFieldName             = "Quantity Field"
		singleSelectFieldName         = "Single Select"
		singleSelectNullFieldName     = "Single Select 2"
		singleSelectRequiredFieldName = "Single Select 3"
		multiSelectFieldName          = "Multi Select"
	)

	if err := s.cli.CreateForm(ctx, &types.FormDefinition{
		Name:       "My Form",
		DatabaseID: db.ID,
		Fields: []*types.FieldDefinition{
			tu.ATextField(tu.FieldName(textFieldName)),
			tu.AMonthField(tu.FieldName(monthFieldName)),
			tu.ADateField(tu.FieldName(dateFieldName)),
			tu.AWeekField(tu.FieldName(weekFieldName)),
			tu.AMultilineTextField(tu.FieldName(multilineTextFieldName)),
			tu.AQuantityField(tu.FieldName(quantityFieldName)),
			tu.ASingleSelectField([]*types.SelectOption{
				{Name: "option1"},
				{Name: "option2"},
			}, tu.FieldName(singleSelectFieldName)),

			tu.ASingleSelectField([]*types.SelectOption{
				{Name: "option1"},
				{Name: "option2"},
			},
				tu.FieldName(singleSelectNullFieldName),
			),
			tu.ASingleSelectField([]*types.SelectOption{
				{Name: "option1"},
				{Name: "option2"},
			},
				tu.FieldName(singleSelectRequiredFieldName),
			),
			tu.AMultiSelectField([]*types.SelectOption{
				{Name: "option1"},
				{Name: "option2"},
				{Name: "option3"},
			}, tu.FieldName(multiSelectFieldName)),
		},
	}, &form); !assert.NoError(s.T(), err) {
		return
	}

	var values types.FieldValues
	values, _ = values.SetValueForFieldName(&form, textFieldName, types.NewStringValue("text value"))
	values, _ = values.SetValueForFieldName(&form, monthFieldName, types.NewStringValue("2010-01"))
	values, _ = values.SetValueForFieldName(&form, dateFieldName, types.NewStringValue("2010-12-31"))
	values, _ = values.SetValueForFieldName(&form, weekFieldName, types.NewStringValue("2020-W10"))
	values, _ = values.SetValueForFieldName(&form, multilineTextFieldName, types.NewStringValue("text\nvalue"))
	values, _ = values.SetValueForFieldName(&form, quantityFieldName, types.NewStringValue("10"))

	singleSelectFieldSimple, _ := form.Fields.GetFieldByName(singleSelectFieldName)
	values, _ = values.SetValueForFieldName(&form, singleSelectFieldName, types.NewStringValue(singleSelectFieldSimple.FieldType.SingleSelect.Options[0].ID))

	singleSelectFieldRequired, _ := form.Fields.GetFieldByName(singleSelectRequiredFieldName)
	values, _ = values.SetValueForFieldName(&form, singleSelectRequiredFieldName, types.NewStringValue(
		singleSelectFieldRequired.FieldType.SingleSelect.Options[1].ID))

	values, _ = values.SetValueForFieldName(&form, singleSelectNullFieldName, types.NewNullValue())

	multiSelectField, _ := form.Fields.GetFieldByName(multiSelectFieldName)
	values, _ = values.SetValueForFieldName(&form, multiSelectFieldName, types.NewArrayValue([]string{
		multiSelectField.FieldType.MultiSelect.Options[0].ID,
		multiSelectField.FieldType.MultiSelect.Options[1].ID,
	}))

	var record types.Record
	if err := s.cli.CreateRecord(ctx, &types.Record{
		Values:     values,
		FormID:     form.ID,
		DatabaseID: form.DatabaseID,
	}, &record); !assert.NoError(s.T(), err) {
		return
	}

	var got types.Record
	if err := s.cli.GetRecord(ctx, types.RecordRef{
		ID:         record.ID,
		DatabaseID: record.DatabaseID,
		FormID:     record.FormID,
	}, &got); !assert.NoError(s.T(), err) {
		return
	}
	assert.Equal(s.T(), record, got)

	var list types.RecordList
	if err := s.cli.ListRecords(ctx, types.RecordListOptions{
		DatabaseID: form.DatabaseID,
		FormID:     form.ID,
	}, &list); !assert.NoError(s.T(), err) {
		return
	}
	assert.Contains(s.T(), list.Items, &got)

}
