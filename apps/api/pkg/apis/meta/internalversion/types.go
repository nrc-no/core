package internalversion

import (
	metav1 "github.com/nrc-no/core/apps/api/pkg/apis/meta/v1"
	"github.com/nrc-no/core/apps/api/pkg/fields"
	"github.com/nrc-no/core/apps/api/pkg/labels"
	"github.com/nrc-no/core/apps/api/pkg/runtime"
)

// ListOptions is the query options to a standard REST list call.
type ListOptions struct {
	metav1.TypeMeta

	// A selector based on labels
	LabelSelector labels.Selector
	// A selector based on fields
	FieldSelector fields.Selector
	// If true, watch for changes to this list
	Watch bool
	// allowWatchBookmarks requests watch events with type "BOOKMARK".
	// Servers that do not implement bookmarks may ignore this flag and
	// bookmarks are sent at the server's discretion. Clients should not
	// assume bookmarks are returned at any specific interval, nor may they
	// assume the server will send any BOOKMARK event during a session.
	// If this is not a watch, this field is ignored.
	// If the feature gate WatchBookmarks is not enabled in apiserver,
	// this field is ignored.
	AllowWatchBookmarks bool
	// resourceVersion sets a constraint on what resource versions a request may be served from.
	// See https://kubernetes.io/docs/reference/using-api/api-concepts/#resource-versions for
	// details.
	ResourceVersion string
	// resourceVersionMatch determines how resourceVersion is applied to list calls.
	// It is highly recommended that resourceVersionMatch be set for list calls where
	// resourceVersion is set.
	// See https://kubernetes.io/docs/reference/using-api/api-concepts/#resource-versions for
	// details.
	ResourceVersionMatch metav1.ResourceVersionMatch

	// Timeout for the list/watch call.
	TimeoutSeconds *int64
	// Limit specifies the maximum number of results to return from the server. The server may
	// not support this field on all resource types, but if it does and more results remain it
	// will set the continue field on the returned list object.
	Limit int64
	// Continue is a token returned by the server that lets a client retrieve chunks of results
	// from the server by specifying limit. The server may reject requests for continuation tokens
	// it does not recognize and will return a 410 error if the token can no longer be used because
	// it has expired.
	Continue string
}

// List holds a list of objects, which may not be known by the server.
type List struct {
	metav1.TypeMeta
	// +optional
	metav1.ListMeta

	Items []runtime.Object
}
