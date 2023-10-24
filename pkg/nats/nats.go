package nats

import (
	"github.com/nats-io/stan.go"
	"log"
)

type Service struct {
	sc  stan.Conn
	sub stan.Subscription
}

type Server struct {
	ClusterID string
	ClientID  string
	Url       string
}

func (srv *Server) Subscribe(subject string) (*Service, error) {
	sc, err := stan.Connect(srv.ClusterID, srv.ClientID, stan.NatsURL(srv.Url))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	sub, err := sc.Subscribe(subject, func(m *stan.Msg) {
		log.Printf("Received a message: %s\n", string(m.Data))
	})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	s := Service{
		sc:  sc,
		sub: sub,
	}
	return &s, nil
}

func (s *Service) Publish(subject string, data []byte) {
	err := s.sc.Publish(subject, data)
	if err != nil {
		log.Println("error published: ", err)
		return
	}
}

func (s *Service) Close() {
	_ = s.sc.Close()
	_ = s.sub.Unsubscribe()
}
