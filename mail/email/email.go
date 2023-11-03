package email

import (
	"time"

	"git.sr.ht/~rockorager/go-jmap"
	"git.sr.ht/~rockorager/go-jmap/mail"
)

func init() {
	jmap.RegisterCapability(&smimeVerify{})
	jmap.RegisterMethod("Email/get", newGetResponse)
	jmap.RegisterMethod("Email/changes", newChangesResponse)
	jmap.RegisterMethod("Email/query", newQueryResponse)
	jmap.RegisterMethod("Email/queryChanges", newQueryChangesResponse)
	jmap.RegisterMethod("Email/set", newSetResponse)
	jmap.RegisterMethod("Email/copy", newCopyResponse)
	jmap.RegisterMethod("Email/import", newImportResponse)
	jmap.RegisterMethod("Email/parse", newParseResponse)
}

type Email struct {
	// The ID of the Email. Note: this is _not_ the Message-ID
	//
	// immutable;server-set
	ID jmap.ID `json:"id,omitempty"`

	// The ID of the raw RFC5322 message
	//
	// immutable;server-set
	BlobID jmap.ID `json:"blobId,omitempty"`

	// The id of the Thread to which this Email belongs.
	//
	// immutable;server-set
	ThreadID jmap.ID `json:"threadId,omitempty"`

	// The set of Mailbox ids this Email belongs to. An Email in the mail
	// store MUST belong to one or more Mailboxes at all times (until it
	// is destroyed).
	MailboxIDs map[jmap.ID]bool `json:"mailboxIds,omitempty"`

	// A set of keywords that apply to the Email. Each key must have an
	// associated value of "true"
	Keywords map[string]bool `json:"keywords,omitempty"`

	// The size, in bytes, of the message
	//
	// immutable;server-set
	Size uint64 `json:"size,omitempty"`

	// The date the Email was received. Equivalent to INTERNAL_DATE in IMAP
	//
	// immutable
	ReceivedAt *time.Time `json:"receivedAt,omitempty"`

	// This is a list of all header fields, in the same order they appear in
	// the message.
	//
	// immutable
	Headers []*Header `json:"headers,omitempty"`

	// The Message-ID of the email. For conforming messages, this will be
	// len() == 1
	//
	// immutable
	MessageID []string `json:"messageId,omitempty"`

	// immutable
	InReplyTo []string `json:"inReplyTo,omitempty"`

	// immutable
	References []string `json:"references,omitempty"`

	// immutable
	Sender []*mail.Address `json:"sender,omitempty"`

	// immutable
	From []*mail.Address `json:"from,omitempty"`

	// immutable
	To []*mail.Address `json:"to,omitempty"`

	// immutable
	CC []*mail.Address `json:"cc,omitempty"`

	// immutable
	BCC []*mail.Address `json:"bcc,omitempty"`

	// immutable
	ReplyTo []*mail.Address `json:"replyTo,omitempty"`

	// immutable
	Subject string `json:"subject,omitempty"`

	// SentAt is the Date header value
	//
	// immutable
	SentAt *time.Time `json:"sentAt,omitempty"`

	// This is the full MIME structure of the message body, without
	// recursing into message/rfc822 or message/global parts.
	//
	// immutable
	BodyStructure *BodyPart `json:"bodyStructure,omitempty"`

	// This is a map of partId to an EmailBodyValue object for none, some,
	// or all text/* parts. Which parts are included and whether the value
	// is truncated is determined by various arguments to Email/get and
	// Email/parse.
	//
	// immutable
	BodyValues map[string]*BodyValue `json:"bodyValues,omitempty"`

	// A list of text/plain, text/html, image/*, audio/*, and/or video/*
	// parts to display (sequentially) as the message body, with a
	// preference for text/plain when alternative versions are available.
	//
	// immutable
	TextBody []*BodyPart `json:"textBody,omitempty"`

	// A list of text/plain, text/html, image/*, audio/*, and/or video/*
	// parts to display (sequentially) as the message body, with a
	// preference for text/html when alternative versions are available.
	//
	// immutable
	HTMLBody []*BodyPart `json:"htmlBody,omitempty"`

	// A list, traversing depth-first, of all parts in bodyStructure that
	// satisfy either of the following conditions:
	//
	//     not of type multipart/* and not included in textBody or htmlBody
	//
	//     of type image/*, audio/*, or video/* and not in both textBody
	//     and htmlBody
	//
	// immutable
	Attachments []*BodyPart `json:"attachments,omitempty"`

	// immutable;server-set
	HasAttachment bool `json:"hasAttachment,omitempty"`

	// A plaintext fragment of the message body.	// This MUST NOT be more
	// than 256 characters in length.
	//
	// immutable;server-set
	Preview string `json:"preview,omitempty"`

	// If empty, there is no S/MIME signature. Otherwise will be one of the
	// following strings
	// - "unknown" - Can be returned for OpenPGP signed messages
	// - "signed" - S/MIME signed but not yet verified
	// - "signed/verified" - Signed and verified per RFC8551 and RFC8550
	// - "signed/failed"
	// - "encrypted+signed/verified"
	// - "encrypted+signed/failed"
	//
	// server-set
	SMIMEStatus string `json:"smimeStatus,omitempty"`

	// If empty, there is no S/MIME signature. Otherwise will be one of the
	// following strings, and represents the status at time of delivery
	// - "unknown" - Can be returned for OpenPGP signed messages
	// - "signed" - S/MIME signed but not yet verified
	// - "signed/verified" - Signed and verified per RFC8551 and RFC8550
	// - "signed/failed"
	// - "encrypted+signed/verified"
	// - "encrypted+signed/failed"
	//
	// server-set
	SMIMEStatusAtDelivery string `json:"smimeStatusAtDelivery,omitempty"`

	// If empty, no errors or no signature. Otherwise, this will contain any
	// errors during verification of SMIME properties
	//
	// server-set
	SMIMEErrors []string `json:"smimeErrors,omitempty"`

	// If empty, no signature or not verified. Otherwise, this is the time
	// the signature was most recently verified
	//
	// server-set
	SMIMEVerifiedAt *time.Time `json:"smimeVerifiedAt,omitempty"`
}

