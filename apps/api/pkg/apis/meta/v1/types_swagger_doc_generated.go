// Hello!

package v1

// This file contains a collection of methods that can be used from go-restful to
// generate Swagger API documentation for its models. Please read this PR for more
// information on the implementation: https://github.com/emicklei/go-restful/pull/215
//
// TODOs are ignored from the parser (e.g. TODO(andronat):... || TODO:...) if and only if
// they are on one line! For multiple line or blocks that you want to ignore use ---.
// Any context after a --- is ignored.
//
// Those methods can be generated by using hack/update-generated-swagger-docs.sh

// AUTO-GENERATED FUNCTIONS START HERE. DO NOT EDIT.
var map_APIGroup = map[string]string{
	"":                           "APIGroup contains the name, the supported versions, and the preferred version of a group.",
	"name":                       "name is the name of the group.",
	"versions":                   "versions are the versions supported in this group.",
	"preferredVersion":           "preferredVersion is the version preferred by the API server, which probably is the storage version.",
	"serverAddressByClientCIDRs": "a map of client CIDR to server address that is serving this group. This is to help clients reach servers in the most network-efficient way possible. Clients can use the appropriate server address as per the CIDR that they match. In case of multiple matches, clients should use the longest matching CIDR. The server returns only those CIDRs that it thinks that the client can match. For example: the master will return an internal IP CIDR only, if the client reaches the server using an internal IP. Server looks at X-Forwarded-For header or X-Real-Ip header or request.RemoteAddr (in that order) to get the client IP.",
}

func (APIGroup) SwaggerDoc() map[string]string {
	return map_APIGroup
}

var map_APIGroupList = map[string]string{
	"":       "APIGroupList is a list of APIGroup, to allow clients to discover the API at /apis.",
	"groups": "groups is a list of APIGroup.",
}

func (APIGroupList) SwaggerDoc() map[string]string {
	return map_APIGroupList
}

var map_APIResource = map[string]string{
	"":                   "APIResource specifies the name of a resource and whether it is namespaced.",
	"name":               "name is the plural name of the resource.",
	"singularName":       "singularName is the singular name of the resource.  This allows clients to handle plural and singular opaquely. The singularName is more correct for reporting status on a single item and both singular and plural are allowed from the kubectl CLI interface.",
	"namespaced":         "namespaced indicates if a resource is namespaced or not.",
	"group":              "group is the preferred group of the resource.  Empty implies the group of the containing resource list. For subresources, this may have a different value, for example: Scale\".",
	"version":            "version is the preferred version of the resource.  Empty implies the version of the containing resource list For subresources, this may have a different value, for example: v1 (while inside a v1beta1 version of the core resource's group)\".",
	"kind":               "kind is the kind for the resource (e.g. 'Foo' is the kind for a resource 'foo')",
	"verbs":              "verbs is a list of supported kube verbs (this includes get, list, watch, create, update, patch, delete, deletecollection, and proxy)",
	"shortNames":         "shortNames is a list of suggested short names of the resource.",
	"categories":         "categories is a list of the grouped resources this resource belongs to (e.g. 'all')",
	"storageVersionHash": "The hash value of the storage version, the version this resource is converted to when written to the data store. Value must be treated as opaque by clients. Only equality comparison on the value is valid. This is an alpha feature and may change or be removed in the future. The field is populated by the apiserver only if the StorageVersionHash feature gate is enabled. This field will remain optional even if it graduates.",
}

func (APIResource) SwaggerDoc() map[string]string {
	return map_APIResource
}

var map_APIResourceList = map[string]string{
	"":             "APIResourceList is a list of APIResource, it is used to expose the name of the resources supported in a specific group and version, and if the resource is namespaced.",
	"groupVersion": "groupVersion is the group and version this APIResourceList is for.",
	"resources":    "resources contains the name of the resources and if they are namespaced.",
}

func (APIResourceList) SwaggerDoc() map[string]string {
	return map_APIResourceList
}

