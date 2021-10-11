package login

import "net/http"

func (s *Server) Render(w http.ResponseWriter, req *http.Request, templateName string, data map[string]interface{}) {
	if err := s.template.ExecuteTemplate(w, templateName, data); err != nil {
		s.Error(w, err)
		return
	}
}