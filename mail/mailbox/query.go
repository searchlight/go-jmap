package mailbox

import (
	"git.sr.ht/~rockorager/go-jmap"
	"git.sr.ht/~rockorager/go-jmap/mail"
)

// This is a standard /query method as described in RFC8620, Section 5.5,
// but with the following additional request argument: sortAsTree, filterAsTree
type Query struct {
	Account jmap.ID `json:"accountId,omitempty"`

	// Determines the set of emails returned in the results.
	Filter Filter `json:"filter,omitempty"`

	// Lists the names of properties to compare between two Email records,
	Sort []*SortComparator `json:"sort,omitempty"`

	// The zero-based index of the first id in the full list of results to
	// return.
	Position int64 `json:"position,omitempty"`

	// An Email id to use along with AnchorOffset.
	Anchor jmap.ID `json:"anchor,omitempty"`

	// The index of the first result to return relative to the index of the
	// anchor, if an anchor is given.
	AnchorOffset int64 `json:"anchorOffset,omitempty"`

	// The maximum number of results to return.
	Limit uint64 `json:"limit,omitempty"`

	// Does the client wish to know the total number of results in the
	// query?
	CalculateTotal bool `json:"calculateTotal,omitempty"`

	// Sort mailboxes according to their tree structure first, and then
	// according to the Sort argument.
	SortAsTree bool `json:"sortAsTree,omitempty"`

	// If true, a Mailbox is only included in the query if all its
	// ancestors are also included in the query according to the filter.
	FilterAsTree bool `json:"filterAsTree,omitempty"`
}

func (m *Query) Name() string { return "Mailbox/query" }

func (m *Query) Requires() []jmap.URI { return []jmap.URI{mail.URI} }

type QueryResponse struct {
	Account jmap.ID `json:"accountId,omitempty"`

	// A string encoding the current state of the query on the server.
	QueryState string `json:"queryState,omitempty"`

	// This is true if the server supports calling Mailbox/queryChanges
	CanCalculateChanges bool `json:"canCalculateChanges,omitempty"`

	// The zero-based index of the first result in the ids array within the
	// complete list of query results.
	Position uint64 `json:"position,omitempty"`

	// The list of ids for each Mailbox in the query results
	IDs []jmap.ID `json:"ids,omitempty"`

	// The total number of Mailboxes in the results (given the filter).
	Total int64 `json:"total,omitempty"`

	// The limit enforced by the server on the maximum number of results to
	// return.
	Limit uint64 `json:"limit,omitempty"`
}

func newQueryResponse() jmap.MethodResponse { return &QueryResponse{} }
