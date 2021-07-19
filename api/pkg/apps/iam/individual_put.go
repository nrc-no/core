package iam

import (
	"github.com/nrc-no/core/pkg/validation"
	"net/http"
)

func (s *Server) putIndividual(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var id string

	if !s.getPathParam("id", w, req, &id) {
		return
	}

	var individual Individual
	if err := s.bind(req, &individual); err != nil {
		s.error(w, err)
		return
	}

	_, err := s.individualStore.get(ctx, id)
	if err != nil {
		s.error(w, err)
		return
	}

	errList := ValidateIndividual(&individual, validation.NewPath(""))
	if len(errList) > 0 {
		status := errList.Status(http.StatusUnprocessableEntity, "invalid individual")
		s.json(w, status.Code, status)
		return
	}

	if err := s.individualStore.upsert(ctx, &individual); err != nil {
		s.error(w, err)
		return
	}

	s.json(w, http.StatusOK, &individual)

}
