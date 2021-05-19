package controllers

import (
	v1 "github.com/nrc-no/core/api/pkg/apis/core/v1"
	"github.com/stretchr/testify/assert"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func int64ptr(value int64) *int64 {
	return &value
}

func TestWalkFormSchemeShouldFlattenProperties(t *testing.T) {
	element := v1.FormElementDefinition{
		Type: v1.SectionType,
		Children: []v1.FormElementDefinition{
			{
				Key:       "prop1",
				Type:      v1.ShortTextType,
				MinLength: 3,
			}, {
				Key:       "prop2",
				Type:      v1.IntegerType,
				MaxLength: int64ptr(10),
			}, {
				Type: v1.SectionType,
				Children: []v1.FormElementDefinition{
					{
						Key:  "",
						Type: v1.SectionType,
						Children: []v1.FormElementDefinition{
							{
								Key:  "prop3",
								Type: v1.ShortTextType,
							},
						},
					},
				},
			},
		},
	}
	jsonProps := apiextensionsv1.JSONSchemaProps{}
	walkFormSchema(element, &jsonProps)

	assert.Equal(t, 3, len(jsonProps.Properties))
	assertHasProperty(t, jsonProps, "prop1")
	assertHasProperty(t, jsonProps, "prop2")
	assertHasProperty(t, jsonProps, "prop3")

	return
}

func TestCreateCrdFromFormDefinition(t *testing.T) {

	formDef := &v1.FormDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: "some-form",
		},
		Spec: v1.FormDefinitionSpec{
			Group: "supergroup",
			Names: v1.FormDefinitionNames{
				Plural:   "superforms",
				Singular: "superform",
				Kind:     "SuperForm",
			},
			Versions: []v1.FormDefinitionVersion{
				{
					Name: "v1",
					Schema: v1.FormDefinitionValidation{
						FormSchema: v1.FormDefinitionSchema{
							Root: v1.FormElementDefinition{
								Children: []v1.FormElementDefinition{
									{
										Key:      "prop1",
										Type:     v1.LongTextType,
										Required: true,
										Description: v1.TranslatedStrings{
											{
												Locale: "en",
												Value:  "Property 1",
											},
										},
									}, {
										Key:  "prop2",
										Type: v1.LongTextType,
										Description: v1.TranslatedStrings{
											{
												Locale: "en",
												Value:  "Property 2",
											},
										},
									}, {
										Key:  "prop3",
										Type: v1.LongTextType,
										Description: v1.TranslatedStrings{
											{
												Locale: "en",
												Value:  "Property 3",
											},
										},
									},
								},
							},
						},
					},
				}, {
					Name: "v2",
					Schema: v1.FormDefinitionValidation{
						FormSchema: v1.FormDefinitionSchema{
							Root: v1.FormElementDefinition{
								Children: []v1.FormElementDefinition{
									{
										Key:      "prop1",
										Type:     v1.LongTextType,
										Required: true,
										Description: v1.TranslatedStrings{
											{
												Locale: "en",
												Value:  "Property 1",
											},
										},
									}, {
										Key:  "prop2",
										Type: v1.LongTextType,
										Description: v1.TranslatedStrings{
											{
												Locale: "en",
												Value:  "Property 2",
											},
										},
									}, {
										Key:  "prop3",
										Type: v1.LongTextType,
										Description: v1.TranslatedStrings{
											{
												Locale: "en",
												Value:  "Property 3",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	expected := &apiextensionsv1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: "some-form",
		},
		Spec: apiextensionsv1.CustomResourceDefinitionSpec{
			Group: "supergroup",
			Names: apiextensionsv1.CustomResourceDefinitionNames{
				Singular: "superform",
				Plural:   "superforms",
				Kind:     "SuperForm",
			},
			Versions: []apiextensionsv1.CustomResourceDefinitionVersion{
				{
					Name: "v1",
					Schema: &apiextensionsv1.CustomResourceValidation{
						OpenAPIV3Schema: &apiextensionsv1.JSONSchemaProps{
							Description: "Schema for the SuperForm api",
							Type:        "object",
							Properties: map[string]apiextensionsv1.JSONSchemaProps{
								"apiVersion": {
									Type: "string",
									Description: `APIVersion defines the versioned schema of this representation
of an object. Servers should convert recognized schemas to the latest internal value, and may
reject unrecognized values.`,
								},
								"kind": {
									Type: "string",
									Description: `Kind is a string value representing the REST resource this 
object represents. Servers may infer this from the endpoint the client submits requests to.
Cannot be updated. In CamelCase.`,
								},
								"metadata": {
									Type: "object",
								},
								"spec": {
									Type:        "object",
									Description: "Defines the desired state fo SuperForm",
									Required: []string{
										"prop1",
									},
									Properties: map[string]apiextensionsv1.JSONSchemaProps{
										"prop1": {
											Type:        "string",
											Description: "Property 1",
										},
										"prop2": {
											Type:        "string",
											Description: "Property 2",
										},
										"prop3": {
											Type:        "string",
											Description: "Property 3",
										},
									},
								},
							},
						},
					},
				},
				{
					Name: "v2",
					Schema: &apiextensionsv1.CustomResourceValidation{
						OpenAPIV3Schema: &apiextensionsv1.JSONSchemaProps{
							Description: "Schema for the SuperForm api",
							Type:        "object",
							Properties: map[string]apiextensionsv1.JSONSchemaProps{
								"apiVersion": {
									Type: "string",
									Description: `APIVersion defines the versioned schema of this representation
of an object. Servers should convert recognized schemas to the latest internal value, and may
reject unrecognized values.`,
								},
								"kind": {
									Type: "string",
									Description: `Kind is a string value representing the REST resource this 
object represents. Servers may infer this from the endpoint the client submits requests to.
Cannot be updated. In CamelCase.`,
								},
								"metadata": {
									Type: "object",
								},
								"spec": {
									Type:        "object",
									Description: "Defines the desired state fo SuperForm",
									Required: []string{
										"prop1",
									},
									Properties: map[string]apiextensionsv1.JSONSchemaProps{
										"prop1": {
											Type:        "string",
											Description: "Property 1",
										},
										"prop2": {
											Type:        "string",
											Description: "Property 2",
										},
										"prop3": {
											Type:        "string",
											Description: "Property 3",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	crd := createCrdFromFormDefinition(formDef)

	assert.Equal(t, expected, crd)

	return

}

func assertHasProperty(t *testing.T, schema apiextensionsv1.JSONSchemaProps, key string) {
	if !assert.NotNil(t, schema.Properties) {
		return
	}
	_, ok := schema.Properties[key]
	assert.True(t, ok)
}
