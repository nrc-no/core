package types

// RecordRef represents a key that allows referencing a single Record.
type RecordRef struct {
	// ID is the Record.ID
	ID string `json:"id"`
	// DatabaseID is the Record.DatabaseID
	DatabaseID string `json:"databaseId"`
	// FormID is the Record.FormID
	FormID string `json:"formId"`
}

// FormRef represents a key that allows referencing a single FormDefinition.
type FormRef struct {
	// DatabaseID represents the FormDefinition.DatabaseID
	DatabaseID string `json:"databaseId"`
	// FormID represents the FormDefinition.ID
	FormID string `json:"formId"`
}

// GetDatabaseID implements FormReference.GetDatabaseID
func (f FormRef) GetDatabaseID() string {
	return f.DatabaseID
}

// GetFormID implements FormReference.GetFormID
func (f FormRef) GetFormID() string {
	return f.FormID
}

// Record represents an entry in a Form.
type Record struct {
	// ID of the record
	ID string `json:"id"`
	// Seq of the Record. This value is automatically increased by the database.
	// The presence of this field allows us to sort the table by insertion order.
	Seq int64 `json:"seq"`
	// DatabaseID of the Record
	DatabaseID string `json:"databaseId"`
	// FormID of the Record. Represents in which Form this record belongs.
	FormID string `json:"formId"`
	// OwnerID represents the owner of the Record. In cases where
	// a Record is part of a SubForm, this field records the "Owner" form ID.
	OwnerID *string `json:"ownerId"`
	// Values is an arbitrary map of values that correspond to the FormDefinition.Fields.
	// The key of the map is the FieldDefinition.ID ! (not the FormDefinition.Name, not
	// the FormDefinition.Code)
	Values map[string]interface{} `json:"values"`
}

// RecordList represents a list of Record
type RecordList struct {
	Items []*Record `json:"items"`
}

// RecordListOptions represents the options for listing Record.
type RecordListOptions struct {
	DatabaseID string
	FormID     string
}
