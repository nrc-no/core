package clients

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/api/mimetypes"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/ory/hydra-client-go/client/admin"
	"net/http"
)

type Handler struct {
	hydraAdmin admin.ClientService
	ws         *restful.WebService
}

func NewHandler(hydraAdmin admin.ClientService) (*Handler, error) {
	h := &Handler{}

	ws := new(restful.WebService).
		Path("/apis/admin.nrc.no/v1/clients").
		Doc("clients.admin.nrc.no API")

	ws.Route(ws.PUT("/{clientId}").To(restfulUpdate(hydraAdmin)).
		Doc(`updates oauth2 client`).
		Param(ws.PathParameter("clientId", "client id").Required(true)).
		Operation("updateClient").
		Consumes(mimetypes.ApplicationJson).
		Produces(mimetypes.ApplicationJson).
		Reads(&types.OAuth2Client{}).
		Writes(&types.OAuth2Client{}).
		Returns(http.StatusOK, "OK", &types.OAuth2Client{}),
	)

	ws.Route(ws.POST("").To(restfulCreate(hydraAdmin)).
		Doc(`creates oauth2 client`).
		Operation("createClient").
		Consumes(mimetypes.ApplicationJson).
		Produces(mimetypes.ApplicationJson).
		Reads(&types.OAuth2Client{}).
		Writes(&types.OAuth2Client{}).
		Returns(http.StatusOK, "OK", &types.OAuth2Client{}),
	)

	ws.Route(ws.GET("/{clientId}").To(restfulGet(hydraAdmin)).
		Doc(`gets oauth2 client`).
		Param(ws.PathParameter("clientId", "client id").Required(true)).
		Produces(mimetypes.ApplicationJson).
		Operation("getClient").
		Writes(&types.OAuth2Client{}).
		Returns(http.StatusOK, "OK", &types.OAuth2Client{}),
	)

	ws.Route(ws.DELETE("/{clientId}").To(restfulDelete(hydraAdmin)).
		Doc(`deletes oauth2 client`).
		Param(ws.PathParameter("clientId", "client id").Required(true)).
		Operation("deleteClient").
		Returns(http.StatusOK, "OK", nil),
	)

	ws.Route(ws.GET("").To(restfulList(hydraAdmin)).
		Doc(`gets oauth2 clients`).
		Operation("listClients").
		Produces(mimetypes.ApplicationJson).
		Writes(&types.OAuth2ClientList{}).
		Returns(http.StatusOK, "OK", &types.OAuth2ClientList{}),
	)

	h.ws = ws
	return h, nil
}

func (h *Handler) WebService() *restful.WebService {
	return h.ws
}
