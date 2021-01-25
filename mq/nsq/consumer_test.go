package nsq

import (
	"context"
	"errors"
	"fmt"
	"github.com/arden/easy/core/logger"
	"testing"

	"github.com/nsqio/go-nsq"
)

func TestNsqConsumer_AddHandler(t *testing.T) {
	consumer := &NsqHandler{
		Topic:            "test",
		Channel:          "ai",
		OpenChannelTopic: true,
	}

	consumer.SetHandle(func(log *logger.Logger, message *nsq.Message) error {
		fmt.Println(string(message.Body))
		return errors.New("error")
	})

	ctx, _ := context.WithCancel(context.Background())
	err := RunMock(ctx, consumer, MockMessage([]byte("hello")))

	if err != nil {
		t.Error(err)
	}
}
