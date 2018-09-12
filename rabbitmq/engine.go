package rabbitmq

import (
	"encoding/json"

	"github.com/Proofsuite/amp-matching-engine/types"
)

func (c *Connection) SubscribeEngineResponses(fn func(*types.EngineResponse) error) error {
	ch := c.GetChannel("erSub")
	q := c.GetQueue(ch, "engineResponse")

	go func() {
		msgs, err := ch.Consume(
			q.Name, // queue
			"",     // consumer
			true,   // auto-ack
			false,  // exclusive
			false,  // no-local
			false,  // no-wait
			nil,    // args
		)

		if err != nil {
			logger.Fatal("Failed to register a consumer:", err)
		}

		forever := make(chan bool)

		go func() {
			for d := range msgs {
				var res *types.EngineResponse
				err := json.Unmarshal(d.Body, &res)
				if err != nil {
					logger.Error(err)
					continue
				}
				go fn(res)
			}
		}()

		<-forever
	}()
	return nil
}

func (c *Connection) PublishEngineResponse(res *types.EngineResponse) error {
	ch := c.GetChannel("erPub")
	q := c.GetQueue(ch, "engineResponse")

	bytes, err := json.Marshal(res)
	if err != nil {
		logger.Error("Failed to marshal engine response: ", err)
		return err
	}

	err = c.Publish(ch, q, bytes)
	if err != nil {
		logger.Error("Failed to publish order: ", err)
		return err
	}

	return nil
}
