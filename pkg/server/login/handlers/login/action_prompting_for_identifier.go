package login

import (
	"github.com/gorilla/sessions"
	"github.com/looplab/fsm"
	"github.com/nrc-no/core/pkg/logging"
	"github.com/nrc-no/core/pkg/server/login/authrequest"
	"github.com/nrc-no/core/pkg/server/login/templates"
	"go.uber.org/zap"
	"net/http"
)

func handlePromptingForIdentifier(w http.ResponseWriter, userSession *sessions.Session, req *http.Request) func(authRequest *authrequest.AuthRequest, evt *fsm.Event) error {
	return func(authRequest *authrequest.AuthRequest, evt *fsm.Event) error {
		ctx := req.Context()
		l := logging.NewLogger(ctx).With(zap.String("state", authrequest.StatePromptingForIdentifier))

		// The fsm automatically saves the session at the end of the transitions.
		// Though, since we're storing this in a cookie, once we execute the template
		// below, the http.ResponseWriter will be closed. So saving the session
		// will have no effect. We simply save the session here to avoid this problem.
		if err := authRequest.Save(w, req, userSession); err != nil {
			l.Error("failed to save session", zap.Error(err))
			return err
		}

		l.Debug("prompting user for identifier")
		err := templates.Template.ExecuteTemplate(w, "login_subject", map[string]interface{}{
			"Error": authRequest.IdentifierError,
		})
		if err != nil {
			l.Error("failed to prompt user for identifier", zap.Error(err))
			return err
		}
		return nil
	}
}
