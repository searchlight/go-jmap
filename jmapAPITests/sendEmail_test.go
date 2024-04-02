package jmapApi

import (
	"bytes"
	"fmt"
	_ "path"
	"testing"
)

func TestSendEmail(t *testing.T) {
	myMail, err := NewEmailBuilder().
		SetSubject("[UNIT TEST] Test Prometheus Alertmanager alert received").
		SetBodyValue("Test body").
		SetAttachment(bytes.NewReader([]byte("Test pod logs"))).
		SetCustomHeader("customHeaderTest", "works").
		SetCustomHeader("anotherHeader", "works as well").
		SetRecipient("testuser1.org@mydomain").
		Build()

	if err != nil {
		t.Error("Error creating email ", err)
		return
	}

	threadID, err2 := SendEmail(&myMail)

	if err2 != nil {
		t.Error("Error Sending email ", err)
		return
	}

	fmt.Println("server returned threadID: ", threadID)
}
