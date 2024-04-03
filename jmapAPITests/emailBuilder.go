package jmapApi

import (
	"errors"
	"fmt"
	"git.sr.ht/~rockorager/go-jmap"
	_ "git.sr.ht/~rockorager/go-jmap/core"
	"git.sr.ht/~rockorager/go-jmap/mail"
	"git.sr.ht/~rockorager/go-jmap/mail/email"
	"io"
)

type emailBuilder struct {
	recipient, subject, bodyValue string
	customHeaders                 []*email.Header
	attachments                   []*email.BodyPart
}

func NewEmailBuilder() *emailBuilder {
	return &emailBuilder{}
}

func (b *emailBuilder) WithSubject(subject string) *emailBuilder {
	b.subject = subject
	return b
}

func (b *emailBuilder) WithBodyValue(body string) *emailBuilder {
	b.bodyValue = body
	return b
}

func (b *emailBuilder) WithRecipient(recipient string) *emailBuilder {
	b.recipient = recipient
	return b
}

func (b *emailBuilder) WithCustomHeader(key string, value string) *emailBuilder {
	customHeader := email.Header{
		Name:  key,
		Value: value,
	}

	b.customHeaders = append(b.customHeaders, &customHeader)
	return b
}

func (b *emailBuilder) WithAttachment(attachmentName string, blob io.Reader) *emailBuilder {
	resp, err := myClient.Upload(userID, blob)

	if err != nil {
		fmt.Println("Error setting attachment ", err.Error())
		return b
	}

	if resp == nil {
		fmt.Println("response is nil")
		return b
	}

	myAttachment := email.BodyPart{
		BlobID: resp.ID,

		Size: resp.Size,

		Type: resp.Type,

		Name: attachmentName,

		Disposition: "attachment",
	}

	b.attachments = append(b.attachments, &myAttachment)
	return b
}

func (b *emailBuilder) Build() (email.Email, error) {
	if b.recipient == "" {
		return email.Email{}, errors.New("No recipient defined")
	}

	from := mail.Address{
		Name:  userEmail,
		Email: userEmail,
	}

	to := mail.Address{
		Name:  b.recipient,
		Email: b.recipient,
	}

	myBodyValue := email.BodyValue{
		Value: b.bodyValue,
	}

	myBodyPart := email.BodyPart{
		PartID: "body",
		Type:   "text/plain",
	}

	myMail := email.Email{
		CustomHeaders: b.customHeaders,

		From: []*mail.Address{
			&from,
		},

		To: []*mail.Address{
			&to,
		},

		Subject: b.subject,

		Keywords: map[string]bool{"$draft": true},

		MailboxIDs: map[jmap.ID]bool{draftMailboxID: true},

		BodyValues: map[string]*email.BodyValue{"body": &myBodyValue},

		TextBody: []*email.BodyPart{&myBodyPart},

		HasAttachment: true,

		Attachments: b.attachments,
	}

	return myMail, nil
}
