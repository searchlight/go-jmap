package mailbox

import (
	"git.sr.ht/~rockorager/go-jmap"
	"git.sr.ht/~rockorager/go-jmap/mail"
)

// This is a standard /changes method as described in RFC8620, Section 5.2.
type Changes struct {
	Account jmap.ID `json:"accountId,omitempty"`

	// The current state of the client.
	SinceState string `json:"sinceState,omitempty"`

	// The maximum number of ids to return in the response.
	MaxChanges uint64 `json:"maxChanges,omitempty"`
}

func (m *Changes) Name() string { return "Mailbox/changes" }

func (m *Changes) Requires() []jmap.URI { return []jmap.URI{mail.URI} }

// This is a standard /changes method as described in RFC8620, Section 5.2
// but with one extra argument to the response: updatedProperties
type ChangesResponse struct {
	Account jmap.ID `json:"accountId,omitempty"`

	// This is the SinceState argument echoed back
	OldState string `json:"oldState,omitempty"`

	// This is the state the client will be in after applying the set of
	// changes to the old state.
	NewState string `json:"newState,omitempty"`

	// If true, the client may call /changes again with the NewState
	// returned to get further updates. If false, NewState is the current
	// server state.
	HasMoreChanges bool `json:"hasMoreChanges,omitempty"`

	// New mailbox IDs.
	Created []jmap.ID `json:"created,omitempty"`

	// Updated mailbox IDs.
	Updated []jmap.ID `json:"updated,omitempty"`

	// Deleted mailbox IDs.
	Destroyed []jmap.ID `json:"destroyed,omitempty"`

	// If only the “totalEmails”, “unreadEmails”, “totalThreads”, and/or
	// “unreadThreads” Mailbox properties have changed since the old state,
	// this will be the list of properties that may have changed.
	UpdatedProperties []string `json:"updatedProperties,omitempty"`
}

func newChangesResponse() jmap.MethodResponse { return &ChangesResponse{} }
