package memory

import "slackbot/storage"

type Cache struct {
	messages map[string][]string
}

func (c *Cache) Add(user string, msg string) {
	if c.messages == nil {
		c.messages = make(map[string][]string, 10)
	}
	current, ok := c.messages[user]

	if !ok {
		b := []string{msg}
		c.messages[user] = b
		return
	}

	c.messages[user] = append(current, msg)
}

func (c *Cache) Get(user string) string {
	set, ok := c.messages[user]
	if !ok {
		return storage.Nothing
	}
	return set[storage.Random(len(set))]
}
