package email

import (
	"encoding/json"
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

	// A set of keywords that apply to the Email. The set is represented as
	// an object, with the keys being the keywords. The value for each key
	// in the object MUST be true.
	Keywords map[string]bool `json:"keywords,omitempty"`

	// The size, in octets, of the raw data for the message [@!RFC5322] (as
	// referenced by the blobId, i.e., the number of octets in the file the
	// user would download).
	//
	// immutable;server-set
	Size uint64 `json:"size,omitempty"`

	// The date the Email was received by the message store. This is the
	// internal date in IMAP [@?RFC3501].
	//
	// immutable
	ReceivedAt *time.Time `json:"receivedAt,omitempty"`

	// This is a list of all header fields [@!RFC5322], in the same order
	// they appear in the message.
	//
	// immutable
	Headers []*Header `json:"headers,omitempty"`

	// For adding custom headers to the email
	CustomHeaders []*Header `json:"omitempty"`

	// The value is identical to the value of
	// header:Message-ID:asMessageIds. For messages conforming to RFC 5322
	// this will be an array with a single entry.
	//
	// immutable
	MessageID []string `json:"messageId,omitempty"`

	// The value is identical to the value of
	// header:In-Reply-To:asMessageIds.
	//
	// immutable
	InReplyTo []string `json:"inReplyTo,omitempty"`

	// The value is identical to the value of
	// header:References:asMessageIds.mailAccount
	//
	// immutable
	References []string `json:"references,omitempty"`

	// The value is identical to the value of header:Sender:asAddresses.
	//
	// immutable
	Sender []*mail.Address `json:"sender,omitempty"`

	// The value is identical to the value of header:From:asAddresses.
	//
	// immutable
	From []*mail.Address `json:"from,omitempty"`

	// The value is identical to the value of header:To:asAddresses.
	//
	// immutable
	To []*mail.Address `json:"to,omitempty"`

	// The value is identical to the value of header:Cc:asAddresses.
	//
	// immutable
	CC []*mail.Address `json:"cc,omitempty"`

	// The value is identical to the value of header:Bcc:asAddresses.
	//
	// immutable
	BCC []*mail.Address `json:"bcc,omitempty"`

	// The value is identical to the value of header:Reply-To:asAddresses.
	//
	// immutable
	ReplyTo []*mail.Address `json:"replyTo,omitempty"`

	// The value is identical to the value of header:Subject:asText.
	//
	// immutable
	Subject string `json:"subject,omitempty"`

	// The value is identical to the value of header:Date:asDate.
	//
	// immutable
	SentAt *time.Time `json:"sentAt,omitempty"`

	// This is the full MIME structure of the message body, without
	// recursing into message/rfc822 or message/global parts. Note that
	// EmailBodyParts may have subParts if they are of type multipart/*.
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
	// None of these parts include subParts, including message/* types.
	// Attached messages may be fetched using the Email/parse method and
	// the blobId.
	//
	// Note that a text/html body part HTML may reference image parts in
	// attachments by using cid: links to reference the Content-Id, as
	// defined in [@!RFC2392], or by referencing the Content-Location.
	//
	// immutable
	Attachments []*BodyPart `json:"attachments,omitempty"`

	// This is true if there are one or more parts in the message that a
	// client UI should offer as downloadable. A server SHOULD set
	// hasAttachment to true if the attachments list contains at least one
	// item that does not have Content-Disposition: inline. The server MAY
	// ignore parts in this list that are processed automatically in some
	// way or are referenced as embedded images in one of the text/html
	// parts of the message.
	//
	// The server MAY set hasAttachment based on implementation-defined or
	// site-configurable heuristics.
	//
	// immutable;server-set
	HasAttachment bool `json:"hasAttachment,omitempty"`

	// A plaintext fragment of the message body. This is intended to be
	// shown as a preview line when listing messages in the mail store and
	// may be truncated when shown. The server may choose which part of the
	// message to include in the preview; skipping quoted sections and
	// salutations and collapsing white space can result in a more useful
	// preview.
	//
	// This MUST NOT be more than 256 characters in length.
	//
	// As this is derived from the message content by the server, and the
	// algorithm for doing so could change over time, fetching this for an
	// Email a second time MAY return a different result. However, the
	// previous value is not considered incorrect, and the change SHOULD
	// NOT cause the Email object to be considered as changed by the
	// server.
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
	// The display-name of the group [@!RFC5322], or null if the addresses
	// are not part of a group. If this is a quoted-string, it is processed
	// the same as the name in the EmailAddress type.
	Name string `json:"name,omitempty"`

	// The mailbox values that belong to this group, represented as
	// EmailAddress objects.
	Addresses []*mail.Address `json:"addresses,omitempty"`
}

