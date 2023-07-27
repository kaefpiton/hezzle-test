package natsClient

import (
	"context"
	"encoding/json"
	"github.com/nats-io/nats.go"
	"hezzle/internal/infrastructure/queue"
	"hezzle/internal/infrastructure/usecase/repository"
	"hezzle/pkg/logger"
)

const EventTopicName = "events"

type EventListener struct {
	sub              queue.Subscriber
	eventsRepository repository.EventsRepository
	logger           logger.Logger
	ctx              context.Context
}

func NewEventListener(ctx context.Context, sub queue.Subscriber, logRepository repository.EventsRepository, logger logger.Logger) *EventListener {
	return &EventListener{
		sub:              sub,
		eventsRepository: logRepository,
		logger:           logger,
		ctx:              ctx,
	}
}

func (l *EventListener) ListenTopic() {
	l.logger.Info("Event Listener started!")
	unsub, err := l.sub.Sub(EventTopicName, func(m *nats.Msg) {
		l.logger.Info("Received a message: %s\n", string(m.Data))

		var itemModel repository.ItemModel
		err := json.Unmarshal(m.Data, &itemModel)
		if err != nil {
			l.logger.Error(err)
		}

		EventModel := repository.ItemModelToEvent(itemModel)
		err = l.eventsRepository.Create(EventModel)
		if err != nil {
			l.logger.Error(err)
		}
	})
	if err != nil {
		l.logger.ErrorF("could not subscribe to topic %s: %w", EventTopicName, err)

	}

	go func() {
		<-l.ctx.Done()
		l.logger.Info("Stop listen events!")
		if err := unsub(); err != nil {
			panic(err)
		}
	}()
}
