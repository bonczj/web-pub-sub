package pubsub

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type PubsubSuite struct{ suite.Suite }

func (suite *PubsubSuite) SetupTest() {
	// easiest way to ensure we have clean data each time is to
	// forceabily clear out any existing instance data within pubsub
	instance = nil
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(PubsubSuite))
}

func (suite *PubsubSuite) TestSingleton() {
	inst := Instance()

	// add a few subscribers
	for i := 0; i < 10; i++ {
		id := uuid.New().String()
		ch := make(chan []byte)
		inst.Subscribe(id, ch)
	}

	inst2 := Instance()

	// both instances must have the same count and IDs of subscribers
	suite.Equal(len(inst.Subscribers), len(inst2.Subscribers))

	for id := range inst.Subscribers {
		_, found := inst2.Subscribers[id]
		suite.True(found)
	}
}

func (suite *PubsubSuite) TestSubscribeSameId() {
	inst := Instance()
	id := uuid.New().String()
	ch := make(chan []byte)

	inst.Subscribe(id, ch)
	inst.Subscribe(id, ch)
	suite.Equal(1, len(inst.Subscribers))
}

func (suite *PubsubSuite) TestUnsubscribe() {
	inst := Instance()

	// add a few subscribers
	for i := 0; i < 10; i++ {
		id := uuid.New().String()
		ch := make(chan []byte)
		inst.Subscribe(id, ch)
		inst.Unsubscribe(id)
	}

	suite.Equal(0, len(inst.Subscribers))
}

func (suite *PubsubSuite) TestClear() {
	inst := Instance()

	// add a few subscribers
	for i := 0; i < 10; i++ {
		id := uuid.New().String()
		ch := make(chan []byte)
		inst.Subscribe(id, ch)
	}

	suite.Equal(10, len(inst.Subscribers))
	inst.Clear()
	suite.Equal(0, len(inst.Subscribers))

}

func (suite *PubsubSuite) TestPublish() {
	inst := Instance()
	ch := make(chan []byte, 10)

	// add a few subscribers
	// we are going to cheat and use a single channel
	// for all subscribers
	for i := 0; i < 3; i++ {
		id := uuid.New().String()
		inst.Subscribe(id, ch)
	}

	// publish several messages
	for i := 0; i < 3; i++ {
		msg := []byte(uuid.NewString())
		inst.Publish(msg)
	}

	// read messages from the channel and count how many there are
	close(ch) // to avoid blocking after reading last message
	count := 0
	done := false
	for !done {
		_, ok := <-ch
		if !ok {
			done = true
		} else {
			count++
		}
	}

	// we should have read exactly 9 messages (3 messages sent to three different subscribers)
	suite.Equal(9, count)
}
