package incoming

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/nlopes/slack"
)

var action = []string{"add", "burn", "delete", "help"}

type Message struct {
	Action     string
	Sender     string
	TargetUser string
	Text       string
}

func Disect(msgEvent *slack.MessageEvent, handle string) (msg Message) {

	fmt.Println("Message:")
	fmt.Println(msgEvent.Text)
	// <@UA5JF5R29> add <@UA5EDQCSW> is gek

	msg.Sender = msgEvent.User

	text := strings.TrimSpace(strings.TrimPrefix(msgEvent.Text, handle))

	// what type of action?
	act, content := substractActionAndContent(text)
	if act != "" {
		msg.Action = act
		msg.Text = content
	}

	// is there a target user incoming the incoming
	reUserHandle := regexp.MustCompile("<@[A-Z0-9]{9}>")
	t := reUserHandle.FindString(text)

	if t != "" {
		msg.TargetUser = t
	}

	return
}

func substractActionAndContent(str string) (string, string) {
	if str == "help" {
		return "help", ""
	}
	for i := range str {
		if str[i] == ' ' {
			act, rest := str[:i], str[i:]
			for _, v := range action {
				if v == act {
					return act, rest
				}
			}
		}
	}
	return "", str
}