func (myMail *Email) MarshalJSON() ([]byte, error) {
	mailMap := make(map[string]interface{})

	mailMap["mailboxIds"] = myMail.MailboxIDs
	mailMap["keywords"] = myMail.Keywords
	mailMap["from"] = myMail.From
	mailMap["to"] = myMail.To
	if myMail.InReplyTo != nil {
		mailMap["inReplyTo"] = myMail.InReplyTo
	}
	if myMail.InReplyTo != nil {
		mailMap["replyTo"] = myMail.ReplyTo
	}
	if myMail.Sender != nil {
		mailMap["sender"] = myMail.Sender
	}
	if myMail.CC != nil {
		mailMap["cc"] = myMail.CC
	}
	if myMail.BCC != nil {
		mailMap["bcc"] = myMail.BCC
	}
	if &myMail.Subject != nil {
		mailMap["subject"] = myMail.Subject
	}
	if myMail.References != nil {
		mailMap["references"] = myMail.References
	}
	if myMail.ReceivedAt != nil {
		mailMap["receivedAt"] = myMail.ReceivedAt
	}
	if myMail.SentAt != nil {
		mailMap["sentAt"] = myMail.SentAt
	}
	if myMail.Headers != nil {
		mailMap["headers"] = myMail.Headers
	}
	if myMail.BodyStructure != nil {
		mailMap["bodyStructure"] = myMail.BodyStructure
	}
	if myMail.BodyValues != nil {
		mailMap["bodyValues"] = myMail.BodyValues
	}
	if myMail.TextBody != nil {
		mailMap["textBody"] = myMail.TextBody
	}
	if myMail.HTMLBody != nil {
		mailMap["htmlBody"] = myMail.HTMLBody
	}
	if myMail.Attachments != nil {
		mailMap["attachments"] = myMail.Attachments
	}
	if &myMail.HasAttachment != nil {
		mailMap["hasAttachment"] = myMail.HasAttachment
	}
	if &myMail.Preview != nil {
		mailMap["preview"] = myMail.Preview
	}

	for _, header := range myMail.CustomHeaders {
		mailMap[header.Name] = header.Value
	}

	return json.Marshal(mailMap)
}

type Header struct {
	// The header field name as defined in [@!RFC5322], with the same
	// capitalization that it has in the message.
	Name string `json:"name,omitempty"`

	// The header field value as defined in [@!RFC5322], in Raw form.
	Value string `json:"value,omitempty"`
}

type HeaderSlice []*Header

// These properties are derived from the message body [@!RFC5322] and its MIME
// entities [@RFC2045].
type BodyPart struct {
	// Identifies this part uniquely within the Email. This is scoped to
	// the emailId and has no meaning outside of the JMAP Email object
	// representation. This is null if, and only if, the part is of type
	// multipart/*.
	PartID string `json:"partId,omitempty"`

	// The id representing the raw octets of the contents of the part,
	// after decoding any known Content-Transfer-Encoding (as defined in
	// [@!RFC2045]), or null if, and only if, the part is of type
	// multipart/*. Note that two parts may be transfer-encoded differently
	// but have the same blob id if their decoded octets are identical and
	// the server is using a secure hash of the data for the blob id. If
	// the transfer encoding is unknown, it is treated as though it had no
	// transfer encoding.
	BlobID jmap.ID `json:"blobId,omitempty"`

	// The size, in octets, of the raw data after content transfer decoding
	// (as referenced by the blobId, i.e., the number of octets in the file
	// the user would download).
	Size uint64 `json:"size,omitempty"`

	// This is a list of all header fields in the part, in the order they
	// appear in the message. The values are in Raw form.
	Headers []*Header `json:"headers,omitempty"`

	// This is the decoded filename parameter of the Content-Disposition
	// header field per [@!RFC2231], or (for compatibility with existing
	// systems) if not present, then it’s the decoded name parameter of the
	// Content-Type header field per [@!RFC2047].
	Name string `json:"name,omitempty"`

	// The value of the Content-Type header field of the part, if present;
	// otherwise, the implicit type as per the MIME standard (text/plain or
	// message/rfc822 if inside a multipart/digest). CFWS is removed and
	// any parameters are stripped.
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
	// otherwise it’s null. CFWS and surrounding angle brackets (<>) are
	// removed. This may be used to reference the content from within a
	// text/html body part HTML using the cid: protocol, as defined in
	// [@!RFC2392].
	CID string `json:"cid,omitempty"`

	// The list of language tags, as defined in [@!RFC3282], in the
	// Content-Language header field of the part, if present.
	Language []string `json:"language,omitempty"`

	// The URI, as defined in [@!RFC2557], in the Content-Location header
	// field of the part, if present.
	Location string `json:"location,omitempty"`

	// If the type is multipart/*, this contains the body parts of each
	// child.
	SubParts []*BodyPart `json:"subParts,omitempty"`
}

// This is a map of partId to an BodyValue object for none, some, or all
// text/* parts. Which parts are included and whether the value is truncated is
// determined by various arguments to Email/get and Email/parse.
type BodyValue struct {
	// The value of the body part after decoding Content-Transfer-Encoding
	// and the Content-Type charset, if both known to the server, and with
	// any CRLF replaced with a single LF. The server MAY use heuristics to
	// determine the charset to use for decoding if the charset is unknown,
	// no charset is given, or it believes the charset given is incorrect.
	// Decoding is best effort; the server SHOULD insert the unicode
	// replacement character (U+FFFD) and continue when a malformed section
	// is encountered.
	//
	// Note that due to the charset decoding and line ending normalisation,
	// the length of this string will probably not be exactly the same as
	// the size property on the corresponding EmailBodyPart.
	Value string `json:"value,omitempty"`

	// This is true if malformed sections were found while decoding the
	// charset, or the charset was unknown, or the
	// content-transfer-encoding was unknown.
	IsEncodingProblem bool `json:"isEncodingProblem,omitempty"`

	// This is true if the value has been truncated
	IsTruncated bool `json:"isTruncated"`
}
