package jmapApi

import (
	"fmt"
	"git.sr.ht/~rockorager/go-jmap"
	_ "git.sr.ht/~rockorager/go-jmap/core"
	"git.sr.ht/~rockorager/go-jmap/mail/email"
)

func IsLatestProblemSolved(customHeader string) bool {
	req := &jmap.Request{
		Using: []jmap.URI{"urn:ietf:params:jmap:core", "urn:ietf:params:jmap:mail"},
	}
	inboxMailboxId, _ := getMailboxIdByTag("inbox")
	myfilter := email.FilterCondition{
		Subject:   "THIS IS THE SUBJECT",
		InMailbox: inboxMailboxId,
		Header:    append([]string{}, customHeader),
	}

	mySortComparator := email.SortComparator{
		Property:    "receivedAt",
		IsAscending: false,
	}

	callID := req.Invoke(&email.Query{
		Account: userID,
		Filter:  &myfilter,
		Sort:    []*email.SortComparator{&mySortComparator},
		Limit:   10,
	})
	req.Invoke(&email.Get{
		Account: userID,
		ReferenceIDs: &jmap.ResultReference{
			ResultOf: callID,        // The CallID of the referenced method
			Name:     "Email/query", // The name of the referenced method
			Path:     "/ids",        // JSON pointer to the location of the reference
		},
		Properties: append([]string{}, "headers"),
	})
	resp, _ := myClient.Do(req)

	flag := false
	// Searching is there any header name 'customHeader' with value ' solved'
	for _, inv := range resp.Responses {
		switch r := inv.Args.(type) {
		case *email.GetResponse:
			for _, eml := range r.List {
				for _, header := range eml.Headers {
					fmt.Println(header.Name, " ", header.Value)
					if header.Name == customHeader && header.Value == " solved" {
						flag = true
						break
					}
				}
			}
		}
	}
	return flag
}
