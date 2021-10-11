package attachments_test

import (
	"github.com/nrc-no/core/internal/generic/server"
	"github.com/nrc-no/core/internal/testutils"
	. "github.com/nrc-no/core/pkg/attachments"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"
	"testing"
)

type Suite struct {
	*server.GenericServerTestSetup
	suite.Suite
	server *Server
	client *RESTAttachmentClient
}

func (s *Suite) SetupSuite() {
	s.GenericServerTestSetup = server.NewGenericServerTestSetup()
	s.server = NewServerOrDie(s.Ctx, s.GenericServerOptions)
	s.client = NewClient(testutils.SetXAuthenticatedUserSubject(s.Port))
	s.Serve(s.T(), s.server)
}

// This will run before each test in the suite but must be called manually before subtests
func (s *Suite) SetupTest() {
	err := s.server.ResetDB(s.Ctx, s.GenericServerOptions.MongoDatabase)
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *Suite) TearDownSuite() {
	s.SetupTest()
}

func TestSuite(t *testing.T) {
	suite.Run(t, &Suite{})
}

//
// Helpers
//

func (s *Suite) NewUUID() string {
	return uuid.NewV4().String()
}

func (s *Suite) uuidSlice(n int) []string {
	t := []string{}
	for i := 0; i < n; i++ {
		t = append(t, s.NewUUID())
	}
	return t
}

func (s *Suite) aMockAttachment(attachedToId string) *Attachment {
	var newUUIDForAttachment = s.NewUUID()
	return &Attachment{
		//ID:           newUUIDForAttachment,
		AttachedToID: attachedToId,
		Body:         "{\"data\":\"" + newUUIDForAttachment + "\"}",
	}
}

func (s *Suite) mockAttachments(n int, attachedToId string) []*Attachment {
	var attachments []*Attachment
	for i := 0; i < n; i++ {
		attachments = append(attachments, s.aMockAttachment(attachedToId))
	}
	return attachments
}