var map_APIVersions = map[string]string{
	"":                           "APIVersions lists the versions that are available, to allow clients to discover the API at /api, which is the root path of the legacy v1 API.",
	"versions":                   "versions are the api versions that are available.",
	"serverAddressByClientCIDRs": "a map of client CIDR to server address that is serving this group. This is to help clients reach servers in the most network-efficient way possible. Clients can use the appropriate server address as per the CIDR that they match. In case of multiple matches, clients should use the longest matching CIDR. The server returns only those CIDRs that it thinks that the client can match. For example: the master will return an internal IP CIDR only, if the client reaches the server using an internal IP. Server looks at X-Forwarded-For header or X-Real-Ip header or request.RemoteAddr (in that order) to get the client IP.",
}

func (APIVersions) SwaggerDoc() map[string]string {
	return map_APIVersions
}

var map_GroupVersionForDiscovery = map[string]string{
	"":             "GroupVersion contains the \"group/version\" and \"version\" string of a version. It is made a struct to keep extensibility.",
	"groupVersion": "groupVersion specifies the API group and version in the form \"group/version\"",
	"version":      "version specifies the version in the form of \"version\". This is to save the clients the trouble of splitting the GroupVersion.",
}

func (GroupVersionForDiscovery) SwaggerDoc() map[string]string {
	return map_GroupVersionForDiscovery
}

var map_LabelSelector = map[string]string{
	"":                 "A label selector is a label query over a set of resources. The result of matchLabels and matchExpressions are ANDed. An empty label selector matches all objects. A null label selector matches no objects.",
	"matchLabels":      "matchLabels is a map of {key,value} pairs. A single {key,value} in the matchLabels map is equivalent to an element of matchExpressions, whose key field is \"key\", the operator is \"In\", and the values array contains only \"value\". The requirements are ANDed.",
	"matchExpressions": "matchExpressions is a list of label selector requirements. The requirements are ANDed.",
}

func (LabelSelector) SwaggerDoc() map[string]string {
	return map_LabelSelector
}

var map_LabelSelectorRequirement = map[string]string{
	"":         "A label selector requirement is a selector that contains values, a key, and an operator that relates the key and values.",
	"key":      "key is the label key that the selector applies to.",
	"operator": "operator represents a key's relationship to a set of values. Valid operators are In, NotIn, Exists and DoesNotExist.",
	"values":   "values is an array of string values. If the operator is In or NotIn, the values array must be non-empty. If the operator is Exists or DoesNotExist, the values array must be empty. This array is replaced during a strategic merge patch.",
}

func (LabelSelectorRequirement) SwaggerDoc() map[string]string {
	return map_LabelSelectorRequirement
}

var map_List = map[string]string{
	"":         "List holds a list of objects, which may not be known by the server.",
	"metadata": "Standard list metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds",
	"items":    "List of objects",
}

func (List) SwaggerDoc() map[string]string {
	return map_List
}

var map_ListMeta = map[string]string{
	"": "ListMeta describes metadata that synthetic resources must have",
}

func (ListMeta) SwaggerDoc() map[string]string {
	return map_ListMeta
}

