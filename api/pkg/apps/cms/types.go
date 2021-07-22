package cms

import (
	"encoding/json"
	"github.com/nrc-no/core/pkg/validation"
	"html/template"
	"time"
)

type Case struct {
	ID         string        `json:"id" bson:"id"`
	CaseTypeID string        `json:"caseTypeId" bson:"caseTypeId"`
	PartyID    string        `json:"partyId" bson:"partyId"`
	Done       bool          `json:"done" bson:"done"`
	ParentID   string        `json:"parentId" bson:"parentId"`
	TeamID     string        `json:"teamId" bson:"teamId"`
	CreatorID  string        `json:"creatorId" bson:"creatorId"`
	FormData   *CaseTemplate `json:"formData" bson:"formData"`
}

type CaseList struct {
	Items []*Case `json:"items" bson:"items"`
}

type CaseType struct {
	ID          string        `json:"id" bson:"id"`
	Name        string        `json:"name" bson:"name"`
	PartyTypeID string        `json:"partyTypeId" bson:"partyTypeId"`
	TeamID      string        `json:"teamId" bson:"teamId"`
	Template    *CaseTemplate `json:"template" bson:"template"`
}

type CaseTypeList struct {
	Items []*CaseType `json:"items" bson:"items"`
}

func (c *CaseType) String() string {
	return c.Name
}

func (c *CaseType) Pretty() string {
	b, err := json.MarshalIndent(c.Template, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(b)
}

func (l *CaseTypeList) FindByID(id string) *CaseType {
	for _, caseType := range l.Items {
		if caseType.ID == id {
			return caseType
		}
	}
	return nil
}

type Comment struct {
	ID        string    `json:"id" bson:"id"`
	CaseID    string    `json:"caseId" bson:"caseId"`
	AuthorID  string    `json:"authorId" bson:"authorId"`
	Body      string    `json:"body" bson:"body"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type CommentList struct {
	Items []*Comment `json:"items"`
}

// Case templates
// https://docs.github.com/en/communities/using-templates-to-encourage-useful-issues-and-pull-requests/syntax-for-githubs-form-schema

type CaseTemplate struct {
	// FormElements is an ordered list of the elements found in the form
	FormElements []FormElement `json:"formElements" bson:"formElements"`
}

type FormElementType string

const (
	Dropdown  FormElementType = "dropdown"
	Textarea  FormElementType = "textarea"
	TextInput FormElementType = "textinput"
	Checkbox  FormElementType = "checkbox"
)

type FormElement struct {
	Type       FormElementType       `json:"type" bson:"type"`
	Attributes FormElementAttribute  `json:"attributes" bson:"attributes"`
	Validation FormElementValidation `json:"validation" bson:"validation"`
}

type FormElementAttribute struct {
	Label           string           `json:"label" bson:"label"`
	ID              string           `json:"id" bson:"id"`
	Description     string           `json:"description" bson:"description"`
	Placeholder     string           `json:"placeholder" bson:"placeholder"`
	Value           []string         `json:"value" bson:"value"`
	Multiple        bool             `json:"multiple" bson:"multiple"`
	Options         []string         `json:"options" bson:"options"`
	CheckboxOptions []CheckboxOption `json:"checkboxOptions" bson:"checkboxOptions"`
}

type FormElementValidation struct {
	Required bool `json:"required" bson:"required"`
}

type CheckboxOption struct {
	Label    string `json:"label" bson:"label"`
	Required bool   `json:"required" bson:"required"`
}

type CaseForm struct {
	Elements []CaseFormElement
}

type CaseFormElement struct {
	FormElement
	Error validation.Error
}

func (c CaseFormElement) RenderError() template.HTML {
	return template.HTML("")
}
