package authrequest

import (
	"encoding/gob"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/looplab/fsm"
	loginstore "github.com/nrc-no/core/pkg/server/login/store"
	"net/http"
	"strings"
	"time"
)

const (
	StateStart                        = ""
	StateLoginRequested               = "login_requested"
	StateRefreshingIdentity           = "refreshing_identity"
	StatePromptingForIdentifier       = "prompting_for_identifier"
	StateValidatingIdentifier         = "validating_identifier"
	StateFindingAuthMethod            = "finding_auth_method"
	StatePromptingForPassword         = "prompting_for_password"
	StateValidatingPassword           = "validating_password"
	StatePasswordLoginSucceeded       = "password_login_succeeded"
	StatePromptingForIdentityProvider = "prompting_for_identity_provider"
	StateAwaitingIdpCallback          = "awaiting_idp_callback"
	StatePerformingAuthCodeExchange   = "performing_auth_code_exchange"
	StateAuthCodeExchangeSucceeded    = "auth_code_exchange_succeeded"
	StateAwaitingConsentChallenge     = "awaiting_consent_challenge"
	StateReceivedConsentChallenge     = "received_consent_challenge"
	StatePresentingConsent            = "presenting_consent"
	StateConsentRequestApproved       = "consent_request_approved"
	StateConsentRequestDeclined       = "consent_request_declined"
	StateAccepted                     = "accepted"
	StateDeclined                     = "declined"
	StateFailed                       = "failed"

	EventRequestLogin             = "request_login"
	EventSkipLoginRequest         = "skip_login_request"
	EventPerformLogin             = "perform_login"
	EventProvideIdentifier        = "provide_identifier"
	EventProvideInvalidIdentifier = "provide_invalid_identifier"
	EventProvideValidIdentifier   = "provide_valid_identifier"
	EventUsePasswordLogin         = "use_password_login"
	EventProvideValidPassword     = "provide_valid_password"
	EventProvideInvalidPassword   = "provide_invalid_password"
	EventProvidePassword          = "provide_password"
	EventUseOidcLogin             = "use_oidc_login"
	EventUseIdentityProvider      = "use_identity_provider"
	EventCallOidcCallback         = "call_oidc_callback"
	EventFailCodeExchange         = "fail_auth_code_exchange"
	EventSucceedCodeExchange      = "succeed_auth_code_exchange"
	EventAcceptLoginRequest       = "accept_login_request"
	EventReceiveConsentChallenge  = "receive_consent_challenge"
	EventSkipConsentChallenge     = "skip_consent_challenge"
	EventPresentConsentChallenge  = "present_consent_challenge"
	EventApproveConsentChallenge  = "approve_consent_challenge"
	EventDeclineConsentChallenge  = "decline_consent_challenge"
	EventSetRefreshedIdentity     = "set_refreshed_identity"
	EventFail                     = "fail"
	EventAccept                   = "accept"
	EventDecline                  = "decline"
)

var allStates = []string{
	StateStart,
	StateLoginRequested,
	StateRefreshingIdentity,
	StatePromptingForIdentifier,
	StateValidatingIdentifier,
	StateFindingAuthMethod,
	StatePromptingForPassword,
	StateValidatingPassword,
	StatePasswordLoginSucceeded,
	StatePromptingForIdentityProvider,
	StateAwaitingIdpCallback,
	StatePerformingAuthCodeExchange,
	StateAuthCodeExchangeSucceeded,
	StateAwaitingConsentChallenge,
	StateReceivedConsentChallenge,
	StatePresentingConsent,
	StateConsentRequestApproved,
	StateConsentRequestDeclined,
	StateAccepted,
	StateDeclined,
	StateFailed,
}

var allEvents = []string{
	EventRequestLogin,
	EventSkipLoginRequest,
	EventPerformLogin,
	EventProvideIdentifier,
	EventProvideInvalidIdentifier,
	EventProvideValidIdentifier,
	EventUsePasswordLogin,
	EventProvideValidPassword,
	EventProvideInvalidPassword,
	EventProvidePassword,
	EventUseOidcLogin,
	EventUseIdentityProvider,
	EventCallOidcCallback,
	EventFailCodeExchange,
	EventSucceedCodeExchange,
	EventAcceptLoginRequest,
	EventReceiveConsentChallenge,
	EventSkipConsentChallenge,
	EventPresentConsentChallenge,
	EventApproveConsentChallenge,
	EventDeclineConsentChallenge,
	EventSetRefreshedIdentity,
	EventFail,
	EventAccept,
	EventDecline,
}

