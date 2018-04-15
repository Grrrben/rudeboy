package outgoing

import (
	"fmt"

	"slackbot/incoming"
	"slackbot/storage"

	"github.com/nlopes/slack"
)

var Storage storage.Storager

type Messenger struct {
	Api *slack.RTM
}

func (r Messenger) ReactOnCall(msg incoming.Message, ch string) {
	switch msg.Action {
	case "help":
		t := "This is the help text\n" +
			"Use _add_, or _burn_ followed by an @userhandle and to add, a burn as well. E.g.\n" +
			"@bot add @username is a fool\n" +
			"@bot burn @username"
		r.Api.SendMessage(r.Api.NewOutgoingMessage(t, ch))

	case "add":
		Storage.Add(msg.TargetUser, msg.Text)
		s := fmt.Sprintf("Added, thnx <@%s>", msg.Sender)
		r.Api.SendMessage(r.Api.NewOutgoingMessage(s, ch))

	case "burn":
		r.Api.SendMessage(r.Api.NewOutgoingMessage(Storage.Get(msg.TargetUser), ch))

	case "delete":
		r.Api.SendMessage(r.Api.NewOutgoingMessage("Delete is not working yet, sorry", ch))

	default:
		r.Api.SendMessage(r.Api.NewOutgoingMessage("I do not understand", ch))
	}
}
