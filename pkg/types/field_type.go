package types

type FieldType struct {
	Text      *FieldTypeText      `json:"text,omitempty"`
	Reference *FieldTypeReference `json:"reference,omitempty"`
	SubForm   *FieldTypeSubForm   `json:"subForm,omitempty"`
	Date      *FieldTypeDate      `json:"date,omitempty"`
}

type FieldTypeReference struct {
	DatabaseID string `json:"databaseId,omitempty"`
	FormID     string `json:"formId,omitempty"`
}

type FieldTypeText struct{}

type FieldTypeMultiLineText struct{}

type FieldTypeDate struct{}

type FieldTypeSubForm struct {
	ID     string             `json:"id"`
	Name   string             `json:"name"`
	Code   string             `json:"code"`
	Fields []*FieldDefinition `json:"fields,omitempty"`
}

func (f *FieldTypeSubForm) GetID() string {
	return f.ID
}
func (f *FieldTypeSubForm) GetFields() []*FieldDefinition {
	return f.Fields
}