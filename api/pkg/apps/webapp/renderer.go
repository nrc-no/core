package webapp

import (
	"errors"
	"fmt"
	"github.com/nrc-no/core/pkg/auth"
	"github.com/nrc-no/core/pkg/sessionmanager"
	"html/template"
	"io"
	"net/http"
	"path"
)

// RenderInterface defines the methods available in the template
type RenderInterface interface {
	IsLoggedIn() bool
	Profile() *Claims
	Notifications() []*sessionmanager.Notification
	Dict(values ...interface{}) (map[string]interface{}, error)
}

// RendererFactory is a factory to create Renderer
type RendererFactory struct {
	template       *template.Template
	sessionManager sessionmanager.Store
}

// RendererFactory must implement RenderInterface so that the methods
// are available at compile time
var _ RenderInterface = &RendererFactory{}

// IsLoggedIn is a mock method that allows compilation of go templates
// Actual implementation is in Renderer.IsLoggedIn
func (r *RendererFactory) IsLoggedIn() bool {
	return false
}

func (r *RendererFactory) Profile() (*Claims, error) {
	return nil, nil
}

func (r *RendererFactory) Notifications() ([]*sessionmanager.Notification, error) {
	return []*sessionmanager.Notification{}, nil
}

func (r *RendererFactory) Dict(values ...interface{}) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}

// NewRendererFactory creates a new instance of the RendererFactory
func NewRendererFactory(templateDirectory string, sessionManager sessionmanager.Store) (*RendererFactory, error) {
	f := &RendererFactory{
		sessionManager: sessionManager,
	}
	t := template.New("")
	t = WithRenderInterface(t, f)
	var err error
	t, err = t.ParseGlob(path.Join(templateDirectory, "*.gohtml"))
	if err != nil {
		return nil, err
	}
	f.template = t
	return f, nil
}

// New creates a new Renderer
func (r *RendererFactory) New(req *http.Request) *Renderer {
	renderer := &Renderer{
		req:            req,
		sessionManager: r.sessionManager,
	}
	renderer.template = WithRenderInterface(r.template, renderer)
	return renderer
}

// Renderer is the actual struct that will render templates
type Renderer struct {
	template       *template.Template
	req            *http.Request
	sessionManager sessionmanager.Store
}

// Renderer must implement RenderInterface so that the methods are available
// in the templates
var _ RenderInterface = &Renderer{}

func (r *Renderer) ExecuteTemplate(w io.Writer, name string, args interface{}) error {
	return r.template.ExecuteTemplate(w, name, args)
}

// IsLoggedIn returns whether the request is made by an authenticated user
func (r *Renderer) IsLoggedIn() bool {
	return auth.IsAuthenticatedRequest(r.req)
}

func (r *Renderer) Profile() (*Claims, error) {

	session, err := r.sessionManager.Get(r.req)
	if err != nil {
		return nil, err
	}

	profileIntf, ok := session.Values["profile"]
	if !ok {
		return nil, fmt.Errorf("profile not found")
	}

	profile, ok := profileIntf.(*Claims)
	if !ok {
		return nil, fmt.Errorf("")
	}

	return profile, nil
}

func (r *Renderer) Notifications() ([]*sessionmanager.Notification, error) {
	return r.sessionManager.ConsumeNotifications(r.req)
}

func (r *Renderer) Dict(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, errors.New("invalid dict call")
	}
	dict := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, errors.New("dict keys must be strings")
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}

// WithRenderInterface adds the RenderInterface methods to the template
func WithRenderInterface(t *template.Template, intf RenderInterface) *template.Template {
	return t.Funcs(map[string]interface{}{
		"IsLoggedIn":    intf.IsLoggedIn,
		"Profile":       intf.Profile,
		"Notifications": intf.Notifications,
		"Dict":          intf.Dict,
	})
}
