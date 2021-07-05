package rabbit

import (
	"context"
	"fmt"
	"github.com/streadway/amqp"
)

type Message struct {
	Queue	string
	Body	*[]byte
	ContentType	string
}

func Publish(ctx context.Context, message Message) error {
	URI := fmt.Sprintf("amqp://%s:%s@%s:%d/", "backendtest", "backendtest","localhost", 5672)
	conn, err := amqp.Dial(URI)
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
		message.Queue, // name
		true,        // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)

	if err != nil {
		return err
	}

	err = ch.Confirm(false)
	if err != nil {
		return err
	}

	//confirms := ch.NotifyPublish(make(chan amqp.Confirmation, 1))

	err = ch.Publish(
		"",     // exchange
		message.Queue, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			DeliveryMode:  amqp.Persistent,
			ContentType:   message.ContentType,
			Body:          *message.Body,
		})
	if err != nil {
		return err
	}

	//confirmed := <-confirms
	//if confirmed.Ack {
	//	_ = "Message published"
	//	return nil
	//}

	return nil
}
