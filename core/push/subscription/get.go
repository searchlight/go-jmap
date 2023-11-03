package subscription

import (
	"git.sr.ht/~rockorager/go-jmap"
	"git.sr.ht/~rockorager/go-jmap/core"
)

// This is a standard “/get” method as described in [@!RFC8620], Section 5.1.
type Get struct {
	IDs        []jmap.ID `json:"ids,omitempty"`
	Properties []string  `json:"properties,omitempty"`
}

func (m *Get) Name() string { return "PushSubscription/get" }

func (m *Get) Requires() []jmap.URI { return []jmap.URI{core.URI} }

// This is a standard “/get” method as described in [@!RFC8620], Section 5.1.
type GetResponse struct {
	List     []*PushSubscription `json:"list,omitempty"`
	NotFound []jmap.ID           `json:"notFound,omitempty"`
}

func newGetResponse() jmap.MethodResponse { return &GetResponse{} }
