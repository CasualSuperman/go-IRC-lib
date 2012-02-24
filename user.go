package irc

type user struct {
	Nick, Username, Name, Pass string
	Hidden bool
}

func NewUser(nick, username string, hidden bool, name, pass string) user {
	return user{nick, username, name, pass, hidden}
}

func (u user) PassMessage() Message {
	return NewPassMessage(u.Pass)
}

func (u user) UserMessage() Message {
	return NewUserMessage(u.Username, u.Hidden, u.Name)
}

func (u user) NickMessage() Message {
	return NewNickMessage(u.Nick)
}
