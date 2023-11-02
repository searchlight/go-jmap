package mailbox

import "git.sr.ht/~rockorager/go-jmap"

// Filter argument for a /query operation (see RFC8620, section 5.5)
type Filter interface {
	implementsFilter()
}

// FilterOperator can be used to create complex filtering (e.g.: return
// mailboxes which are subscribed and NOT named Inbox)
type FilterOperator struct {
	// jmap.OperatorOR, jmap.OperatorAND or jmap.OperatorNOT
	Operator jmap.Operator `json:"operator,omitempty"`

	// List of nested FilterOperator or FilterCondition.
	Conditions []Filter `json:"conditions,omitempty"`
}

func (fo *FilterOperator) implementsFilter() {}

// See RFC8621, section 4.4.1.
type FilterCondition struct {
	// The Mailbox parentId property must match the given value exactly.
	ParentID jmap.ID `json:"parentId,omitempty"`
	// The Mailbox name property contains the given string.
	Name string `json:"name,omitempty"`
	// The Mailbox role property must match the given value exactly.
	Role Role `json:"role,omitempty"`
	// If true, a Mailbox matches if it has any non-null value for its role
	// property.
	HasAnyRole bool `json:"hasAnyRole,omitempty"`
	// The isSubscribed property of the Mailbox must be identical to the
	// value given to match the condition.
	IsSubscribed bool `json:"isSubscribed,omitempty"`
}

func (fc *FilterCondition) implementsFilter() {}
