include $(GOROOT)/src/Make.inc

TARG=irc
GOFILES=\
    conn.go\
	process.go\
	user.go\
	messages/join.go\
	messages/message.go\
	messages/part.go\
	messages/ping.go\
	messages/private.go\
	messages/quit.go\
	messages/user.go\
	messages/pass.go\
	messages/nick.go

include $(GOROOT)/src/Make.pkg
