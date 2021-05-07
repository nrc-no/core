package testing

import (
	"context"
	"github.com/EventStore/EventStore-Client-Go/client"
	"github.com/nrc-no/core/apps/api/pkg/apis/core/v1"
	"github.com/nrc-no/core/apps/api/pkg/client/nrc"
	"github.com/nrc-no/core/apps/api/pkg/client/rest"
	"github.com/nrc-no/core/apps/api/pkg/endpoints/handlers/formdefinitions"
	"github.com/nrc-no/core/apps/api/pkg/runtime"
	"github.com/nrc-no/core/apps/api/pkg/runtime/serializer/json"
	"github.com/nrc-no/core/apps/api/pkg/server"
	mongostorage "github.com/nrc-no/core/apps/api/pkg/storage/mongo"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http/httptest"
	"testing"
)

type MainTestSuite struct {
	suite.Suite
	ctx                context.Context
	httpServer         *httptest.Server
	apiServer          *server.Server
	nrcClient          *nrc.NrcCoreClient
	mongoClient        *mongo.Client
	eventStoreDBClient *client.Client
	store              *mongostorage.Store
}

func TestMainSuite(t *testing.T) {
	suite.Run(t, new(MainTestSuite))
}

func (s *MainTestSuite) SetupSuite() {
	ctx := context.Background()
	s.ctx = ctx

	// Create API server
	apiServer := server.NewServer()
	s.apiServer = apiServer

	// Create HTTP server
	httpServer := httptest.NewServer(apiServer)
	s.httpServer = httpServer

	// Create client
	nrcClient, err := nrc.NewForConfig(&rest.Config{
		ContentConfig: rest.DefaultContentConfig,
		Host:          httpServer.URL,
	})
	if err != nil {
		s.T().Errorf("unable to create rest client: %v", err)
		return
	}
	s.nrcClient = nrcClient

	// Create eventdb client
	eventStoreDBClient, err := client.NewClient(&client.Configuration{
		Address:    "localhost:2113",
		DisableTLS: true,
	})
	if err != nil {
		s.T().Errorf("failed to create eventstoredb client: %v", err)
		return
	}
	if err := eventStoreDBClient.Connect(); err != nil {
		s.T().Errorf("failed to connect to evenstore: %v", err)
		return
	}
	s.eventStoreDBClient = eventStoreDBClient

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:30001"))
	if err != nil {
		s.T().Errorf("could not connect to mongo: %v", err)
		return
	}
	s.mongoClient = mongoClient

	scheme := runtime.NewScheme()
	if err := v1.AddToScheme(scheme); err != nil {
		s.T().Errorf("unable to register scheme: %v", err)
		return
	}

	serializer := json.NewSerializer(json.DefaultMetaFactory, scheme, scheme)

	// Create storage
	formDefinitionsStore := mongostorage.NewStore(
		eventStoreDBClient,
		mongoClient,
		"core",
		"core.nrc.no/formdefinitions",
		func() runtime.Object { return &v1.FormDefinition{} })
	s.store = formDefinitionsStore

	// Install FormDefinitions api
	formdefinitions.Install(
		apiServer.Container,
		formDefinitionsStore,
		v1.SchemeGroupVersion.WithKind("FormDefinition"),
		v1.SchemeGroupVersion.WithResource("formdefinitions"),
		scheme,
		scheme,
		serializer,
	)

}

func (s *MainTestSuite) TearDownSuite() {
	defer s.httpServer.Close()
	defer s.mongoClient.Disconnect(s.ctx)
	defer s.eventStoreDBClient.Close()
}
