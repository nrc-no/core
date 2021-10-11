package webapp

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nrc-no/core/internal/form"
	"github.com/nrc-no/core/internal/sessionmanager"
	"github.com/nrc-no/core/internal/validation"
	"github.com/nrc-no/core/pkg/cms"
	iam2 "github.com/nrc-no/core/pkg/iam"
	"github.com/satori/go.uuid"
	"github.com/xeonx/timeago"
	"golang.org/x/sync/errgroup"
	"net/http"
	"net/url"
)

func (s *Server) Cases(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	cmsClient, err := s.CMSClient(req)
	if err != nil {
		s.Error(w, err)
		return
	}

	iamClient, err := s.IAMClient(req)
	if err != nil {
		s.Error(w, err)
		return
	}

	options := &CasesListOptions{}
	if err := options.UnmarshalQueryParams(req.URL.Query()); err != nil {
		s.Error(w, err)
		return
	}

	var kases *cms.CaseList
	var caseTypes *cms.CaseTypeList
	var partyList *iam2.PartyList
	var teams *iam2.TeamList

	caseListOptions := &cms.CaseListOptions{}
	caseListOptions.Done = options.Closed
	caseListOptions.CaseTypeIDs = options.CaseTypeIDs
	caseListOptions.TeamIDs = options.TeamIDs

	wg, wgCtx := errgroup.WithContext(ctx)
	wg.Go(func() error {
		var err error
		kases, err = cmsClient.Cases().List(wgCtx, *caseListOptions)
		return err
	})
	wg.Go(func() error {
		var err error
		caseTypes, err = cmsClient.CaseTypes().List(wgCtx, cms.CaseTypeListOptions{})
		return err
	})
	wg.Go(func() error {
		var err error
		partyList, err = iamClient.Parties().List(wgCtx, iam2.PartyListOptions{})
		return err
	})
	wg.Go(func() error {
		var err error
		teams, err = iamClient.Teams().List(wgCtx, iam2.TeamListOptions{})
		return err
	})

	if err := wg.Wait(); err != nil {
		s.Error(w, err)
		return
	}

	if err := s.renderFactory.New(req, w).ExecuteTemplate(w, "cases", map[string]interface{}{
		"Cases":         kases,
		"CaseTypes":     caseTypes,
		"Parties":       partyList,
		"Teams":         teams,
		"FilterOptions": options,
	}); err != nil {
		s.Error(w, err)
		return
	}

}

func (s *Server) Case(w http.ResponseWriter, req *http.Request) {
	var (
		err              error
		recipientParty   *iam2.Party
		team             *iam2.Team
		kase             *cms.Case
		parent           *cms.Case
		kaseTypes        *cms.CaseTypeList
		referrals        *cms.CaseList
		referralCaseType *cms.CaseType
		creator          *iam2.Party
		comments         *cms.CommentList
		commentAuthors   *iam2.PartyList
		commentAuthorMap = make(map[string]*iam2.Party)
	)

	ctx := req.Context()

	cmsClient, err := s.CMSClient(req)
	if err != nil {
		s.Error(w, err)
		return
	}

	iamClient, err := s.IAMClient(req)
	if err != nil {
		s.Error(w, err)
		return
	}

	caseID, referralCaseTypeID, err := s.getCaseIds(req)
	if err != nil {
		s.Error(w, err)
		return
	}

	kase, err = s.getCase(ctx, cmsClient, caseID, referralCaseTypeID)
	if err != nil {
		s.Error(w, err)
		return
	}

	g, waitCtx := errgroup.WithContext(ctx)

	// Comments
	g.Go(func() error {
		if comments, err = cmsClient.Comments().List(waitCtx, cms.CommentListOptions{
			CaseID: caseID,
		}); err == nil {
			var commentAuthorIDMap = make(map[string]bool)
			var commentAuthorIDs []string
			// Get all comment author IDs
			for _, comment := range comments.Items {
				commentAuthorIDMap[comment.AuthorID] = true
			}
			for authorID := range commentAuthorIDMap {
				commentAuthorIDs = append(commentAuthorIDs, authorID)
			}
			commentAuthors, err = iamClient.Parties().Search(waitCtx, iam2.PartySearchOptions{
				PartyTypeIDs: []string{iam2.IndividualPartyType.ID},
				PartyIDs:     commentAuthorIDs,
			})
			for _, author := range commentAuthors.Items {
				commentAuthorMap[author.ID] = author
			}
		}
		return err
	})

	// Creator
	g.Go(func() error {
		if len(kase.CreatorID) > 0 {
			creator, err = iamClient.Parties().Get(waitCtx, kase.CreatorID)
		}
		return err
	})

	// Team
	g.Go(func() error {
		team, err = iamClient.Teams().Get(waitCtx, kase.TeamID)
		return err
	})

	// Parent Case
	g.Go(func() error {
		if len(kase.ParentID) > 0 {
			parent, err = cmsClient.Cases().Get(waitCtx, kase.ParentID)
		}
		return err
	})

	// Recipient & CaseTypes
	g.Go(func() error {
		if recipientParty, err = iamClient.Parties().Get(waitCtx, kase.PartyID); err == nil {
			kaseTypes, err = cmsClient.CaseTypes().List(waitCtx, cms.CaseTypeListOptions{
				PartyTypeIDs: recipientParty.PartyTypeIDs,
			})
		}
		return err
	})

	// Referrals
	g.Go(func() error {
		referrals, err = cmsClient.Cases().List(waitCtx, cms.CaseListOptions{ParentID: kase.ID})
		return err
	})
	g.Go(func() error {
		if len(referralCaseTypeID) > 0 {
			referralCaseType, err = cmsClient.CaseTypes().Get(waitCtx, referralCaseTypeID)
		}
		return err
	})

	if err = g.Wait(); err != nil {
		s.Error(w, err)
		return
	}

	valuedKase := form.NewValidatedForm(kase.Form, kase.FormData, nil)
	if err := s.renderFactory.New(req, w).ExecuteTemplate(w, "case", map[string]interface{}{
		"Case":             kase,
		"CaseForm":         valuedKase,
		"Parent":           parent,
		"CaseTypes":        kaseTypes,
		"Recipient":        recipientParty,
		"ReferralCaseType": referralCaseType,
		"Referrals":        referrals,
		"Team":             team,
		"CreatedBy":        creator,
		"Comments":         displayComments(comments, commentAuthorMap),
	}); err != nil {
		s.Error(w, err)
	}
}

