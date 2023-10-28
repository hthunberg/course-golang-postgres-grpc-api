package integration

import (
	"fmt"

	tc "github.com/testcontainers/testcontainers-go"
)

const lastMessage = "DONE"

type TestLogConsumer struct {
	Msgs []string
	Done chan bool

	// Accepted provides a blocking way of ensuring the logs messages have been consumed.
	// This allows for proper synchronization during Test_StartStop in particular.
	// Please see func devNullAcceptorChan if you're not interested in this synchronization.
	Accepted chan string
}

func newTestLogConsumer(msgs []string, done chan bool) TestLogConsumer {
	return TestLogConsumer{
		Msgs:     msgs,
		Done:     done,
		Accepted: devNullAcceptorChan(),
	}
}

func (g *TestLogConsumer) Accept(l tc.Log) {
	s := string(l.Content)
	if s == fmt.Sprintf("echo %s\n", lastMessage) {
		g.Done <- true
		return
	}
	g.Accepted <- s
	g.Msgs = append(g.Msgs, s)
}

// devNullAcceptorChan returns string channel that essentially sends all strings to dev null
func devNullAcceptorChan() chan string {
	c := make(chan string)
	go func(c <-chan string) {
		//revive:disable
		for range c {
			// do nothing, just pull off channel
		}
		//revive:enable
	}(c)
	return c
}
