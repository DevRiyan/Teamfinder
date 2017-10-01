package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func stringRPC(n string) (res string, err error) {
	conn, err := amqp.Dial("amqp://admin:rian@185.53.130.233:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)

	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	err = ch.Publish(
		"",          // exchange
		"rpc_queue", // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			ReplyTo:     q.Name,
			Body:        []byte(n),
		})

	failOnError(err, "Failed to publish a message")

	for d := range msgs {
		//	if corrId == d.CorrelationId {
		res = (string(d.Body))
		log.Printf(" [.] Got1 = %s", res)
		// failOnError(err, "Failed to convert body to integer")
		break
		// }
	}
	log.Printf(" [.] Got2 = %s", res)
	return
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	n := "aloha"

	log.Printf(" [x] Send Messange(%s)", n)
	res, err := stringRPC(n)
	failOnError(err, "Failed to handle RPC request")

	log.Printf(" [.] Got = %s", res)
}
