package whaleconf

import (
	"fmt"
	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/fatih/color"
)

func (l Conf) CreateSubscriber() *Subscriber {
	if l.MQ.UserServiceReader.Endpoint == "" {
		fmt.Println(color.HiRedString("MQ is not configured"))
		return &Subscriber{}
	}
	userLevelChange := l.MQ.UserServiceReader.CreateConsumer("level-change", "matching")
	return &Subscriber{
		UserLevelChange: userLevelChange,
	}
}

type Subscriber struct {
	UserLevelChange pulsar.Consumer
}
