package main

import (
	"fmt"
	"os"
	"strings"

	"slackbot/incoming"
	"slackbot/outgoing"
	"slackbot/storage/bolt"

	"path/filepath"

	"github.com/nlopes/slack"
)

var messenger outgoing.Messenger

func main() {
	api := slack.New(getToken())
	api.SetDebug(false)

	rtm := api.NewRTM()
	messenger.Api = rtm
	go rtm.ManageConnection()

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}

	b := &bolt.Bolt{}
	b.Path = filepath.Dir(ex) + "/db/bolt.db"
	outgoing.Storage = b
	//outgoing.Storage = &memory.Cache{}

	listen(rtm)
}

func listen(rtm *slack.RTM) {
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.MessageEvent:
				info := rtm.GetInfo()
				handle := fmt.Sprintf("<@%s> ", info.User.ID)
				if ev.User != info.User.ID && strings.Contains(ev.Text, handle) {
					ic := incoming.Disect(ev, handle)
					messenger.ReactOnCall(ic, ev.Channel)
				}

			case *slack.RTMError:
				fmt.Printf("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				fmt.Println("Invalid credentials")
				os.Exit(1)

			default:
				//Take no action
			}
		}
	}
}

func getToken() string {
	token := os.Getenv("SLACKKEY")

	if token == "" {
		fmt.Println("Empty token. (env SLACKKEY needed)")
		os.Exit(1)
	}

	return token
}
