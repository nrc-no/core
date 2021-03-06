package types

// CaseRecipientDefinition represents the specification for allowed recipients of a CaseDefinition
// Introducing a struct here since it would be possible
// to have recipients that are not records in other forms in the future.
// For example, we might have a recipient that is a link to multiple people
// (Extensibility)
type CaseRecipientDefinition struct {
	// FormRef indicates that the recipient is a record in the form referenced
	// by the FormRef.DatabaseID and FormRef.FormID
	FormRef *FormRef `json:"formRef,omitempty"`
}

// CaseDefinition represents a special type of FormDefinition that represents a Case.
// A Case is an instance of a service given to a person, household, group, ...
// A person, household, group, ... is the recipient of the Case
type CaseDefinition struct {
	// FormDefinition represents the specification of the form for collecting
	// data about this case.
	// We embed the FormDefinition here since a Case is also a Form.
	FormDefinition FormDefinition `json:"formDefinition"`
	// RecipientDefinitions represents the types of recipient that this case allows.
	// For example, if the CaseDefinition is Colombia Individual Intake, then
	// the RecipientDefinitions could be a link to a ColombiaIndividualRecipient form.
	RecipientDefinitions []CaseRecipientDefinition `json:"recipientDefinition"`
}
