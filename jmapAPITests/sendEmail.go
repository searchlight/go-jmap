package jmapApi

import (
	"errors"
	"git.sr.ht/~rockorager/go-jmap"
	_ "git.sr.ht/~rockorager/go-jmap/core"
	"git.sr.ht/~rockorager/go-jmap/mail/email"
	"git.sr.ht/~rockorager/go-jmap/mail/emailsubmission"
)

func SendEmail(myMail *email.Email) (jmap.ID, error) {
	req := &jmap.Request{
		Using: []jmap.URI{"urn:ietf:params:jmap:core", "urn:ietf:params:jmap:mail"},
	}

	invokeSetDraftEMail(req, userID, myMail)
	invokeSendEmail(req, userID)

	resp, err := myClient.Do(req)
	if err != nil {
		return "", err
	}

	var threadID jmap.ID

	for _, inv := range resp.Responses {
		if threadID != "" {
			break
		}

		switch r := inv.Args.(type) {
		case *email.SetResponse:
			if emailResponse, ok := r.Created["draft"]; ok {
				threadID = emailResponse.ThreadID
			}
		}
	}

	if threadID == "" {
		return "", errors.New("email submission failed / not created")
	}

	return threadID, nil
}

func invokeSetDraftEMail(req *jmap.Request, id jmap.ID, myMail *email.Email) {
	myMap := map[jmap.ID]*email.Email{
		"draft": &(*myMail),
	}

	req.Invoke(&email.Set{
		Account: id,
		Create:  myMap,
	})
}

func invokeSendEmail(req *jmap.Request, id jmap.ID) {
	myEmailSubmission := emailsubmission.EmailSubmission{
		EmailID: "#draft",
	}

	req.Invoke(&emailsubmission.Set{
		Account: id,

		Create: map[jmap.ID]*emailsubmission.EmailSubmission{
			"sendIt": &myEmailSubmission,
		},

		OnSuccessUpdateEmail: map[jmap.ID]jmap.Patch{
			"#sendIt": {
				"mailboxIds/" + string(draftMailboxID): nil,
				"mailboxIds/" + string(sentMailboxID):  true,
				"keywords/$seen":                       nil,
				"keywords/$draft":                      nil,
			},
		},
	})
}
