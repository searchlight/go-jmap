package jmapApi

import (
	"errors"
	"fmt"
	"git.sr.ht/~rockorager/go-jmap"
	"git.sr.ht/~rockorager/go-jmap/mail/mailbox"
	"strconv"
)

// ##### UNTESTED #####
func getAllMailboxes() ([]*mailbox.Mailbox, error) {
	req := &jmap.Request{
		Using: []jmap.URI{"urn:ietf:params:jmap:core", "urn:ietf:params:jmap:mail"},
	}

	req.Invoke(&mailbox.Get{
		Account: userID,
	})

	resp, err := myClient.Do(req)
	if err != nil || resp == nil {
		fmt.Println("Error fetching mailboxes")
		return []*mailbox.Mailbox{}, err
	}

	if len(resp.Responses) > 1 {
		return []*mailbox.Mailbox{}, errors.New("Multiple responses received on function getAllMailboxes")
	}

	///len(resp.Responses) == 1
	var mailboxList []*mailbox.Mailbox

	if getResp := resp.Responses[0].Args.(*mailbox.GetResponse); getResp == nil {
		return []*mailbox.Mailbox{}, errors.New("get response is null")
	} else {
		mailboxList = getResp.List
	}

	return mailboxList, nil
}

func getMailboxIdByTag(tag string) (jmap.ID, error) {
	req := &jmap.Request{
		Using: []jmap.URI{"urn:ietf:params:jmap:core", "urn:ietf:params:jmap:mail"},
	}

	myFilter := mailbox.FilterCondition{
		Role: mailbox.Role(tag),
	}

	req.Invoke(&mailbox.Query{
		Account: userID,
		Filter:  &myFilter,
	})

	resp, err := myClient.Do(req)
	if err != nil || resp == nil {
		fmt.Println("error fetching mailbox with tag: " + tag)
		return "", err
	}

	if len(resp.Responses) > 1 {
		return "", errors.New("Multiple responses received on function getMailboxIdByTag with tag: " + tag)
	}

	///len(resp.Responses) == 1
	var requiredID jmap.ID

	if queryResp := resp.Responses[0].Args.(*mailbox.QueryResponse); queryResp == nil {
		return "", errors.New("query response is null")
	} else if mailBoxCount := len(queryResp.IDs); mailBoxCount != 1 {
		return "", errors.New("expected exactly 1 mailbox with the role: " + tag + ", found " + strconv.FormatInt(int64(mailBoxCount), 10))
	} else {
		requiredID = queryResp.IDs[0]
	}

	return requiredID, nil
}
