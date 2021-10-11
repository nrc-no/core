package webapp

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nrc-no/core/internal/form"
	"github.com/nrc-no/core/internal/registrationctrl"
	"github.com/nrc-no/core/internal/seeder"
	"github.com/nrc-no/core/internal/sessionmanager"
	"github.com/nrc-no/core/internal/teamstatusctrl"
	"github.com/nrc-no/core/internal/utils"
	"github.com/nrc-no/core/internal/validation"
	"github.com/nrc-no/core/pkg/cms"
	iam2 "github.com/nrc-no/core/pkg/iam"
	"golang.org/x/sync/errgroup"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type IndividualWithStatuses struct {
	Individual         *iam2.Individual
	RegistrationStatus *registrationctrl.Status
	TeamStatusActions  []teamstatusctrl.TeamStatusAction
}

func (s *Server) Individuals(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	iamClient, err := s.IAMClient(req)
	if err != nil {
		s.Error(w, err)
		return
	}

	if req.Method == "POST" {
		s.Individual(w, req)
		return
	}

	var listOptions iam2.IndividualListOptions
	if err := listOptions.UnmarshalQueryParameters(req.URL.Query()); err != nil {
		s.Error(w, err)
		return
	}
	listOptions.PartyTypeIDs = []string{iam2.BeneficiaryPartyType.ID}

	list, err := iamClient.Individuals().List(ctx, listOptions)
	if err != nil {
		s.Error(w, err)
		return
	}

	var beneficiaryWithStatuses []IndividualWithStatuses
	for _, b := range list.Items {
		rc, err := s.GetRegistrationController(w, req, b)
		if err != nil {
			s.Error(w, err)
			return
		}

		tsc, err := s.GetTeamStatusController(req, b)
		if err != nil {
			s.Error(w, err)
			return
		}

		beneficiaryWithStatuses = append(beneficiaryWithStatuses, IndividualWithStatuses{
			Individual:         b,
			RegistrationStatus: rc.Status(),
			TeamStatusActions:  tsc.GetTeamStatusActions(),
		})
	}

	if err := s.renderFactory.New(req, w).ExecuteTemplate(w, "individuals", map[string]interface{}{
		"Individuals":             list,
		"IndividualsWithStatuses": beneficiaryWithStatuses,
		"Page":                    "list",
	}); err != nil {
		s.Error(w, err)
		return
	}

}

func (s *Server) IndividualIdentificationDocuments(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	iamClient, err := s.IAMClient(req)
	if err != nil {
		s.Error(w, err)
		return
	}

	id, ok := mux.Vars(req)["id"]
	if !ok || len(id) == 0 {
		err := fmt.Errorf("no id in path")
		s.Error(w, err)
		return
	}

	var identificationDocuments *iam2.IdentificationDocumentList
	var identificationDocumentTypes *iam2.IdentificationDocumentTypeList
	var identificationDocumentTypesMap = map[string]string{}
	var individual *iam2.Individual

	g, waitCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		identificationDocumentTypes, err = iamClient.IdentificationDocumentTypes().List(waitCtx, iam2.IdentificationDocumentTypeListOptions{})
		for _, idt := range identificationDocumentTypes.Items {
			identificationDocumentTypesMap[idt.ID] = idt.Name
		}
		return err
	})

	g.Go(func() error {
		var err error
		identificationDocuments, err = iamClient.IdentificationDocuments().List(waitCtx, iam2.IdentificationDocumentListOptions{PartyIDs: []string{id}})
		return err
	})

	g.Go(func() error {
		var err error
		individual, err = iamClient.Individuals().Get(ctx, id)
		return err
	})

	if err := g.Wait(); err != nil {
		s.Error(w, err)
		return
	}

	if req.Method == "POST" {
		s.PostIndividualIdentificationDocuments(w, req, individual.ID)
		return
	}

	if req.Method == "DELETE" {
		s.DeleteIndividualIdentificationDocuments(w, req, individual.ID)
	}

	if err := s.renderFactory.New(req, w).ExecuteTemplate(w, "individual_identification_documents", map[string]interface{}{
		"Page":                           "identification_documents",
		"IsNew":                          id == "new",
		"Individual":                     individual,
		"IdentificationDocuments":        identificationDocuments,
		"IdentificationDocumentTypes":    identificationDocumentTypes,
		"IdentificationDocumentTypesMap": identificationDocumentTypesMap,
	}); err != nil {
		s.Error(w, err)
		return
	}
}

