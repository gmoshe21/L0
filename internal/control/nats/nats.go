package nats

import (
	"L0/internal/control"
	"L0/internal/models"
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/nats-io/stan.go"
)

type subscriber struct {
	log      *log.Logger
	natsConn stan.Conn
	orderUC  control.UseCase
}

func NewSubscriber(natsConn stan.Conn, orderUC control.UseCase) *subscriber {
	return &subscriber{natsConn: natsConn, orderUC: orderUC}
}

func (s *subscriber) Subscribe(subject, qgroup string, workerNum int, cb stan.MsgHandler) {
	wg := &sync.WaitGroup{}

	for i := 0; i <= workerNum; i++ {
		wg.Add(1)
		go s.runWorker(
			wg,
			i,
			s.natsConn,
			subject,
			qgroup,
			cb,
		)
	}
	wg.Wait()
}

func (s *subscriber) runWorker(
	wg *sync.WaitGroup,
	workerID int,
	conn stan.Conn,
	subject string,
	qgroup string,
	cb stan.MsgHandler,
	opts ...stan.SubscriptionOption,
) {
	defer wg.Done()

	_, err := conn.QueueSubscribe(subject, qgroup, cb, opts...)
	if err != nil {
		log.Printf("WorkerID: %v, QueueSubscribe: %v", workerID, err)
		if err := conn.Close(); err != nil {
			log.Printf("WorkerID: %v, conn.Close error: %v", workerID, err)
		}
	}

}

func (s *subscriber) Run(ctx context.Context) {
	go s.Subscribe("create", "service", 0, s.processCreateOrder(ctx))
}

func (s *subscriber) processCreateOrder(ctx context.Context) stan.MsgHandler {
	return func(msg *stan.Msg) {
		var params models.Order

			err := json.Unmarshal(msg.Data, &params)
			if err != nil {
				log.Println(err)
				log.Println("subscriber.Run.Unmarshal():", err)
			}
			
			err = s.orderUC.NewOrder(ctx, params)
			if err != nil {
				log.Println(err)
			}
	}
}