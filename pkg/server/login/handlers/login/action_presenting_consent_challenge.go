package login

import (
	"github.com/gorilla/sessions"
	"github.com/looplab/fsm"
	"github.com/nrc-no/core/pkg/logging"
	"github.com/nrc-no/core/pkg/server/login/authrequest"
	"github.com/nrc-no/core/pkg/server/login/templates"
	"github.com/ory/hydra-client-go/models"
	"go.uber.org/zap"
	"net/http"
)

func handlePresentingConsentChallenge(
	w http.ResponseWriter,
	req *http.Request,
	userSession *sessions.Session,
	getConsentRequest func(consentChallenge string) (*models.ConsentRequest, error),
) func(authRequest *authrequest.AuthRequest, evt *fsm.Event) error {

	return func(authRequest *authrequest.AuthRequest, evt *fsm.Event) error {

		ctx := req.Context()
		l := logging.NewLogger(ctx).With(zap.String("state", authrequest.StatePresentingConsent))

		l.Debug("getting consent request")
		consentRequest, err := getConsentRequest(authRequest.ConsentChallenge)
		if err != nil {
			l.Error("failed to get consent request", zap.Error(err))
			return err
		}

		// The fsm automatically saves the session at the end of the transitions.
		// Though, since we're storing this in a cookie, once we execute the template
		// below, the http.ResponseWriter will be closed. So saving the session
		// will have no effect. We simply save the session here to avoid this problem.
		if err := authRequest.Save(w, req, userSession); err != nil {
			l.Error("failed to save session", zap.Error(err))
			return err
		}

		l.Debug("prompting user for consent")
		err = templates.Template.ExecuteTemplate(w, "challenge", map[string]interface{}{
			"Scopes":     consentRequest.RequestedScope,
			"ClientName": consentRequest.Client.ClientName,
		})
		if err != nil {
			l.Error("failed to prompt user for consent", zap.Error(err))
			return err
		}

		return nil

	}
}