func (s *Server) PostIndividualIdentificationDocuments(w http.ResponseWriter, req *http.Request, partyID string) {

	ctx := req.Context()

	if err := req.ParseForm(); err != nil {
		s.Error(w, err)
		return
	}
	values := req.Form
	documentNumber := values.Get("documentNumber")
	documentTypeID := values.Get("documentTypeId")

	if len(documentNumber) == 0 || len(documentTypeID) == 0 {
		s.Error(w, fmt.Errorf("invalid data"))
		return
	}

	var newIdentificationDocument = &iam2.IdentificationDocument{
		PartyID:                      partyID,
		DocumentNumber:               documentNumber,
		IdentificationDocumentTypeID: documentTypeID,
	}

	iamClient, err := s.IAMClient(req)
	if err != nil {
		s.Error(w, err)
		return
	}

	if _, err := iamClient.IdentificationDocuments().Create(ctx, newIdentificationDocument); err != nil {
		s.Error(w, err)
		return
	}

	http.Redirect(w, req, "/individuals/"+partyID+"/identificationdocuments", http.StatusSeeOther)
}

func (s *Server) DeleteIndividualIdentificationDocuments(w http.ResponseWriter, req *http.Request, partyID string) {
	ctx := req.Context()

	id := req.URL.Query().Get("id")

	if len(id) == 0 {
		s.Error(w, fmt.Errorf("invalid data"))
		return
	}

	iamClient, err := s.IAMClient(req)
	if err != nil {
		s.Error(w, err)
		return
	}

	if err := iamClient.IdentificationDocuments().Delete(ctx, id); err != nil {
		s.Error(w, err)
		return
	}

	http.Redirect(w, req, "/individuals/"+partyID+"/identificationdocuments", http.StatusSeeOther)
}

func (s *Server) IndividualCredentials(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	iamClient, err := s.IAMClient(req)
	if err != nil {
		s.Error(w, err)
		return
	}

	id, ok := mux.Vars(req)["id"]
	if !ok || len(id) == 0 {
		err := fmt.Errorf("no id in path")
		s.Error(w, err)
		return
	}

	individual, err := iamClient.Individuals().Get(ctx, id)
	if err != nil {
		s.Error(w, err)
		return
	}

	if req.Method == "POST" {
		s.PostIndividualCredentials(w, req, individual.ID)
		return
	}

	if err := s.renderFactory.New(req, w).ExecuteTemplate(w, "individual_credentials", map[string]interface{}{
		"Page":       "credentials",
		"Individual": individual,
	}); err != nil {
		s.Error(w, err)
		return
	}

}

func (s *Server) PostIndividualCredentials(w http.ResponseWriter, req *http.Request, partyID string) {
	ctx := req.Context()

	if err := req.ParseForm(); err != nil {
		s.Error(w, err)
		return
	}
	values := req.Form
	password := values.Get("password")

	if err := s.login.Login().SetCredentials(ctx, partyID, password); err != nil {
		s.Error(w, err)
		return
	}

	http.Redirect(w, req, "/individuals/"+partyID+"/credentials", http.StatusSeeOther)
}

