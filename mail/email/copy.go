package email

import (
	"git.sr.ht/~rockorager/go-jmap"
	"git.sr.ht/~rockorager/go-jmap/mail"
)

// This is a standard "/copy" method as described in [RFC8620], Section 5.4,
type Copy struct {
	// The id of the account to copy records from.
	FromAccount   jmap.ID `json:"fromAccountId,omitempty"`
	IfFromInState string  `json:"ifFromInState,omitempty"`

	// The id of the account to copy records to. This MUST be different to
	// the fromAccountId.
	Account                  jmap.ID            `json:"accountId,omitempty"`
	IfInState                string             `json:"ifInState,omitempty"`
	Create                   map[jmap.ID]*Email `json:"create,omitempty"`
	OnSuccessDestroyOriginal bool               `json:"onSuccessDestroyOriginal,omitempty"`
	DestroyFromIfInState     string             `json:"destroyFromIfInState,omitempty"`
}

func (m *Copy) Name() string { return "Email/copy" }

func (m *Copy) Requires() []jmap.URI { return []jmap.URI{mail.URI} }

type CopyResponse struct {
	// The id of the account records were copied from.
	FromAccount jmap.ID `json:"fromAccountId,omitempty"`

	// The id of the account records were copied to.
	Account    jmap.ID                    `json:"accountId,omitempty"`
	OldState   string                     `json:"oldState,omitempty"`
	NewState   string                     `json:"newState,omitempty"`
	Created    map[jmap.ID]*Email         `json:"created,omitempty"`
	NotCreated map[jmap.ID]*jmap.SetError `json:"notCreated,omitempty"`
}

func newCopyResponse() jmap.MethodResponse { return &CopyResponse{} }