func (s *Server) NewCase(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()
	cmsClient, err := s.CMSClient(req)
	if err != nil {
		s.Error(w, err)
		return
	}

	iamClient, err := s.IAMClient(req)
	if err != nil {
		s.Error(w, err)
		return
	}

	var caseTypes *cms.CaseTypeList
	var p *iam2.PartyList

	g, waitCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		caseTypes, err = cmsClient.CaseTypes().List(waitCtx, cms.CaseTypeListOptions{})
		return err
	})

	if err := g.Wait(); err != nil {
		s.Error(w, err)
		return
	}

	qry := req.URL.Query()
	caseTypeID := qry.Get("caseTypeId")
	partyTypeID := ""
	teamID := ""
	for _, caseType := range caseTypes.Items {
		if caseType.ID == caseTypeID {
			partyTypeID = caseType.PartyTypeID
			teamID = caseType.TeamID
			break
		}
	}

	listOptions := iam2.PartyListOptions{
		PartyTypeID: partyTypeID,
	}

	p, err = iamClient.Parties().List(ctx, listOptions)
	if err != nil {
		s.Error(w, err)
	}

	var caseType *cms.CaseType
	if caseTypeID != "" {
		caseType, err = cmsClient.CaseTypes().Get(ctx, caseTypeID)
		if err != nil {
			s.Error(w, err)
			return
		}
	}

	var team *iam2.Team
	if len(teamID) > 0 {
		team, err = iamClient.Teams().Get(ctx, teamID)
		if err != nil {
			s.Error(w, err)
			return
		}
	}

	var phorm form.Form
	if caseType != nil {
		phorm = form.NewValidatedForm(caseType.Form, nil, nil)
	}

	if err := s.renderFactory.New(req, w).ExecuteTemplate(w, "casenew", map[string]interface{}{
		"PartyID":   qry.Get("partyId"),
		"CaseType":  caseType,
		"Form":      phorm,
		"Team":      team,
		"CaseTypes": caseTypes,
		"Parties":   p,
	}); err != nil {
		s.Error(w, err)
		return
	}
}

