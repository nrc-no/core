package login

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/nrc-no/core/internal/generic/server"
	"github.com/nrc-no/core/internal/rest"
	"github.com/nrc-no/core/internal/utils"
	iam2 "github.com/nrc-no/core/pkg/iam"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"html/template"
	"net/http"
	"path"
)

type ServerOptions struct {
	*server.GenericServerOptions
	BCryptCost        int
	AdminHTTPClient   *http.Client
	IAMHost           string
	IAMScheme         string
	TemplateDirectory string
}

type Server struct {
	HydraAdmin              admin.ClientService
	BCryptCost              int
	router                  *mux.Router
	template                *template.Template
	iam                     iam2.Interface
	HydraHTTPClient         *http.Client
	credentialsCollectionFn func() (*mongo.Collection, error)
}

func NewServer(ctx context.Context, o *ServerOptions) (*Server, error) {

	iamCli := iam2.NewClientSet(&rest.RESTConfig{
		Scheme:     o.IAMScheme,
		Host:       o.IAMHost,
		HTTPClient: o.AdminHTTPClient,
	})

	srv := &Server{
		HydraAdmin:      o.HydraAdminClient.Admin,
		BCryptCost:      o.BCryptCost,
		iam:             iamCli,
		HydraHTTPClient: o.HydraHTTPClient,
		credentialsCollectionFn: func() (*mongo.Collection, error) {
			mongoClient, err := o.MongoClientFn(ctx)
			if err != nil {
				logrus.WithError(err).Errorf("unable to get mongo client")
				return nil, err
			}
			collection := mongoClient.Database(o.MongoDatabase).Collection("credentials")
			return collection, nil
		},
	}

	router := mux.NewRouter()
	router.Path("/auth/logout").Methods("GET").HandlerFunc(srv.GetLogoutForm)
	router.Path("/auth/login").Methods("GET").HandlerFunc(srv.GetLoginForm)
	router.Path("/auth/login").Methods("POST").HandlerFunc(srv.PostLoginForm)
	router.Path("/auth/consent").Methods("GET").HandlerFunc(srv.GetConsent)
	router.Path("/auth/consent").Methods("POST").HandlerFunc(srv.PostConsent)
	router.Path("/apis/login/v1/credentials").
		Methods("POST").
		Handler(srv.WithAuth()(http.HandlerFunc(srv.PostCredentials)))

	srv.router = router

	tpl, err := template.ParseGlob(path.Join(o.TemplateDirectory, "*.gohtml"))
	if err != nil {
		logrus.WithError(err).Errorf("failed to parse templates")
		return nil, err
	}

	srv.template = tpl

	return srv, nil

}

func (s *Server) json(w http.ResponseWriter, status int, data interface{}) {
	utils.JSONResponse(w, status, data)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.router.ServeHTTP(w, req)
}

func (s *Server) Error(w http.ResponseWriter, err error) {
	logrus.WithError(err).Error()
	utils.ErrorResponse(w, err)
}

func (s *Server) Bind(req *http.Request, into interface{}) error {
	return utils.BindJSON(req, into)
}