package iam

import (
	"github.com/nrc-no/core/internal/validation"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func (s *Server) postRelationshipType(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var payload RelationshipType
	if err := s.bind(req, &payload); err != nil {
		s.error(w, err)
		return
	}

	p := &payload
	p.ID = uuid.NewV4().String()

	errList := ValidateRelationshipType(p, validation.NewPath(""))
	if len(errList) > 0 {
		status := errList.Status(http.StatusUnprocessableEntity, "invalid relationship type")
		s.error(w, &status)
		return
	}

	if err := s.relationshipTypeStore.Create(ctx, p); err != nil {
		s.error(w, err)
		return
	}

	s.json(w, http.StatusOK, p)
}