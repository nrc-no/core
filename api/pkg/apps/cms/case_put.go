package cms

import (
	"github.com/nrc-no/core/pkg/validation"
	"net/http"
)

func (s *Server) PutCase(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var id string
	if !s.GetPathParam("id", w, req, &id) {
		return
	}

	var payload Case
	if err := s.Bind(req, &payload); err != nil {
		s.Error(w, err)
		return
	}

	kase, err := s.caseStore.Get(ctx, id)
	if err != nil {
		s.Error(w, err)
		return
	}

	kase.Template = payload.Template

	if !kase.BypassValidation {
		errList := ValidateCase(kase, &validation.Path{})
		if len(errList) > 0 {
			status := errList.Status(http.StatusUnprocessableEntity, "invalid case")
			s.Error(w, &status)
			return
		}
	} else {
		kase.BypassValidation = false
	}

	// if no validation errors, assume the case is Done
	kase.Done = true

	if err := s.caseStore.Update(ctx, kase); err != nil {
		s.Error(w, err)
		return
	}

	s.JSON(w, http.StatusOK, kase)

}
