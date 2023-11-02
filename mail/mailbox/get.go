package mailbox

import (
	"git.sr.ht/~rockorager/go-jmap"
	"git.sr.ht/~rockorager/go-jmap/mail"
)

// Objects of type Mailbox are fetched via a call to Mailbox/get The ids
// argument may be null to fetch all at once.
type Get struct {
	Account jmap.ID `json:"accountId,omitempty"`

	// The ids of the Mailbox objects to return.
	IDs []jmap.ID `json:"ids,omitempty"`

	// If supplied, only the properties listed in the array are returned
	// for each Mailbox object.
	Properties []string `json:"properties,omitempty"`

	// Use IDs from a previous call
	ReferenceIDs *jmap.ResultReference `json:"#ids,omitempty"`

	// Use Properties from a previous call
	ReferenceProperties *jmap.ResultReference `json:"#properties,omitempty"`
}

func (m *Get) Name() string {
	return "Mailbox/get"
}

func (m *Get) Requires() []jmap.URI { return []jmap.URI{mail.URI} }

type GetResponse struct {
	Account jmap.ID `json:"accountId,omitempty"`

	// A string representing the state on the server for all the data of
	// this type in the account.
	State string `json:"state,omitempty"`

	// List of the Mailbox objects requested.
	List []*Mailbox `json:"list,omitempty"`

	// List of Mailbox IDs that do not exist.
	NotFound []string `json:"notFound,omitempty"`
}

func newGetResponse() jmap.MethodResponse { return &GetResponse{} }