func NewAuthRequest(callbacks fsm.Callbacks) *AuthRequest {
	return newAuthRequest(StateStart, callbacks)
}

func newAuthRequest(state string, callbacks fsm.Callbacks) *AuthRequest {
	authRequest := &AuthRequest{}
	authRequest.fsm = fsm.NewFSM(state,
		[]fsm.EventDesc{

			// START
			// start -- request login --> login requested
			{Src: []string{StateStart}, Name: EventRequestLogin, Dst: StateLoginRequested},
			// accepted -- request login --> awaiting consent challenge
			{Src: []string{StateAccepted}, Name: EventRequestLogin, Dst: StateLoginRequested},
			// login requested -- skip login --> refreshing identity
			{Src: []string{StateLoginRequested}, Name: EventSkipLoginRequest, Dst: StateAwaitingConsentChallenge},
			// login requested -- skip login --> awaiting consent challenge
			{Src: []string{StateRefreshingIdentity}, Name: EventSetRefreshedIdentity, Dst: StateAwaitingConsentChallenge},
			// login requested -- perform login --> identifier needed
			{Src: []string{StateLoginRequested}, Name: EventPerformLogin, Dst: StatePromptingForIdentifier},

			// IDENTIFIER
			// identifier needed -- provide identifier --> validating identifier
			{Src: []string{StatePromptingForIdentifier}, Name: EventProvideIdentifier, Dst: StateValidatingIdentifier},
			// validating identifier -- invalid identifier --> identifier needed
			{Src: []string{StateValidatingIdentifier}, Name: EventProvideInvalidIdentifier, Dst: StatePromptingForIdentifier},
			// validating identifier -- valid identifier --> finding auth method
			{Src: []string{StateValidatingIdentifier}, Name: EventProvideValidIdentifier, Dst: StateFindingAuthMethod},

			// PASSWORD
			// finding auth method -- use password login --> password needed
			{Src: []string{StateFindingAuthMethod}, Name: EventUsePasswordLogin, Dst: StatePromptingForPassword},
			// password needed -- provide password --> validating password
			{Src: []string{StatePromptingForPassword}, Name: EventProvidePassword, Dst: StateValidatingPassword},
			// validating password -- invalid password --> password needed
			{Src: []string{StateValidatingPassword}, Name: EventProvideInvalidPassword, Dst: StatePromptingForPassword},
			// validating password -- valid password --> password login succeeded
			{Src: []string{StateValidatingPassword}, Name: EventProvideValidPassword, Dst: StatePasswordLoginSucceeded},
			// password login succeeded -- accept password login --> login request accepted
			{Src: []string{StatePasswordLoginSucceeded}, Name: EventAcceptLoginRequest, Dst: StateAwaitingConsentChallenge},

			// OIDC
			// finding auth method -- use oidc login --> identity provided needed
			{Src: []string{StateFindingAuthMethod}, Name: EventUseOidcLogin, Dst: StatePromptingForIdentityProvider},
			// identity provided needed -- use identity provider -- awaiting idp callback
			{Src: []string{StatePromptingForIdentityProvider}, Name: EventUseIdentityProvider, Dst: StateAwaitingIdpCallback},
			// awaiting callback -- call idp callback -- performing auth code exchange
			{Src: []string{StateAwaitingIdpCallback}, Name: EventCallOidcCallback, Dst: StatePerformingAuthCodeExchange},
			// performing auth code exchange -- fail code exchange -- auth code exchange failed
			{Src: []string{StatePerformingAuthCodeExchange}, Name: EventFailCodeExchange, Dst: StateFindingAuthMethod},
			// performing auth code exchange -- succeed code exchange -- auth code exchange succeeded
			{Src: []string{StatePerformingAuthCodeExchange}, Name: EventSucceedCodeExchange, Dst: StateAuthCodeExchangeSucceeded},
			// auth code exchange succeeded -- accept oidc login --> login request accepted
			{Src: []string{StateAuthCodeExchangeSucceeded}, Name: EventAcceptLoginRequest, Dst: StateAwaitingConsentChallenge},

			// CONSENT
			// awaiting consent challenge -- receive consent challenge --> received consent challenge
			{Src: []string{StateAwaitingConsentChallenge}, Name: EventReceiveConsentChallenge, Dst: StateReceivedConsentChallenge},
			// received consent challenge -- skip consent challenge --> accepted
			{Src: []string{StateReceivedConsentChallenge}, Name: EventSkipConsentChallenge, Dst: StateConsentRequestApproved},
			// received consent challenge -- present consent challenge --> presenting consent challenge
			{Src: []string{StateReceivedConsentChallenge}, Name: EventPresentConsentChallenge, Dst: StatePresentingConsent},
			// presenting consent challenge -- decline consent --> declined
			{Src: []string{StatePresentingConsent}, Name: EventDeclineConsentChallenge, Dst: StateConsentRequestDeclined},
			// presenting consent challenge -- approve consent --> accepted
			{Src: []string{StatePresentingConsent}, Name: EventApproveConsentChallenge, Dst: StateConsentRequestApproved},
			// consent challenge approved -- accept --> accepted
			{Src: []string{StateConsentRequestApproved}, Name: EventAccept, Dst: StateAccepted},
			// consent challenge decline -- decline --> declined
			{Src: []string{StateConsentRequestDeclined}, Name: EventDecline, Dst: StateDeclined},

			// GENERIC FAIL
			{Src: allStates, Name: EventFail, Dst: StateFailed},
		},
		callbacks,
	)

	return authRequest

}

