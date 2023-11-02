package mailbox

import (
	"git.sr.ht/~rockorager/go-jmap"
	"git.sr.ht/~rockorager/go-jmap/mail"
)

type QueryChanges struct {
	Account jmap.ID `json:"accountId,omitempty"`

	// The filter argument that was used with Mailbox/query.
	Filter Filter `json:"filter,omitempty"`

	// The sort argument that was used with Mailbox/query.
	Sort []*SortComparator `json:"sort,omitempty"`

	// The current state of the query in the client. This is the string
	// that was returned as the queryState argument in the Mailbox/query
	// response with the same sort/filter. The server will return the
	// changes made to the query since this state.
	SinceQueryState string `json:"sinceQueryState,omitempty"`

	// The maximum number of changes to return in the response.
	MaxChanges uint64 `json:"maxChanges,omitempty"`

	// The last (highest-index) id the client currently has cached from the
	// query results.
	UpToID jmap.ID `json:"upToId,omitempty"`

	// Does the client wish to know the total number of results now in the
	// query?
	CalculateTotal bool `json:"calculateTotal,omitempty"`
}

func (m *QueryChanges) Name() string { return "Mailbox/queryChanges" }

func (m *QueryChanges) Requires() []jmap.URI { return []jmap.URI{mail.URI} }

type QueryChangesResponse struct {
	Account jmap.ID `json:"accountId,omitempty"`

	// This is the SinceQueryState argument echoed back
	OldQueryState string `json:"oldQueryState,omitempty"`

	// This is the state the query will be in after applying the set of
	// changes to the old state.
	NewQueryState string `json:"newQueryState,omitempty"`

	// Deleted Mailbox IDs
	Removed []jmap.ID `json:"removed,omitempty"`

	// Added Mailbox IDs
	Added []*jmap.AddedItem `json:"added,omitempty"`
}

func newQueryChangesResponse() jmap.MethodResponse { return &QueryChangesResponse{} }
