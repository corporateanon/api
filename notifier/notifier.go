package notifier

import (
	"log"

	"github.com/hibiken/asynq"
	"github.com/my1562/api/config"
	"github.com/my1562/api/models"
	"github.com/my1562/queue"
)

type Notifier struct {
	client *asynq.Client
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func NewNotifier(config *config.Config) *Notifier {
	client := asynq.NewClient(
		asynq.RedisClientOpt{
			Addr: config.Redis,
		},
	)
	return &Notifier{client: client}
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
		me.client.Enqueue(queue.NewNotifyTask(chatID, fullMessageText))
	}

	if err != nil {
		return err
	}
	return nil
}
