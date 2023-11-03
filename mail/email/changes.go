package email

import (
	"git.sr.ht/~rockorager/go-jmap"
	"git.sr.ht/~rockorager/go-jmap/mail"
)

// This is a standard "/changes" method as described in [RFC8620], Section 5.2.
type Changes struct {
	// The id of the account to use.
	Account jmap.ID `json:"accountId,omitempty"`

	// The current state of the client
	SinceState string `json:"sinceState,omitempty"`

	// The maximum number of ids to return in the response
	MaxChanges uint64 `json:"maxChanges,omitempty"`
}

func (m *Changes) Name() string { return "Email/changes" }

func (m *Changes) Requires() []jmap.URI { return []jmap.URI{mail.URI} }

// This is a standard "/changes" method as described in [RFC8620], Section 5.2.
type ChangesResponse struct {
	// The id of the account used for the call.
	Account jmap.ID `json:"accountId,omitempty"`

	// This is the sinceState argument echoed back
	OldState string `json:"oldState,omitempty"`

	// The state the client will be in after applying the Changes
	NewState string `json:"newState,omitempty"`

	// If true, not all changes were returned in this response
	HasMoreChanges bool `json:"hasMoreChanges,omitempty"`

	Created   []jmap.ID `json:"created,omitempty"`
	Updated   []jmap.ID `json:"updated,omitempty"`
	Destroyed []jmap.ID `json:"destroyed,omitempty"`
}

func newChangesResponse() jmap.MethodResponse { return &ChangesResponse{} }
