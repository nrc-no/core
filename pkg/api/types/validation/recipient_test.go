package validation

import (
	"testing"

	uuid "github.com/satori/go.uuid"

	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/validation"
	"github.com/stretchr/testify/assert"
)

func Test_ValidateRecipientForm_SkippedWhenNotRecipientForm(t *testing.T) {

	validFormButInvalidRecipientForm := &types.FormDefinition{
		Type:       types.DefaultFormType,
		Name:       "test",
		DatabaseID: uuid.NewV4().String(),
		Fields: []*types.FieldDefinition{
			{
				Key:      true,
				Name:     "bla",
				Required: true,
				FieldType: types.FieldType{
					Text: &types.FieldTypeText{},
				},
			},
			{
				Key:      true,
				Required: true,
				Name:     "bli",
				FieldType: types.FieldType{
					Text: &types.FieldTypeText{},
				},
			},
		},
	}

	// assert that the form is valid
	assert.Equal(t, validation.ErrorList(nil), ValidateForm(validFormButInvalidRecipientForm))

	// change the type to "Recipient"
	validFormButInvalidRecipientForm.Type = types.RecipientFormType

	// assert that the form is now invalid
	assert.False(t, ValidateForm(validFormButInvalidRecipientForm).IsEmpty())
}

func Test_validateRecipientFormHasSingleKeyField(t *testing.T) {
	tests := []struct {
		name           string
		fields         types.FieldDefinitions
		path           *validation.Path
		wantField      *types.FieldDefinition
		wantFieldIndex int
		wantErrs       validation.ErrorList
	}{
		{
			name:           "without key field",
			fields:         []*types.FieldDefinition{},
			path:           validation.NewPath("fields"),
			wantField:      nil,
			wantFieldIndex: -1,
			wantErrs:       nil,
		},
		{
			name: "with single key field",
			fields: []*types.FieldDefinition{
				{
					Key: true,
					FieldType: types.FieldType{
						Text: &types.FieldTypeText{},
					},
				},
			},
			path: validation.NewPath("fields"),
			wantField: &types.FieldDefinition{
				Key: true,
				FieldType: types.FieldType{
					Text: &types.FieldTypeText{},
				},
			},
			wantFieldIndex: 0,
			wantErrs:       nil,
		},
		{
			name: "with multiple key field",
			fields: []*types.FieldDefinition{
				{
					Key: true,
					FieldType: types.FieldType{
						Text: &types.FieldTypeText{},
					},
				},
				{
					Key: true,
					FieldType: types.FieldType{
						Text: &types.FieldTypeText{},
					},
				},
			},
			path:           validation.NewPath("fields"),
			wantField:      nil,
			wantFieldIndex: -1,
			wantErrs: validation.ErrorList{
				validation.Invalid(
					validation.NewPath("fields").Index(1).Child("key"),
					true,
					errRecipientMultipleKeyFields,
				),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := validateRecipientFormHasSingleKeyField(tt.fields, tt.path)
			assert.Equalf(t, tt.wantField, got, "validateRecipientFormHasSingleKeyField(%v, %v)", tt.fields, tt.path)
			assert.Equalf(t, tt.wantFieldIndex, got1, "validateRecipientFormHasSingleKeyField(%v, %v)", tt.fields, tt.path)
			assert.Equalf(t, tt.wantErrs, got2, "validateRecipientFormHasSingleKeyField(%v, %v)", tt.fields, tt.path)
		})
	}
}
