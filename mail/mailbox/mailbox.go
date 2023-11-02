package mailbox

import "git.sr.ht/~rockorager/go-jmap"

func init() {
	jmap.RegisterMethod("Mailbox/get", newGetResponse)
	jmap.RegisterMethod("Mailbox/changes", newChangesResponse)
	jmap.RegisterMethod("Mailbox/query", newQueryResponse)
	jmap.RegisterMethod("Mailbox/queryChanges", newQueryChangesResponse)
	jmap.RegisterMethod("Mailbox/set", newSetResponse)
}

// A Mailbox represents a named set of Emails. This is the primary mechanism
// for organising Emails within an account. It is analogous to a folder or a
// label in other systems.
//
// See RFC8621, section 2.
type Mailbox struct {
	// The id of the Mailbox.
	ID jmap.ID `json:"id,omitempty"`

	// User-visible name for the Mailbox, e.g., “Inbox”.
	Name string `json:"name,omitempty"`

	// The Mailbox id for the parent of this Mailbox, or null if this
	// Mailbox is at the top level.
	ParentID jmap.ID `json:"parentId,omitempty"`

	// Identifies Mailboxes that have a particular common purpose (e.g.,
	// the “inbox”), regardless of the name property (which may be
	// localised).
	Role Role `json:"role,omitempty"`

	// Defines the sort order of Mailboxes when presented in the client’s
	// UI, so it is consistent between devices.
	//
	// A Mailbox with a lower order should be displayed before a Mailbox
	// with a higher order (that has the same parent) in any Mailbox
	// listing in the client’s UI.
	SortOrder uint64 `json:"sortOrder,omitempty"`

	// The number of Emails in this Mailbox.
	TotalEmails uint64 `json:"totalEmails,omitempty"`

	// The number of Emails in this Mailbox that have neither the $seen
	// keyword nor the $draft keyword.
	UnreadEmails uint64 `json:"unreadEmails,omitempty"`

	// The number of Threads where at least one Email in the Thread is in
	// this Mailbox.
	TotalThreads uint64 `json:"totalThreads,omitempty"`

	// An indication of the number of “unread” Threads in the Mailbox.
	UnreadThreads uint64 `json:"unreadThreads,omitempty"`

	// The set of rights (Access Control Lists (ACLs)) the user has in
	// relation to this Mailbox.
	Rights *Rights `json:"myRights,omitempty"`

	// true if the user indicated they wish to see this Mailbox in their
	// client.
	IsSubscribed bool `json:"isSubscribed,omitempty"`
}

// Access Control Lists (ACLs)
type Rights struct {
	MayReadItems bool `json:"mayReadItems,omitempty"`

	MayAddItems bool `json:"mayAddItems,omitempty"`

	MayRemoveItems bool `json:"mayRemoveItems,omitempty"`

	MaySetSeen bool `json:"maySetSeen,omitempty"`

	MaySetKeywords bool `json:"maySetKeywords,omitempty"`

	MayCreateChild bool `json:"mayCreateChild,omitempty"`

	MayRename bool `json:"mayRename,omitempty"`

	MayDelete bool `json:"mayDelete,omitempty"`

	MaySubmit bool `json:"maySubmit,omitempty"`
}

// Identifies Mailboxes that have a particular common purpose (e.g., the
// “inbox”), regardless of the name property (which may be localised).
type Role string

const (
	RoleAll Role = "all"

	RoleArchive Role = "archive"

	RoleDrafts Role = "drafts"

	RoleFlagged Role = "flagged"

	RoleHasChildren Role = "haschildren"

	RoleHasNoChildren Role = "hasnochildren"

	RoleImportant Role = "important"

	RoleInbox Role = "inbox"

	RoleJunk Role = "junk"

	RoleMarked Role = "marked"

	RoleNoInferiors Role = "noinferiors"

	RoleNonExistent Role = "nonexistent"

	RoleNoSelect Role = "noselect"

	RoleRemote Role = "remote"

	RoleSent Role = "sent"

	RoleSubscribed Role = "subscribed"

	RoleTrash Role = "trash"

	RoleUnmarked Role = "unmarked"
)
