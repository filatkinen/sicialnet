package consumer

import (
	"errors"
	"log"

	"github.com/filatkinen/socialnet/internal/rabbit"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Comsumer struct {
	log     *log.Logger
	config  rabbit.Config
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
	chExit  chan struct{}
	chMsg   <-chan amqp.Delivery
}

func NewConsumer(config rabbit.Config, log *log.Logger, clientID string) (*Comsumer, error) {
	c := Comsumer{log: log, config: config}

	configAmpq := amqp.Config{
		Vhost:      "/",
		Properties: amqp.NewConnectionProperties(),
	}
	configAmpq.Properties.SetClientConnectionName(clientID)

	conn, err := amqp.Dial(config.GetDSN())
	if err != nil {
		e := c.Close()
		return nil, errors.Join(err, e)
	}
	c.conn = conn

	channel, err := conn.Channel()
	if err != nil {
		e := c.Close()
		return nil, errors.Join(err, e)
	}
	c.channel = channel

	err = channel.ExchangeDeclare(
		config.ExchangeName, // name
		"direct",            // type
		true,                // durable
		false,               // auto-deleted
		false,               // internal
		false,               // no-wait
		nil,                 // arguments
	)
	if err != nil {
		e := c.Close()
		return nil, errors.Join(err, e)
	}

	queue, err := channel.QueueDeclare(
		"",    // name of the queue
		false, // durable
		true,  // delete when unused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		e := c.Close()
		return nil, errors.Join(err, e)
	}
	c.queue = queue

	err = channel.QueueBind(
		queue.Name,          // queue name
		clientID,            // routing key
		config.ExchangeName, // exchange
		false,
		nil,
	)
	if err != nil {
		e := c.Close()
		return nil, errors.Join(err, e)
	}

	c.chExit = make(chan struct{})

	msgs, err := c.channel.Consume(
		queue.Name, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		e := c.Close()
		return nil, errors.Join(err, e)
	}
	c.chMsg = msgs
	return &c, nil
}

func (c *Comsumer) Start(f func([]byte)) {
	c.log.Println("Starting ws message reader")
	go func() {
		for d := range c.chMsg {
			f(d.Body)
		}
	}()
	<-c.chExit
}

func (c *Comsumer) Stop() {
	c.log.Println("Stopping ws message reader")
	c.chExit <- struct{}{}
}

func (c *Comsumer) Close() (err error) {
	c.log.Println("Closing ws message reader")
	if c.channel != nil {
		if e := c.channel.Close(); e != nil {
			err = errors.Join(err, e)
		}
	}
	if c.conn != nil {
		if e := c.conn.Close(); e != nil {
			err = errors.Join(err, e)
		}
	}
	return err
}
