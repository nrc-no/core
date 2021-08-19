package webapp

import (
	"github.com/sirupsen/logrus"
	"net/http"
)

func (s *Server) Login(w http.ResponseWriter, req *http.Request) {
	state, err := s.createHydraStateVariable()
	if err != nil {
		logrus.WithError(err).Errorf("failed to make a new state variable for hydra login flow")
		s.Error(w, err)
		return
	}

	session, err := s.sessionManager.Get(req)
	if err != nil {
		s.Error(w, err)
		return
	}

	session.Values["state"] = state
	if err := session.Save(req, w); err != nil {
		s.Error(w, err)
		return
	}

	redirectUrl := s.publicOauth2Config.AuthCodeURL(state)
	http.Redirect(w, req, redirectUrl, http.StatusTemporaryRedirect)

}
