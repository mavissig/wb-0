package nats

import (
	"github.com/nats-io/stan.go"
	"log"
)

type Service struct {
	sc  stan.Conn
	sub stan.Subscription
}

type Consumer interface {
	Consume(data []byte) error
}

func (srv *Service) Connect(clientID string) error {
	sc, err := stan.Connect("wb_cluster", clientID, stan.NatsURL("0.0.0.0:4222"))
	if err != nil {
		log.Println(clientID, "err: ", err)
		return err
	}

	srv.sc = sc
	return err
}

func (srv *Service) Subscribe(subject string, consumer Consumer) error {
	sub, err := srv.sc.Subscribe(subject, func(m *stan.Msg) {
		err := consumer.Consume(m.Data)
		if err != nil {
			log.Println(err)
			return
		}
	})
	if err != nil {
		log.Println(err)
		return err
	}

	srv.sub = sub
	return nil
}

func (srv *Service) Publish(subject string, data []byte) {
	err := srv.sc.Publish(subject, data)
	if err != nil {
		log.Println("error published: ", err)
		return
	}
}

func (srv *Service) Close() {
	if srv.sc != nil {
		_ = srv.sc.Close()
	}
	if srv.sub != nil {
		_ = srv.sub.Unsubscribe()
	}

}