func (s *Server) PostCase(w http.ResponseWriter, req *http.Request) {
	var err error
	ctx := req.Context()

	cmsClient, err := s.CMSClient(req)
	if err != nil {
		return
	}

	caseID, referralCaseTypeID, err := s.getCaseIds(req)
	if err != nil {
		s.Error(w, err)
		return
	}

	kase, err := s.getCase(ctx, cmsClient, caseID, referralCaseTypeID)
	if err != nil {
		s.Error(w, err)
		return
	}

	if err := req.ParseForm(); err != nil {
		return
	}
	values := req.Form

	if kase.CaseTypeID == "" {
		// we need to get the caseTypeId from the form data
		kase.CaseTypeID = values.Get("caseTypeId")
		if kase.CaseTypeID == "" {
			err = fmt.Errorf("unable to detect case type id for new case")
			return
		}
	}

	caseType, err := cmsClient.CaseTypes().Get(ctx, kase.CaseTypeID)
	if err != nil {
		return
	}

	kase.TeamID = values.Get("teamId")
	kase.PartyID = values.Get("partyId")
	kase.Done = values.Get("done") == "on"
	kase.Form = caseType.Form
	kase.FormData = req.Form

	var storedCase *cms.Case
	var isNewCase = kase.ID == ""
	if isNewCase {
		kase.IntakeCase = caseType.IntakeCaseType
		storedCase, err = cmsClient.Cases().Create(ctx, kase)
	} else {
		storedCase, err = cmsClient.Cases().Update(ctx, kase)
	}
	if err != nil {
		if status, ok := err.(*validation.Status); ok {
			formValidation := makeFormValidation(status.Errors, kase.Form)
			s.json(w, status.Code, formValidation)
		} else {
			s.Error(w, err)
		}
		return
	}

	s.redirectAfterPost(w, req, storedCase, isNewCase)

	return
}

func (s *Server) getCaseIds(req *http.Request) (caseID string, referralCaseTypeID string, err error) {
	qry := req.URL.Query()
	referralCaseTypeID = qry.Get("referralCaseTypeId")
	caseID, ok := mux.Vars(req)["id"]
	if !ok || len(caseID) == 0 {
		if req.Method != "POST" {
			err := fmt.Errorf("no id in path")
			return "", "", err
		}
	}
	return caseID, referralCaseTypeID, nil
}

func (s *Server) getCase(ctx context.Context, cmsClient cms.Interface, caseID string, referralCaseTypeID string) (*cms.Case, error) {
	isNewCase := len(caseID) == 0
	kase := &cms.Case{}
	if isNewCase {
		if len(referralCaseTypeID) > 0 {
			kase.CaseTypeID = referralCaseTypeID
		}
	} else {
		var err error
		kase, err = cmsClient.Cases().Get(ctx, caseID)
		if err != nil {
			return nil, err
		}
	}
	return kase, nil
}

func (s *Server) redirectAfterPost(w http.ResponseWriter, req *http.Request, posted *cms.Case, isNewCase bool) {
	var action string
	if isNewCase {
		action = "created"
	} else {
		action = "updated"
	}
	if err := s.sessionManager.AddNotification(req, w, &sessionmanager.Notification{
		Message: fmt.Sprintf("Case successfully %s", action),
		Theme:   "success",
	}); err != nil {
		s.Error(w, err)
		return
	}

	if len(posted.ParentID) > 0 {
		w.Header().Set("Location", "/cases/"+posted.ParentID)
	} else {
		w.Header().Set("Location", "/cases/"+posted.ID)
	}
	w.WriteHeader(http.StatusSeeOther)
	return
}

type displayComment struct {
	*cms.Comment
	Author  *iam2.Party
	TimeAgo string
}

func displayComments(comments *cms.CommentList, authorMap map[string]*iam2.Party) []*displayComment {
	var displayComments []*displayComment
	for _, item := range comments.Items {
		c := &displayComment{
			Comment: item,
			TimeAgo: timeago.English.Format(item.CreatedAt),
			Author:  authorMap[item.AuthorID],
		}
		displayComments = append(displayComments, c)
	}
	return displayComments
}

type CasesListOptions struct {
	Closed      *bool
	TeamIDs     []string
	CaseTypeIDs []string
}

func (c *CasesListOptions) ClosedOnly() bool {
	return c.Closed != nil && *c.Closed == true
}

func (c *CasesListOptions) OpenOnly() bool {
	return c.Closed != nil && *c.Closed == false
}

func (c *CasesListOptions) UnmarshalQueryParams(values url.Values) error {

	if len(values["status"]) == 1 {
		closed := values["status"][0] == "closed"
		c.Closed = &closed
	}

	for _, teamId := range values["teamId"] {
		if _, err := uuid.FromString(teamId); err == nil {
			c.TeamIDs = append(c.TeamIDs, teamId)
		}
	}

	for _, caseTypeId := range values["caseTypeId"] {
		if _, err := uuid.FromString(caseTypeId); err == nil {
			c.CaseTypeIDs = append(c.CaseTypeIDs, caseTypeId)
		}
	}

	return nil

}