func (s *Server) Individual(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	iamClient, err := s.IAMClient(req)
	if err != nil {
		s.Error(w, err)
		return
	}

	cmsClient, err := s.CMSClient(req)
	if err != nil {
		s.Error(w, err)
		return
	}

	id, ok := mux.Vars(req)["id"]
	if (!ok || len(id) == 0) && req.Method != "POST" {
		err := fmt.Errorf("no id in path")
		s.Error(w, err)
		return
	}

	var individual *iam2.Individual
	var parties *iam2.IndividualList
	var caseTypes *cms.CaseTypeList
	var cList *cms.CaseList
	var partyTypes *iam2.PartyTypeList
	var relationshipsForIndividual *iam2.RelationshipList
	var relationshipTypes *iam2.RelationshipTypeList
	var attrs *iam2.PartyAttributeDefinitionList
	var teams *iam2.TeamList
	var situationAnalysis *cms.Case
	var individualResponse *cms.Case
	var saForm form.Form
	var irForm form.Form

	g, waitCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		if id == "new" {
			individual = iam2.NewIndividual("")
			return nil
		}
		var err error
		individual, err = iamClient.Individuals().Get(waitCtx, id)
		return err
	})

	g.Go(func() error {
		var err error
		parties, err = iamClient.Individuals().List(waitCtx, iam2.IndividualListOptions{})
		return err
	})

	g.Go(func() error {
		var err error
		teams, err = iamClient.Teams().List(waitCtx, iam2.TeamListOptions{})
		return err
	})

	g.Go(func() error {
		var err error
		partyTypes, err = iamClient.PartyTypes().List(waitCtx, iam2.PartyTypeListOptions{})
		return err
	})

	g.Go(func() error {
		if id == "new" {
			relationshipsForIndividual = &iam2.RelationshipList{
				Items: []*iam2.Relationship{},
			}
			return nil
		}
		var err error
		relationshipsForIndividual, err = iamClient.Relationships().List(waitCtx, iam2.RelationshipListOptions{EitherPartyID: id})
		return err
	})

	g.Go(func() error {
		var err error
		relationshipTypes, err = iamClient.RelationshipTypes().List(waitCtx, iam2.RelationshipTypeListOptions{
			PartyTypeID: iam2.IndividualPartyType.ID,
		})
		return err
	})

	countryID := s.GetCountryFromLoginUser(w, req)

	g.Go(func() error {
		var err error
		attrs, err = iamClient.PartyAttributeDefinitions().List(waitCtx, iam2.PartyAttributeDefinitionListOptions{
			PartyTypeIDs: []string{iam2.IndividualPartyType.ID},
			CountryIDs:   []string{iam2.GlobalCountry.ID, countryID},
		})
		return err
	})

	g.Go(func() error {
		var err error
		caseTypes, err = cmsClient.CaseTypes().List(ctx, cms.CaseTypeListOptions{})
		return err
	})

	g.Go(func() error {
		var err error
		cList, err = cmsClient.Cases().List(ctx, cms.CaseListOptions{PartyIDs: []string{id}})
		return err
	})

	g.Go(func() error {
		var err error
		returnedCases, err := cmsClient.Cases().List(ctx, cms.CaseListOptions{
			PartyIDs:    []string{id},
			CaseTypeIDs: []string{seeder.UGSituationalAnalysisCaseType.ID},
		})
		if err != nil {
			return err
		}
		if len(returnedCases.Items) == 1 {
			kase := returnedCases.Items[0]
			situationAnalysis = kase
			saForm = form.NewValidatedForm(kase.Form, kase.FormData, nil)
		}
		return err
	})

	g.Go(func() error {
		var err error
		cases, err := cmsClient.Cases().List(ctx, cms.CaseListOptions{
			PartyIDs:    []string{id},
			CaseTypeIDs: []string{seeder.UGIndividualResponseCaseType.ID},
		})
		if err != nil {
			return err
		}
		if len(cases.Items) == 1 {
			kase := cases.Items[0]
			individualResponse = kase
			irForm = form.NewValidatedForm(kase.Form, kase.FormData, nil)
		}
		return err
	})

	g.Go(func() error {
		var err error
		cList, err = cmsClient.Cases().List(ctx, cms.CaseListOptions{PartyIDs: []string{id}})
		return err
	})

	if err := g.Wait(); err != nil {
		s.Error(w, err)
		return
	}

	filteredRelationshipTypes := PrepRelationshipTypeDropdown(relationshipTypes)

	if req.Method == "POST" {
		individual, err = s.PostIndividual(ctx, attrs, id, w, req)
		if err != nil {
			if status, ok := err.(*validation.Status); ok {
				validatedAttrs := sumbittedFormFromErrors(&status.Errors, attrs, req.Form)
				s.json(w, status.Code, validatedAttrs)
			} else {
				s.Error(w, err)
			}
			return
		}
	}

	type DisplayCase struct {
		Case     *cms.Case
		CaseType *cms.CaseType
	}

	ctMap := map[string]*cms.CaseType{}
	for _, item := range caseTypes.Items {
		ctMap[item.ID] = item
	}

	var displayCases []*DisplayCase
	for _, item := range cList.Items {
		d := DisplayCase{
			item,
			ctMap[item.CaseTypeID],
		}
		displayCases = append(displayCases, &d)
	}

	rc, err := s.GetRegistrationController(w, req, individual)
	if err != nil {
		s.Error(w, err)
		return
	}
	status := rc.Status()
	progressLabel := status.Label
	progress := status.Progress

	attributes := form.Form{}
	for _, attribute := range attrs.Items {
		ctrl := attribute.FormControl
		ctrl.Value = individual.GetAttribute(attribute.ID)
		attributes.Controls = append(attributes.Controls, ctrl)
	}

	if err := s.renderFactory.New(req, w).ExecuteTemplate(w, "individual", map[string]interface{}{
		"IsNew":                     id == "new",
		"Individual":                individual,
		"Parties":                   parties,
		"Teams":                     teams,
		"PartyTypes":                partyTypes,
		"RelationshipTypes":         relationshipTypes,
		"FilteredRelationshipTypes": filteredRelationshipTypes,
		"Relationships":             relationshipsForIndividual,
		"Attributes":                attributes,
		"Cases":                     displayCases,
		"CaseTypes":                 caseTypes,
		"FullNameAttribute":         iam2.FullNameAttribute,
		"DisplayNameAttribute":      iam2.DisplayNameAttribute,
		"Page":                      "general",
		"Constants":                 s.Constants,
		"IndividualPartyTypeID":     iam2.IndividualPartyType.ID,
		"HouseholdPartyTypeID":      iam2.HouseholdPartyType.ID,
		"TeamPartyTypeID":           iam2.TeamPartyType.ID,
		"SituationAnalysis":         situationAnalysis,
		"SituationAnalysisForm":     saForm,
		"IndividualResponse":        individualResponse,
		"IndividualResponseForm":    irForm,
		"Status":                    status,
		"ProgressLabel":             progressLabel,
		"Progress":                  progress,
	}); err != nil {
		s.Error(w, err)
		return
	}
}

