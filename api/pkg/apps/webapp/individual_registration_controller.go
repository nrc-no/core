package webapp

import (
	"fmt"
	"github.com/nrc-no/core/pkg/apps/seeder"
	"github.com/nrc-no/core/pkg/registrationctrl"
	"net/http"
)

func (s *Server) GetRegistrationController(w http.ResponseWriter, req *http.Request) (*registrationctrl.RegistrationController, error) {
	var individualId string
	if !s.GetPathParam("id", w, req, &individualId) {
		return nil, fmt.Errorf("cannot find id in path")
	}

	iamClient, err := s.IAMClient(req)
	if err != nil {
		return nil, err
	}

	i, err := iamClient.Individuals().Get(req.Context(), individualId)
	if err != nil {
		return nil, err
	}

	irh := NewIndividualRegistrationHandler(s, i, req)

	return registrationctrl.NewRegistrationController(irh, seeder.UgandaRegistrationFlow), nil
}
