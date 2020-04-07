package notifier

import (
	"encoding/json"
	"log"

	"github.com/my1562/api/config"
	"github.com/my1562/api/models"
	"github.com/streadway/amqp"
)

type Notifier struct {
	qName string
	ch    *amqp.Channel
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func NewNotifier(config *config.Config) *Notifier {
	conn, err := amqp.Dial(config.RabbitmqURL)
	failOnError(err, "Failed to connect to RabbitMQ")
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")
	return &Notifier{
		qName: q.Name,
		ch:    ch,
	}
}

func (me *Notifier) NotifyServiceMessageChange(chatIDs []int64, message string, addressString string, addressStatus models.AddressArCheckStatus) error {
	var err error

	introduction := ""
	emojiIcon := ""

	if addressStatus == models.AddressStatusNoWork {
		introduction = "Работы не проводятся"
		emojiIcon = "✅"
	}
	if addressStatus == models.AddressStatusWork {
		introduction = ""
		emojiIcon = "🛠"
	}

	fullMessageText := emojiIcon + " " + addressString + ": " + introduction + "\n\n" + message

	for _, chatID := range chatIDs {
		body, err := json.Marshal(map[string]interface{}{
			"ChatID":  chatID,
			"Message": fullMessageText,
		})
		if err != nil {
			continue
		}

		err = me.ch.Publish(
			"",       // exchange
			me.qName, // routing key
			false,    // mandatory
			false,    // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        body,
			})
	}

	if err != nil {
		return err
	}
	return nil
}
