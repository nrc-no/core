package webapp

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nrc-no/core-kafka/pkg/memberships"
	"github.com/nrc-no/core-kafka/pkg/parties/parties"
	"github.com/nrc-no/core-kafka/pkg/teams"
	"golang.org/x/sync/errgroup"
	"net/http"
	"sync"
)

func (h *Handler) Teams(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	t, err := h.teamClient.List(ctx, teams.ListOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.template.ExecuteTemplate(w, "teams", map[string]interface{}{
		"Teams": t,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (h *Handler) Team(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	id, ok := mux.Vars(req)["id"]
	if !ok || len(id) == 0 {
		err := fmt.Errorf("no id found in path")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t, err := h.teamClient.Get(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	m, err := h.membershipClient.List(ctx, memberships.ListOptions{
		TeamID: id,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var members []*parties.Party
	lock := sync.Mutex{}

	g, ctx := errgroup.WithContext(ctx)

	for _, item := range m.Items {
		i := item
		g.Go(func() error {
			individual, err := h.partyClient.Get(ctx, i.IndividualID)
			if err != nil {
				return err
			}
			lock.Lock()
			defer lock.Unlock()
			members = append(members, individual)
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.template.ExecuteTemplate(w, "team", map[string]interface{}{
		"Team":    t,
		"Members": members,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
