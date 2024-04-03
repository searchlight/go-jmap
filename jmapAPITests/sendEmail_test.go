package jmapApi

import (
	"bytes"
	"fmt"
	_ "path"
	"testing"
)

func TestSendEmail(t *testing.T) {
	myMail, err := NewEmailBuilder().
		WithSubject("[UNIT TEST] Test Prometheus Alertmanager alert received").
		WithBodyValue("Test body").
		WithAttachment("podLogs.txt", bytes.NewReader([]byte("Test pod logs"))).
		WithAttachment("podLogs2.txt", bytes.NewReader([]byte("test logs 2"))).
		WithCustomHeader("customHeaderTest", "works").
		WithCustomHeader("anotherHeader", "works as well").
		WithRecipient("testuser1.org@mydomain").
		Build()

	/*	tmp, _ := myMail.MarshalJSON()
		fmt.Println(string(tmp[:]))*/

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