// sumbittedFormFromErrors returns a slice of form.Control populated with validated attributes.
func sumbittedFormFromErrors(errorList *validation.ErrorList, attributes *iam2.PartyAttributeDefinitionList, values url.Values) form.Form {
	var controls []form.Control
	for _, attribute := range attributes.Items {
		errs := errorList.FindFamily(attribute.ID)
		value := values.Get(attribute.ID)
		if len(*errs) > 0 && !shouldIgnoreValidationError(attribute, []string{value}) {
			control := attribute.FormControl
			control.Errors = errs
			controls = append(controls, control)
		}
	}
	return form.Form{Controls: controls}
}

func shouldIgnoreValidationError(attribute *iam2.PartyAttributeDefinition, values []string) bool {
	// Ignore validation errors on empty optional fields
	if !attribute.FormControl.Validation.Required && utils.AllEmpty(values) {
		return true
	}
	return false
}

func PrepRelationshipTypeDropdown(relationshipTypes *iam2.RelationshipTypeList) *iam2.RelationshipTypeList {

	// TODO:
	// Currently Core only works with individuals
	// this function should be changed when this is no longer
	// the case

	var newList = iam2.RelationshipTypeList{}
	for _, relType := range relationshipTypes.Items {

		relTypeOnlyForIndividuals := true

		for _, rule := range relType.Rules {
			if !(rule.PartyTypeRule.FirstPartyTypeID == iam2.IndividualPartyType.ID && rule.PartyTypeRule.SecondPartyTypeID == iam2.IndividualPartyType.ID) {
				relTypeOnlyForIndividuals = false
			}
		}

		if relTypeOnlyForIndividuals {
			newList.Items = append(newList.Items, relType)
			if relType.IsDirectional {
				newList.Items = append(newList.Items, relType.Mirror())
			}
		}
	}
	return &newList
}

