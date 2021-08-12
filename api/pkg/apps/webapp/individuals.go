package webapp

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nrc-no/core/pkg/apps/cms"
	"github.com/nrc-no/core/pkg/apps/iam"
	"github.com/nrc-no/core/pkg/apps/seeder"
	"github.com/nrc-no/core/pkg/registrationctrl"
	"github.com/nrc-no/core/pkg/sessionmanager"
	"github.com/nrc-no/core/pkg/utils"
	"github.com/nrc-no/core/pkg/validation"
	"golang.org/x/sync/errgroup"
	"net/http"
	"strconv"
	"strings"
)

type IndividualWithStatuses struct {
	Individual         *iam.Individual
	RegistrationStatus *registrationctrl.Status
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

	var listOptions iam.IndividualListOptions
	if err := listOptions.UnmarshalQueryParameters(req.URL.Query()); err != nil {
		s.Error(w, err)
		return
	}

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
		beneficiaryWithStatuses = append(beneficiaryWithStatuses, IndividualWithStatuses{
			Individual:         b,
			RegistrationStatus: rc.Status(),
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

	var b *iam.Individual
	var bList *iam.IndividualList
	var ctList *cms.CaseTypeList
	var cList *cms.CaseList
	var partyTypes *iam.PartyTypeList
	var relationshipsForIndividual *iam.RelationshipList
	var relationshipTypes *iam.RelationshipTypeList
	var attrs *iam.AttributeList
	var tList *iam.TeamList
	var individualAssessment *cms.Case
	var situationAnalysis *cms.Case

	g, waitCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		if id == "new" {
			b = iam.NewIndividual("")
			return nil
		}
		var err error
		b, err = iamClient.Individuals().Get(waitCtx, id)
		return err
	})

	g.Go(func() error {
		var err error
		bList, err = iamClient.Individuals().List(waitCtx, iam.IndividualListOptions{})
		return err
	})

	g.Go(func() error {
		var err error
		tList, err = iamClient.Teams().List(waitCtx, iam.TeamListOptions{})
		return err
	})

	g.Go(func() error {
		var err error
		partyTypes, err = iamClient.PartyTypes().List(waitCtx, iam.PartyTypeListOptions{})
		return err
	})

	g.Go(func() error {
		if id == "new" {
			relationshipsForIndividual = &iam.RelationshipList{
				Items: []*iam.Relationship{},
			}
			return nil
		}
		var err error
		relationshipsForIndividual, err = iamClient.Relationships().List(waitCtx, iam.RelationshipListOptions{EitherPartyID: id})
		return err
	})

	g.Go(func() error {
		var err error
		relationshipTypes, err = iamClient.RelationshipTypes().List(waitCtx, iam.RelationshipTypeListOptions{
			PartyTypeID: iam.IndividualPartyType.ID,
		})
		return err
	})

	g.Go(func() error {
		var err error
		attrs, err = iamClient.Attributes().List(waitCtx, iam.AttributeListOptions{
			PartyTypeIDs: []string{iam.IndividualPartyType.ID},
		})
		return err
	})

	g.Go(func() error {
		var err error
		ctList, err = cmsClient.CaseTypes().List(ctx, cms.CaseTypeListOptions{})
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
			CaseTypeIDs: []string{seeder.UGIndividualResponseCaseType.ID},
		})
		if err != nil {
			return err
		}
		if len(returnedCases.Items) == 1 {
			individualAssessment = returnedCases.Items[0]
		}
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
			situationAnalysis = returnedCases.Items[0]
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
		err := s.PostIndividual(ctx, attrs, id, w, req)
		if err != nil {
			if status, ok := err.(*validation.Status); ok {
				zipErrorsAndAttributes(&status.Errors, attrs)
				s.json(w, status.Code, attrs)
				return
			} else {
				s.Error(w, err)
				return
			}
		} else {
			return
		}
	}

	type DisplayCase struct {
		Case     *cms.Case
		CaseType *cms.CaseType
	}

	ctMap := map[string]*cms.CaseType{}
	for _, item := range ctList.Items {
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

	rc, err := s.GetRegistrationController(w, req, b)
	if err != nil {
		s.Error(w, err)
		return
	}
	status := rc.Status()
	progressLabel := status.Label
	progress := status.Progress

	// Write Individual attribute values (if any) to Attributes object
	for _, attribute := range attrs.Items {
		values := b.GetAttribute(attribute.ID)
		attribute.Attributes.Value = values
	}

	if err := s.renderFactory.New(req, w).ExecuteTemplate(w, "individual", map[string]interface{}{
		"IsNew":                     id == "new",
		"Individual":                b,
		"Parties":                   bList,
		"Teams":                     tList,
		"PartyTypes":                partyTypes,
		"RelationshipTypes":         relationshipTypes,
		"FilteredRelationshipTypes": filteredRelationshipTypes,
		"Relationships":             relationshipsForIndividual,
		"Attributes":                attrs,
		"Cases":                     displayCases,
		"CaseTypes":                 ctList,
		"FirstNameAttribute":        iam.FirstNameAttribute,
		"LastNameAttribute":         iam.LastNameAttribute,
		"Page":                      "general",
		"Constants":                 s.Constants,
		"IndividualPartyTypeID":     iam.IndividualPartyType.ID,
		"HouseholdPartyTypeID":      iam.HouseholdPartyType.ID,
		"TeamPartyTypeID":           iam.TeamPartyType.ID,
		"IndividualAssessment":      individualAssessment,
		"SituationAnalysis":         situationAnalysis,
		"ProgressLabel":             progressLabel,
		"Progress":                  progress,
	}); err != nil {
		s.Error(w, err)
		return
	}

}