type AuthRequest struct {
	fsm                *fsm.FSM
	Identifier         string
	IdentifierError    error
	EmailDomain        string
	LoginChallenge     string
	StateVariable      string
	IdentityProviderId string
	IDToken            string
	AccessToken        string
	RefreshToken       string
	TokenExpiry        time.Time
	TokenType          string
	Claims             *Claims
	Identity           *loginstore.Identity
	ConsentChallenge   string
	PostConsentURL     string
	Error              error
	Subject            string
}

func (a *AuthRequest) ToDotGraph() string {
	result := strings.Builder{}
	result.WriteString("digraph A {\n")

	for _, state := range allStates {
		result.WriteString(fmt.Sprintf("  %s\n", state))
	}
	for _, state := range allStates {
		for _, trans := range a.fsm.AvailableTransitions() {
			a.fsm.SetState(state)
			_ = a.fsm.Event(trans)
			currentState := a.currentState()
			result.WriteString(fmt.Sprintf("%s -> %s [label=\"%s\"]", state, currentState, trans))
		}
	}

	result.WriteString("}\n")
	return result.String()

}

type StoredIdentity struct {
	ID          string
	State       string
	Credentials []StoredCredential
}

type StoredCredential struct {
	ID                  string
	CreatedAt           time.Time
	UpdatedAt           time.Time
	IdentityID          string
	Kind                string
	HashedPassword      *string
	Issuer              *string
	Identifiers         []StoredCredentialIdentifier
	InitialAccessToken  *string
	InitialRefreshToken *string
	InitialIdToken      *string
}

type StoredCredentialIdentifier struct {
	ID           string
	CredentialID string
	Identifier   string
}

type StoredAuthRequest struct {
	CurrentState       string
	Identifier         string
	IdentifierError    error
	EmailDomain        string
	LoginChallenge     string
	StateVariable      string
	IdentityProviderId string
	IDToken            string
	AccessToken        string
	RefreshToken       string
	Claims             *Claims
	Identity           *StoredIdentity
	ConsentChallenge   string
	PostConsentURL     string
	TokenExpiry        time.Time
	TokenType          string
	Subject            string
}

func init() {
	gob.Register(&StoredAuthRequest{})
	gob.Register(&StoredIdentity{})
	gob.Register(StoredCredential{})
	gob.Register([]StoredCredential{})
	gob.Register(StoredCredentialIdentifier{})
	gob.Register([]StoredCredentialIdentifier{})
	gob.Register(&Claims{})
}