func (s *Server) PostIndividual(ctx context.Context, attrs *iam2.PartyAttributeDefinitionList, id string, w http.ResponseWriter, req *http.Request) (*iam2.Individual, error) {

	iamClient, err := s.IAMClient(req)
	if err != nil {
		return nil, err
	}

	if err := req.ParseForm(); err != nil {
		return nil, err
	}
	values := req.Form

	var individual *iam2.Individual
	if len(id) == 0 {
		individual = iam2.NewIndividual("")
	} else {
		individual, err = iamClient.Individuals().Get(ctx, id)
		if err != nil {
			return nil, err
		}
	}

	attributeMap := map[string][]string{}
	for _, attribute := range attrs.Items {
		value := values[attribute.FormControl.Name]
		attributeMap[attribute.ID] = value
	}
	individual.Attributes = attributeMap

	type RelationshipEntry struct {
		*iam2.Relationship
		MarkedForDeletion bool
	}

	//goland:noinspection GoPreferNilSlice
	var rels = []*RelationshipEntry{}

	f := req.Form
	for attrId, value := range f {
		// Retrieve party relationships
		if strings.HasPrefix(attrId, "relationships[") {

			keyParts := strings.Split(attrId, ".")
			if len(keyParts) != 2 {
				err := fmt.Errorf("unexpected form value key: %s", attrId)
				return nil, err
			}

			if !strings.HasSuffix(keyParts[0], "]") {
				err := fmt.Errorf("unexpected form value key: %s", attrId)
				return nil, err
			}

			relationshipIdxStr := keyParts[0][14 : len(keyParts[0])-1]
			relationshipIdx, err := strconv.Atoi(relationshipIdxStr)
			if err != nil {
				return nil, err
			}

			var rel *RelationshipEntry
			for {
				if (relationshipIdx + 1) > len(rels) {
					rels = append(rels, &RelationshipEntry{
						Relationship: &iam2.Relationship{},
					})
				} else {
					rel = rels[relationshipIdx]
					break
				}
			}

			attrName := keyParts[1]

			switch attrName {
			case "markedForDeletion":
				rel.MarkedForDeletion = value[0] == "true"
			case "id":
				rel.ID = value[0]
			case "secondPartyId":
				rel.SecondPartyID = value[0]
			case "relationshipTypeId":
				rel.RelationshipTypeID = value[0]
			default:
				err := fmt.Errorf("unexpected relationship attribute: %s", attrName)
				return nil, err
			}
		}
	}

	// Update or create the individual
	var storageAction string
	if id == "" {
		individual.PartyTypeIDs = append(individual.PartyTypeIDs, iam2.BeneficiaryPartyType.ID)
		individual, err = iamClient.Individuals().Create(ctx, individual)
		if err != nil {
			return nil, err
		}
		err = s.createDefaultIndividualIntakeCases(req, individual)
		if err != nil {
			return nil, err
		}
		storageAction = "created"
	} else {
		if individual, err = iamClient.Individuals().Update(ctx, individual); err != nil {
			return nil, err
		}
		storageAction = "updated"
	}

	// Update, create or delete the relationships

	for _, rel := range rels {

		if len(rel.ID) == 0 {
			// Create the relationship
			relationship := rel.Relationship
			relationship.FirstPartyID = individual.ID
			if _, err := iamClient.Relationships().Create(ctx, rel.Relationship); err != nil {
				return nil, err
			}
		} else if !rel.MarkedForDeletion {
			// Update the relationship
			relationship := rel.Relationship
			relationship.FirstPartyID = individual.ID
			if _, err := iamClient.Relationships().Update(ctx, rel.Relationship); err != nil {
				return nil, err
			}
		} else {
			// Delete the relationship
			if err := iamClient.Relationships().Delete(ctx, rel.Relationship.ID); err != nil {
				return nil, err
			}
		}

	}

	// Set flash notification
	if err := s.sessionManager.AddNotification(req, w, &sessionmanager.Notification{
		Message: fmt.Sprintf("Individual \"%s\" successfully %s", individual.String(), storageAction),
		Theme:   "success",
	}); err != nil {
		return nil, err
	}

	if storageAction == "created" {
		w.Header().Set("Location", "/individuals/"+individual.ID)
		w.WriteHeader(http.StatusSeeOther)
	}
	return individual, nil
}

func (s *Server) createDefaultIndividualIntakeCases(req *http.Request, individual *iam2.Individual) error {
	var situationAnalysisCaseType = &seeder.UGSituationalAnalysisCaseType
	var individualResponseCaseType = &seeder.UGIndividualResponseCaseType

	cmsClient, err := s.CMSClient(req)
	if err != nil {
		return err
	}
	creatorId := ""
	ctx := req.Context()
	subject := ctx.Value("Subject")
	if subject != nil {
		creatorId = subject.(string)
	}
	// Create UG Intake Cases for new individual
	if situationAnalysisCaseType != nil {
		_, err = cmsClient.Cases().Create(ctx, &cms.Case{
			CaseTypeID: situationAnalysisCaseType.ID,
			PartyID:    individual.ID,
			Done:       false,
			TeamID:     situationAnalysisCaseType.TeamID,
			CreatorID:  creatorId,
			IntakeCase: situationAnalysisCaseType.IntakeCaseType,
			Form:       situationAnalysisCaseType.Form,
		})
		if err != nil {
			return err
		}
	}
	if individualResponseCaseType != nil {
		_, err = cmsClient.Cases().Create(ctx, &cms.Case{
			CaseTypeID: individualResponseCaseType.ID,
			PartyID:    individual.ID,
			Done:       false,
			TeamID:     individualResponseCaseType.TeamID,
			CreatorID:  creatorId,
			IntakeCase: individualResponseCaseType.IntakeCaseType,
			Form:       individualResponseCaseType.Form,
		})
		if err != nil {
			return err
		}
	}
	return nil
}