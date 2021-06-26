package webapp

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nrc-no/core-kafka/pkg/apps/iam"
	"github.com/nrc-no/core-kafka/pkg/parties/partytypes"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"strings"
)

func (h *Handler) Attributes(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	if req.Method == "POST" {
		h.PostAttribute(ctx, &iam.Attribute{}, w, req)
		return
	}

	list, err := h.iam.Attributes().List(ctx, iam.AttributeListOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	partyTypes, err := h.iam.PartyTypes().List(ctx, iam.PartyTypeListOptions{})
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

	a, err := h.iam.Attributes().Get(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if req.Method == "POST" {
		h.PostAttribute(ctx, a, w, req)
	}

	if err := h.renderFactory.New(req).ExecuteTemplate(w, "attribute", map[string]interface{}{
		"Attribute": a,
		"PartyTypes": iam.PartyTypeList{
			Items: []*iam.PartyType{
				&iam.IndividualPartyType,
			},
		},
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (h *Handler) PostAttribute(ctx context.Context, attribute *iam.Attribute, w http.ResponseWriter, req *http.Request) {

	if err := req.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	values := req.Form

	translationMap := map[string]*iam.AttributeTranslation{}

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
			translationMap[locale] = &iam.AttributeTranslation{
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

	var translations []iam.AttributeTranslation
	for _, translation := range translationMap {
		translations = append(translations, *translation)
	}

	isNew := false
	if len(attribute.ID) == 0 {
		attribute.ID = uuid.NewV4().String()
		isNew = true
	}

	attribute.Name = values.Get("name")
	attribute.PartyTypeIDs = values["partyTypes"]
	attribute.Translations = translations
	attribute.IsPersonallyIdentifiableInfo = values.Get("isPii") == "on"

	var out *iam.Attribute

	if isNew {
		var err error
		out, err = h.iam.Attributes().Create(ctx, attribute)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		var err error
		out, err = h.iam.Attributes().Update(ctx, attribute)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Location", "/settings/attributes/"+out.ID)
	w.WriteHeader(http.StatusSeeOther)

}
