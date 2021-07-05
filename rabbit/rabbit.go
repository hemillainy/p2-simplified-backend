package rabbit

import (
	"context"
	"fmt"
	"github.com/streadway/amqp"
	"os"
)

type Message struct {
	Queue       string
	Body        *[]byte
	ContentType string
}

func Publish(ctx context.Context, message Message) error {
	URI := fmt.Sprintf("amqp://%s:%s@%s:%d/",
		os.Getenv("BACKEND_BROKER_USER"),
		os.Getenv("BACKEND_BROKER_PASSWORD"),
		os.Getenv("BACKEND_BROKER_HOST"),
		os.Getenv("BACKEND_BROKER_PORT"),
	)
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
		message.Queue, 
		true,          
		false,         
		false,         
		false,       
		nil,           
	)

	if err != nil {
		return err
	}

	err = ch.Confirm(false)
	if err != nil {
		return err
	}
	err = ch.Publish(
		"",            
		message.Queue, 
		false,        
		false,        
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  message.ContentType,
			Body:         *message.Body,
		})
	if err != nil {
		return err
	}

	return nil
}
