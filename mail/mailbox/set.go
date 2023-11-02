package mailbox

import (
	"git.sr.ht/~rockorager/go-jmap"
	"git.sr.ht/~rockorager/go-jmap/mail"
)

// This is a standard /set method as described in RFC8620, Section 5.3,
type Set struct {
	Account jmap.ID `json:"accountId,omitempty"`

	// This is a state string as returned by the Mailbox/get method
	IfInState string `json:"ifInState,omitempty"`

	// A map of a creation id (a temporary id set by the client) to Mailbox
	// objects, or null if no objects are to be created.
	Create map[jmap.ID]*Mailbox `json:"create,omitempty"`

	// A map of an id to a Patch object to apply to the current Mailbox
	// object with that id, or null if no objects are to be updated.
	Update map[jmap.ID]jmap.Patch `json:"update,omitempty"`

	// A list of ids for Mailbox objects to permanently delete
	Destroy []jmap.ID `json:"destroy,omitempty"`

	// If false, any attempt to destroy a Mailbox that still has Emails in
	// it will be rejected with a mailboxHasEmail SetError. If true, any
	// Emails that were in the Mailbox will be removed from it, and if in
	// no other Mailboxes, they will be destroyed when the Mailbox is
	// destroyed.
	OnDestroyRemoveEmails bool `json:"onDestroyRemoveEmails,omitempty"`
}

func (m *Set) Name() string { return "Mailbox/set" }

func (m *Set) Requires() []jmap.URI { return []jmap.URI{mail.URI} }

type SetResponse struct {
	Account jmap.ID `json:"accountId,omitempty"`

	// The state string that would have been returned by Mailbox/get before
	// making the requested changes
	OldState string `json:"oldState,omitempty"`

	// The state string that will now be returned by Mailbox/get.
	NewState string `json:"newState,omitempty"`

	// Created mailboxes
	Created map[jmap.ID]*Mailbox `json:"created,omitempty"`

	// Updated mailboxes
	Updated map[jmap.ID]*Mailbox `json:"updated,omitempty"`

	// Deleted mailbox ids
	Destroyed []jmap.ID `json:"destroyed,omitempty"`

	// A map of ID to a SetError for each record that failed to be created
	NotCreated map[jmap.ID]*jmap.SetError `json:"notCreated,omitempty"`

	// A map of ID to a SetError for each record that failed to be updated
	NotUpdated map[jmap.ID]*jmap.SetError `json:"notUpdated,omitempty"`

	// A map of ID to a SetError for each record that failed to be destroyed
	NotDestroyed map[jmap.ID]*jmap.SetError `json:"notDestroyed,omitempty"`
}

func newSetResponse() jmap.MethodResponse { return &SetResponse{} }
