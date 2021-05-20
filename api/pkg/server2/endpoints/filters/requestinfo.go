package filters

import (
	"fmt"
	"github.com/nrc-no/core/api/pkg/server2/endpoints/request"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	"net/http"
)

func WithRequestInfo(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		info, err := request.NewRequestInfo(req)
		if err != nil {
			responsewriters.InternalError(w, req, fmt.Errorf("failed to create RequestInfo: %v", err))
			return
		}
		req = req.WithContext(request.WithRequestInfo(ctx, info))
		handler.ServeHTTP(w, req)
	})
}