type AddressGroup struct {
	// The display-name of the group
	Name      string          `json:"name,omitempty"`
	Addresses []*mail.Address `json:"addresses,omitempty"`
}

type Header struct {
	// The header field name, with the same capitalization that it has in
	// the message.
	Name string `json:"name,omitempty"`

	// The header field value in Raw form.
	Value string `json:"value,omitempty"`
}

type BodyPart struct {
	// Identifies this part uniquely within the Email. This is scoped to
	// the emailId and has no meaning outside of the JMAP Email object
	//
	// Multipart messages do not have a PartID
	PartID string `json:"partId,omitempty"`

	// The Blob ID representing this Part
	BlobID jmap.ID `json:"blobId,omitempty"`

	// The number of bytes the user would download
	Size uint64 `json:"size,omitempty"`

	// This is a list of all header fields in the part, in the order they
	// appear in the message. The values are in Raw form.
	Headers []*Header `json:"headers,omitempty"`

	// The filename associated with this Part, if given
	Name string `json:"name,omitempty"`

	// The value of the Content-Type header field of the part, if present;
	// otherwise, the implicit type as per the MIME standard (text/plain or
	// message/rfc822 if inside a multipart/digest)
	Type string `json:"type,omitempty"`

	// The value of the charset parameter of the Content-Type header
	// field, if present, or null if the header field is present but not
	// of type text/*. If there is no Content-Type header field, or it
	// exists and is of type text/* but has no charset parameter, this is
	// the implicit charset as per the MIME standard: us-ascii.
	Charset string `json:"charset,omitempty"`

	// The value of the Content-Disposition header field of the part, if
	// present; otherwise, it’s null. CFWS is removed and any parameters
	// are stripped.
	Disposition string `json:"disposition,omitempty"`

	// The value of the Content-Id header field of the part, if present;
	// otherwise it’s null.
	CID string `json:"cid,omitempty"`

	// The list of language tags, as defined in RFC3282, in the
	// Content-Language header field of the part, if present.
	Language []string `json:"language,omitempty"`

	// The URI, as defined in RFC2557, in the Content-Location header field
	// of the part, if present.
	Location string `json:"location,omitempty"`

	// If the type is multipart/*, this contains the body parts of each
	// child.
	SubParts []*BodyPart `json:"subParts,omitempty"`
}

type BodyValue struct {
	// The value of the BodyValue
	Value string `json:"value,omitempty"`

	// True if there was an encoding problem
	IsEncodingProblem bool `json:"isEncodingProblem,omitempty"`

	// This is true if the value has been truncated
	IsTruncated bool `json:"isTruncated"`
}
