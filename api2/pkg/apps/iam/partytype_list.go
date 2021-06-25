package iam

import (
	"net/http"
)

func (s *Server) ListPartyTypes(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	ret, err := s.PartyTypeStore.List(ctx)
	if err != nil {
		s.Error(w, err)
		return
	}

	s.JSON(w, http.StatusOK, ret)
}
