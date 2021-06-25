package webapp

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nrc-no/core-kafka/pkg/expressions"
	"github.com/nrc-no/core-kafka/pkg/parties/attributes"
	"github.com/nrc-no/core-kafka/pkg/parties/partytypes"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"strings"
)

func (h *Handler) Attributes(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	if req.Method == "POST" {
		h.PostAttribute(ctx, &attributes.Attribute{}, w, req)
		return
	}

	list, err := h.attributeClient.List(ctx, attributes.ListOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	partyTypes, err := h.partyTypeClient.List(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.renderFactory.New(req).ExecuteTemplate(w, "attributes", map[string]interface{}{
		"Attributes": list,
		"PartyTypes": partyTypes,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) NewAttribute(w http.ResponseWriter, req *http.Request) {
	if err := h.renderFactory.New(req).ExecuteTemplate(w, "attribute", map[string]interface{}{
		"PartyTypes": partytypes.PartyTypeList{
			Items: []*partytypes.PartyType{
				&partytypes.IndividualPartyType,
			},
		},
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Attribute(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	id, ok := mux.Vars(req)["id"]
	if !ok || len(id) == 0 {
		err := fmt.Errorf("No id in path")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	a, err := h.attributeClient.Get(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if req.Method == "POST" {
		h.PostAttribute(ctx, a, w, req)
	}

	if err := h.renderFactory.New(req).ExecuteTemplate(w, "attribute", map[string]interface{}{
		"Attribute": a,
		"PartyTypes": partytypes.PartyTypeList{
			Items: []*partytypes.PartyType{
				&partytypes.IndividualPartyType,
			},
		},
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (h *Handler) PostAttribute(ctx context.Context, attribute *attributes.Attribute, w http.ResponseWriter, req *http.Request) {

	if err := req.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	values := req.Form

	translationMap := map[string]*attributes.AttributeTranslation{}

	for key, v := range values {
		if !strings.HasPrefix(key, "translations.") {
			continue
		}
		parts := strings.Split(key, ".")
		if len(parts) != 3 {
			http.Error(w, "unexpected translation key. Expected 'translation.{locale}.{short/long}' format", http.StatusInternalServerError)
			return
		}

		locale := parts[1]
		part := parts[2]

		if part != "long" && part != "short" {
			http.Error(w, "unexpected translation key. Expected 'translation.{locale}.{short/long}' format", http.StatusInternalServerError)
			return
		}

		if _, ok := translationMap[locale]; !ok {
			translationMap[locale] = &attributes.AttributeTranslation{
				Locale: locale,
			}
		}
		t := translationMap[locale]

		if part == "long" {
			t.LongFormulation = v[0]
		} else if part == "short" {
			t.ShortFormulation = v[0]
		} else {
			http.Error(w, "unexpected translation key. Expected 'translation.{locale}.{short/long}' format", http.StatusInternalServerError)
			return
		}

	}

	var translations []attributes.AttributeTranslation
	for _, translation := range translationMap {
		translations = append(translations, *translation)
	}

	isNew := false
	if len(attribute.ID) == 0 {
		attribute.ID = uuid.NewV4().String()
		isNew = true
	}

	attribute.Name = values.Get("name")
	attribute.ValueType = expressions.ValueType{}
	attribute.PartyTypeIDs = values["partyTypes"]
	attribute.Translations = translations
	attribute.IsPersonallyIdentifiableInfo = values.Get("isPii") == "on"

	var out *attributes.Attribute

	if isNew {
		var err error
		out, err = h.attributeClient.Create(ctx, attribute)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		var err error
		out, err = h.attributeClient.Update(ctx, attribute)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Location", "/settings/attributes/"+out.ID)
	w.WriteHeader(http.StatusSeeOther)

}
