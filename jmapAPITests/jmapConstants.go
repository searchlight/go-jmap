package jmapApi

import (
	"git.sr.ht/~rockorager/go-jmap"
	"git.sr.ht/~rockorager/go-jmap/mail"
	"log"
)

const sessionEndpoint = "http://james.appscode.ninja:80/jmap/session"
const bearerToken = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZXN0dXNlcjEub3JnQG15ZG9tYWluIiwiaWF0IjoxNjE2MjM5MDIyLCJleHAiOjE5NDYyMzkwMjJ9.ZQuDcCFyccpSvAtp0M88TeATePlx4Gvs97_asg75cn46Xkoe-ENVYM-N_ZeWs8G7joxCQS2R2BOQ1LPXY7rFYubEA1Ne8xMMJvS6LW8NktdYrSUZm3J7oyeC6K42cHqYlspK2RaiUl-XKaxy0J38p9MhIKvzOETTEpaABFNp5MplNqkVoAnYnOnF8Ez7AfoCnieBuxmR_qmZpRNdl6XXpoyoihv0azY2RxAyk6uRrnNbzirICuthnkFdllPrbffTcRk9NXkV9PyGlmdAVczxeUADIpd2WMWyoB3yFVe4s1dl9UQp7QsJqpBUUuY3LhE295PwzGDq5WKMWzdII55N2BnVe9LS33GXMU9oGgWqpBmxrTzSwc-jiewhorcc32keXkrxhtxL1uNXuYzg2lqFv7he7aTKWowkI2MJVXLhUN-xpIuaKgAgWY99VBYcf5gfhLfwSYVVZUKIQmYI7Jh0bWNoTHyXoHj0wyFnzNe3SsPECip2WohIDD7_CNLydGDvFmdVp7_aIDAHpiA2BrLysm2Go9GHy-ESdVsma6aG95EeICp1ngLSFHx4IOQ3Oisr4t4DStv7INQQQDkuI6hjhVctbt75rcTa-541a237N3RuBKhOPhEbf3eipUD2t-x28Wi_EHDQpFb0D0nZjvI8Q0AoNm0K9rDzVKvUiyNSIr8"

var draftMailboxID jmap.ID
var sentMailboxID jmap.ID
var userEmail string

var myClient *jmap.Client
var userID jmap.ID

func init() {
	/*	envErr := godotenv.Load()
		if envErr != nil {
			log.Fatal("Error loading .env file")
		}
	*/
	myClient = &jmap.Client{
		SessionEndpoint: sessionEndpoint,
	}

	myClient.WithAccessToken(bearerToken)

	if err := myClient.Authenticate(); err != nil {
		log.Fatal("unable to authenticate user with the given credentials", err)
	}

	userID = myClient.Session.PrimaryAccounts[mail.URI]
	userEmail = myClient.Session.Accounts[userID].Name

	///Not using RFC8621's "Mailbox/query filter" because it's not properly implemented in the James distributed server 3.8.0
	var err error
	if draftMailboxID, err = getMailboxIdByTag("Drafts"); err != nil {
		log.Fatal(err)
	}

	if sentMailboxID, err = getMailboxIdByTag("Sent"); err != nil {
		log.Fatal(err)
	}
}
