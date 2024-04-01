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
	uploadResponse                *jmap.UploadResponse
}

func NewEmailBuilder() *emailBuilder {
	ret := &emailBuilder{}
	ret.customHeaders = []*email.Header{}
	return ret
}

func (b *emailBuilder) SetSubject(subject string) *emailBuilder {
	b.subject = subject
	return b
}

func (b *emailBuilder) SetBodyValue(body string) *emailBuilder {
	b.bodyValue = body
	return b
}

func (b *emailBuilder) SetRecipient(recipient string) *emailBuilder {
	b.recipient = recipient
	return b
}

func (b *emailBuilder) SetCustomHeader(key string, value string) *emailBuilder {
	key = "header:" + key
	customHeader := email.Header{
		Name:  key,
		Value: value,
	}

	b.customHeaders = []*email.Header{&customHeader}
	return b
}

func (b *emailBuilder) SetAttachment(blob io.Reader) *emailBuilder {
	resp, err := myClient.Upload(userID, blob)
	b.uploadResponse = nil

	if err != nil {
		fmt.Println("Error setting attachment ", err.Error())
		return b
	}

	if resp == nil {
		fmt.Println("response is nil")
		return b
	}

	b.uploadResponse = resp
	return b
}

func (b *emailBuilder) Build() (email.Email, error) {
	if b.recipient == "" {
		return email.Email{}, errors.New("No recipient defined")
	}

	if b.uploadResponse == nil {
		return email.Email{}, errors.New("Error setting attachment")
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

	myAttachment := email.BodyPart{
		BlobID: b.uploadResponse.ID,

		Size: b.uploadResponse.Size,

		Type: b.uploadResponse.Type,

		Name: "pod_Logs.txt",

		Disposition: "attachment",
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

		Attachments: []*email.BodyPart{&myAttachment},
	}

	return myMail, nil
}
