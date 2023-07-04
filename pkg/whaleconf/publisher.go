package whaleconf

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/fatih/color"
)

func (l Conf) MatchingPublisher() *Publisher {
	if l.MQ.Endpoint == "" {
		fmt.Println(color.HiRedString("MQ is not configured"))
		return &Publisher{}
	}
	producer := l.MQ.CreateProducer("event")
	return &Publisher{publisher: producer}
}

type Publisher struct {
	publisher pulsar.Producer
}

func (p *Publisher) Pub(ctx context.Context, key string, data any) error {
	if p.publisher == nil {
		return nil
	}
	buff, err := json.Marshal(data)
	if err != nil {
		_, err = p.publisher.Send(ctx, &pulsar.ProducerMessage{
			Payload: buff,
			Key:     key,
		})
	}
	return err
}