// zipErrorsAndAttributes populates the Attribute Errors field if any matching errors are found
func zipErrorsAndAttributes(errorList *validation.ErrorList, attributes *iam.AttributeList) {
	for _, attribute := range attributes.Items {
		errs := errorList.Find(attribute.Attributes.Name)
		if errs != nil && !shouldIgnoreValidationError(attribute, errs) {
			attribute.Errors = errs
		}
	}
}

func shouldIgnoreValidationError(attribute *iam.Attribute, errorList *validation.ErrorList) bool {
	// Ignore validation errors on empty optional fields
	if !attribute.Validation.Required && utils.AllEmpty(attribute.Attributes.Value) {
		return true
	}
	return false
}

func PrepRelationshipTypeDropdown(relationshipTypes *iam.RelationshipTypeList) *iam.RelationshipTypeList {

	// TODO:
	// Currently Core only works with individuals
	// this function should be changed when this is no longer
	// the case

	var newList = iam.RelationshipTypeList{}
	for _, relType := range relationshipTypes.Items {

		relTypeOnlyForIndividuals := true

		for _, rule := range relType.Rules {
			if !(rule.PartyTypeRule.FirstPartyTypeID == iam.IndividualPartyType.ID && rule.PartyTypeRule.SecondPartyTypeID == iam.IndividualPartyType.ID) {
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

func (s *Server) PostIndividual(
	ctx context.Context,
	attrs *iam.AttributeList,
	id string,
	w http.ResponseWriter,
	req *http.Request) error {

	iamClient, err := s.IAMClient(req)
	if err != nil {
		return err
	}

	if err := req.ParseForm(); err != nil {
		return err
	}

	b := iam.NewIndividual(id)

	attributeMap := map[string]*iam.Attribute{}
	for _, attribute := range attrs.Items {
		attributeMap[attribute.ID] = attribute
	}

	type RelationshipEntry struct {
		*iam.Relationship
		MarkedForDeletion bool
	}

	//goland:noinspection GoPreferNilSlice
	var rels = []*RelationshipEntry{}

	// Validate Attributes
	attrErrs := validation.ErrorList{}
	f := req.Form
	for key, vals := range f {

		// Populate the Party.attributes
		// as well as the Attributes object (for validation)
		if attr := attrs.FindByName(key); attr != nil {
			b.Attributes[attr.ID] = vals
			attr.Attributes.Value = vals
			attrErrs = append(attrErrs, iam.ValidateAttribute(attr, validation.NewPath(""))...)
		}

		// Retrieve party relationships
		if strings.HasPrefix(key, "relationships[") {

			keyParts := strings.Split(key, ".")
			if len(keyParts) != 2 {
				err := fmt.Errorf("unexpected form value key: %s", key)
				return err
			}

			if !strings.HasSuffix(keyParts[0], "]") {
				err := fmt.Errorf("unexpected form value key: %s", key)
				return err
			}

			relationshipIdxStr := keyParts[0][14 : len(keyParts[0])-1]
			relationshipIdx, err := strconv.Atoi(relationshipIdxStr)
			if err != nil {
				return err
			}

			var rel *RelationshipEntry
			for {
				if (relationshipIdx + 1) > len(rels) {
					rels = append(rels, &RelationshipEntry{
						Relationship: &iam.Relationship{},
					})
				} else {
					rel = rels[relationshipIdx]
					break
				}
			}

			attrName := keyParts[1]

			switch attrName {
			case "markedForDeletion":
				rel.MarkedForDeletion = vals[0] == "true"
			case "id":
				rel.ID = vals[0]
			case "secondPartyId":
				rel.SecondPartyID = vals[0]
			case "relationshipTypeId":
				rel.RelationshipTypeID = vals[0]
			default:
				err := fmt.Errorf("unexpected relationship attribute: %s", attrName)
				return err
			}
		}
	}

	var individual *iam.Individual
	// Verify attribute validation and act accordingly
	if len(attrErrs) > 0 {
		status := attrErrs.Status(http.StatusUnprocessableEntity, "invalid case")
		return &status
	}

	// Update or create the individual
	if id == "" {
		var err error
		individual, err = iamClient.Individuals().Create(ctx, b)
		if err != nil {
			return err
		}
		err = s.createDefaultIndividualIntakeCases(req, individual)
		if err != nil {
			return err
		}

		if err := s.sessionManager.AddNotification(req, w, &sessionmanager.Notification{
			Message: fmt.Sprintf("Individual \"%s\" successfully created", b.String()),
			Theme:   "success",
		}); err != nil {
			return err
		}

	} else {
		var err error
		if individual, err = iamClient.Individuals().Update(ctx, b); err != nil {
			return err
		}

		if err := s.sessionManager.AddNotification(req, w, &sessionmanager.Notification{
			Message: fmt.Sprintf("Individual \"%s\" successfully updated", b.String()),
			Theme:   "success",
		}); err != nil {
			return err
		}

	}

	// Update, create or delete the relationships

	for _, rel := range rels {

		if len(rel.ID) == 0 {
			// Create the relationship
			relationship := rel.Relationship
			relationship.FirstPartyID = individual.ID
			if _, err := iamClient.Relationships().Create(ctx, rel.Relationship); err != nil {
				return err
			}
		} else if !rel.MarkedForDeletion {
			// Update the relationship
			relationship := rel.Relationship
			relationship.FirstPartyID = individual.ID
			if _, err := iamClient.Relationships().Update(ctx, rel.Relationship); err != nil {
				return err
			}
		} else {
			// Delete the relationship
			if err := iamClient.Relationships().Delete(ctx, rel.Relationship.ID); err != nil {
				return err
			}
		}

	}

	w.Header().Set("Location", "/individuals/"+individual.ID)
	w.WriteHeader(http.StatusSeeOther)
	return nil
}

func (s *Server) createDefaultIndividualIntakeCases(req *http.Request, individual *iam.Individual) error {
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
			IntakeCase: situationAnalysisCaseType.IntakeCaseType,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