var map_ListOptions = map[string]string{
	"":                     "ListOptions is the query options to a standard REST list call.",
	"labelSelector":        "A selector to restrict the list of returned objects by their labels. Defaults to everything.",
	"fieldSelector":        "A selector to restrict the list of returned objects by their fields. Defaults to everything.",
	"watch":                "Watch for changes to the described resources and return them as a stream of add, update, and remove notifications. Specify resourceVersion.",
	"allowWatchBookmarks":  "allowWatchBookmarks requests watch events with type \"BOOKMARK\". Servers that do not implement bookmarks may ignore this flag and bookmarks are sent at the server's discretion. Clients should not assume bookmarks are returned at any specific interval, nor may they assume the server will send any BOOKMARK event during a session. If this is not a watch, this field is ignored.",
	"resourceVersion":      "resourceVersion sets a constraint on what resource versions a request may be served from. See https://kubernetes.io/docs/reference/using-api/api-concepts/#resource-versions for details.\n\nDefaults to unset",
	"resourceVersionMatch": "resourceVersionMatch determines how resourceVersion is applied to list calls. It is highly recommended that resourceVersionMatch be set for list calls where resourceVersion is set See https://kubernetes.io/docs/reference/using-api/api-concepts/#resource-versions for details.\n\nDefaults to unset",
	"timeoutSeconds":       "Timeout for the list/watch call. This limits the duration of the call, regardless of any activity or inactivity.",
	"limit":                "limit is a maximum number of responses to return for a list call. If more items exist, the server will set the `continue` field on the list metadata to a value that can be used with the same initial query to retrieve the next set of results. Setting a limit may return fewer than the requested amount of items (up to zero items) in the event all requested objects are filtered out and clients should only use the presence of the continue field to determine whether more results are available. Servers may choose not to support the limit argument and will return all of the available results. If limit is specified and the continue field is empty, clients may assume that no more results are available. This field is not supported if watch is true.\n\nThe server guarantees that the objects returned when using continue will be identical to issuing a single list call without a limit - that is, no objects created, modified, or deleted after the first request is issued will be included in any subsequent continued requests. This is sometimes referred to as a consistent snapshot, and ensures that a client that is using limit to receive smaller chunks of a very large result can ensure they see all possible objects. If objects are updated during a chunked list the version of the object that was present at the time the first list result was calculated is returned.",
	"continue":             "The continue option should be set when retrieving more results from the server. Since this value is server defined, clients may only use the continue value from a previous query result with identical query parameters (except for the value of continue) and the server may reject a continue value it does not recognize. If the specified continue value is no longer valid whether due to expiration (generally five to fifteen minutes) or a configuration change on the server, the server will respond with a 410 ResourceExpired error together with a continue token. If the client needs a consistent list, it must restart their list without the continue field. Otherwise, the client may send another list request with the token received with the 410 error, the server will respond with a list starting from the next key, but from the latest snapshot, which is inconsistent from the previous list results - objects that are created, modified, or deleted after the first list request will be included in the response, as long as their keys are after the \"next key\".\n\nThis field is not supported when watch is true. Clients may start a watch from the last resourceVersion value returned by the server and not miss any modifications.",
}

func (ListOptions) SwaggerDoc() map[string]string {
	return map_ListOptions
}

var map_ServerAddressByClientCIDR = map[string]string{
	"":              "ServerAddressByClientCIDR helps the client to determine the server address that they should use, depending on the clientCIDR that they match.",
	"clientCIDR":    "The CIDR with which clients can match their IP to figure out the server address that they should use.",
	"serverAddress": "Address of this server, suitable for a client that matches the above CIDR. This can be a hostname, hostname:port, IP or IP:port.",
}

func (ServerAddressByClientCIDR) SwaggerDoc() map[string]string {
	return map_ServerAddressByClientCIDR
}

var map_Status = map[string]string{
	"status":  "Status is the status of the operation. One of \"Failure\" or \"Success\"",
	"message": "Message is a human-readable description of this operation",
	"reason":  "Reason is a machine-readable description of why this operation is in the \"Failure\" status.",
	"details": "Details represents extended data associated with the reason.",
	"code":    "Suggested HTTP status code.",
}

func (Status) SwaggerDoc() map[string]string {
	return map_Status
}

var map_TypeMeta = map[string]string{
	"":           "TypeMeta represents an individual object in an API response or request. It represents the API schema version and kind/type of object",
	"apiVersion": "APIVersion defines the versioned schema of this representation of an object.",
	"kind":       "Kind is a string value representing the REST resource this object represents",
}

func (TypeMeta) SwaggerDoc() map[string]string {
	return map_TypeMeta
}

// AUTO-GENERATED FUNCTIONS END HERE
