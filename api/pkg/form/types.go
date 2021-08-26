package form

import "github.com/nrc-no/core/pkg/validation"

type FieldType string

const (
	Text          FieldType = "text"
	Email         FieldType = "email"
	Phone         FieldType = "tel"
	URL           FieldType = "url"
	Date          FieldType = "date"
	Textarea      FieldType = "textarea"
	Dropdown      FieldType = "dropdown"
	Checkbox      FieldType = "checkbox"
	Radio         FieldType = "radio"
	TaxonomyInput FieldType = "taxonomyinput"
	File          FieldType = "file"
	CustomDiv     FieldType = "div"
)

type FormElement struct {
	Type       FieldType             `json:"type" bson:"type"`
	Attributes FormElementAttributes `json:"attributes" bson:"attributes"`
	Validation FormElementValidation `json:"validation" bson:"validation"`
	Errors     *validation.ErrorList `json:"errors"`
	Readonly   bool
}

type I18nString struct {
	Locale string `json:"locale" bson:"locale"`
	Value  string `json:"value" bson:"value"`
}

// TODO COR-209 change static string fields (labels, descriptors) to []I18nString?
type FormElementAttributes struct {
	Label           string           `json:"label" bson:"label"`
	Name            string           `json:"name" bson:"name"`
	Value           []string         `json:"value" bson:"value"`
	Description     string           `json:"description" bson:"description"`
	Placeholder     string           `json:"placeholder" bson:"placeholder"`
	Multiple        bool             `json:"multiple" bson:"multiple"`
	Options         []string         `json:"options" bson:"options"`
	CheckboxOptions []CheckboxOption `json:"checkboxOptions" bson:"checkboxOptions"`
}

type CheckboxOption struct {
	Label    string `json:"label" bson:"label"`
	Required bool   `json:"required" bson:"required"`
}

type FormElementValidation struct {
	Required bool `json:"required" bson:"required"`
}