type Claims struct {
	Subject       string `json:"sub"`
	DisplayName   string `json:"display_name"`
	FullName      string `json:"name"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
}

func (a *AuthRequest) Save(w http.ResponseWriter, req *http.Request, session *sessions.Session) error {

	stored := StoredAuthRequest{
		CurrentState:       a.fsm.Current(),
		Identifier:         a.Identifier,
		IdentifierError:    a.IdentifierError,
		EmailDomain:        a.EmailDomain,
		LoginChallenge:     a.LoginChallenge,
		StateVariable:      a.StateVariable,
		IdentityProviderId: a.IdentityProviderId,
		IDToken:            a.IDToken,
		AccessToken:        a.AccessToken,
		TokenExpiry:        a.TokenExpiry,
		TokenType:          a.TokenType,
		RefreshToken:       a.RefreshToken,
		Claims:             a.Claims,
		ConsentChallenge:   a.ConsentChallenge,
		PostConsentURL:     a.PostConsentURL,
		Subject:            a.Subject,
	}

	if a.Identity != nil {
		identity := &StoredIdentity{}
		identity.ID = a.Identity.ID
		identity.State = string(a.Identity.State)
		for _, credential := range a.Identity.Credentials {
			cred := StoredCredential{}
			cred.IdentityID = credential.IdentityID
			cred.ID = credential.ID
			cred.CreatedAt = credential.CreatedAt
			cred.UpdatedAt = credential.UpdatedAt
			cred.IdentityID = credential.IdentityID
			cred.Kind = string(credential.Kind)
			cred.HashedPassword = credential.HashedPassword
			cred.Issuer = credential.Issuer
			cred.InitialAccessToken = credential.InitialAccessToken
			cred.InitialRefreshToken = credential.InitialRefreshToken
			cred.InitialIdToken = credential.InitialIdToken
			for _, identifier := range credential.Identifiers {
				var iden = StoredCredentialIdentifier{}
				iden.ID = identifier.ID
				iden.Identifier = identifier.Identifier
				iden.CredentialID = identifier.CredentialID
				cred.Identifiers = append(cred.Identifiers, iden)
			}
			identity.Credentials = append(identity.Credentials, cred)
		}
		stored.Identity = identity
	}

	session.Values["auth-request"] = stored
	return session.Save(req, w)
}

func CreateOrRestore(session *sessions.Session, callbacks fsm.Callbacks) *AuthRequest {
	authRequestIntf, ok := session.Values["auth-request"]
	if !ok {
		return NewAuthRequest(callbacks)
	}
	storedAuthRequest, ok := authRequestIntf.(*StoredAuthRequest)
	if !ok {
		return NewAuthRequest(callbacks)
	}
	authRequest := newAuthRequest(storedAuthRequest.CurrentState, callbacks)
	authRequest.Identifier = storedAuthRequest.Identifier
	authRequest.IdentifierError = storedAuthRequest.IdentifierError
	authRequest.EmailDomain = storedAuthRequest.EmailDomain
	authRequest.LoginChallenge = storedAuthRequest.LoginChallenge
	authRequest.StateVariable = storedAuthRequest.StateVariable
	authRequest.IdentityProviderId = storedAuthRequest.IdentityProviderId
	authRequest.IDToken = storedAuthRequest.IDToken
	authRequest.AccessToken = storedAuthRequest.AccessToken
	authRequest.RefreshToken = storedAuthRequest.RefreshToken
	authRequest.Claims = storedAuthRequest.Claims
	authRequest.ConsentChallenge = storedAuthRequest.ConsentChallenge
	authRequest.PostConsentURL = storedAuthRequest.PostConsentURL
	authRequest.TokenExpiry = storedAuthRequest.TokenExpiry
	authRequest.TokenType = storedAuthRequest.TokenType
	authRequest.Subject = storedAuthRequest.Subject

	if storedAuthRequest.Identity != nil {
		iden := &loginstore.Identity{}
		iden.ID = storedAuthRequest.Identity.ID
		iden.State = loginstore.IdentityState(storedAuthRequest.Identity.State)
		for _, credential := range storedAuthRequest.Identity.Credentials {
			cred := &loginstore.Credential{}
			cred.Identity = iden
			cred.IdentityID = credential.IdentityID
			cred.ID = credential.ID
			cred.CreatedAt = credential.CreatedAt
			cred.UpdatedAt = credential.UpdatedAt
			cred.IdentityID = credential.IdentityID
			cred.Kind = loginstore.CredentialKind(credential.Kind)
			cred.HashedPassword = credential.HashedPassword
			cred.Issuer = credential.Issuer
			cred.InitialAccessToken = credential.InitialAccessToken
			cred.InitialRefreshToken = credential.InitialRefreshToken
			cred.InitialIdToken = credential.InitialIdToken
			for _, identifier := range credential.Identifiers {
				var iden = &loginstore.CredentialIdentifier{}
				iden.Credential = cred
				iden.ID = identifier.ID
				iden.Identifier = identifier.Identifier
				iden.CredentialID = identifier.CredentialID
				cred.Identifiers = append(cred.Identifiers, iden)
			}
			iden.Credentials = append(iden.Credentials, cred)
		}
		authRequest.Identity = iden
	}

	return authRequest
}

func (a *AuthRequest) currentState() string {
	return a.fsm.Current()
}

type RequestLoginProps struct {
	LoginChallenge string
	ClientID       string
	Scope          string
	Subject        *string
}

func (a *AuthRequest) Event(evt string) error {
	return a.fsm.Event(evt)
}

func (a *AuthRequest) State() string {
	return a.fsm.Current()
}
