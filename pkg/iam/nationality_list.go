package iam

import "net/http"

func (s *Server) listNationalities(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var listOptions NationalityListOptions
	if err := listOptions.UnmarshalQueryParameters(req.URL.Query()); err != nil {
		s.error(w, err)
		return
	}

	ret, err := s.nationalityStore.list(ctx, listOptions)
	if err != nil {
		s.error(w, err)
		return
	}

	s.json(w, http.StatusOK, ret)
}