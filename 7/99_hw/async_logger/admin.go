package main

import "time"

type ServeAdmin struct {
	*Serve
}

func (s *ServeAdmin) RemoveListener(server Admin_LoggingServer) {
	for i, listener := range *s.listeners {
		if listener == server {
			s.mu.Lock()
			*s.listeners = append((*s.listeners)[:i], (*s.listeners)[i+1:]...)
			s.mu.Unlock()
			break
		}
	}
}

func (s *ServeAdmin) Logging(nothing *Nothing, server Admin_LoggingServer) error {
	ctx := server.Context()
	time.Sleep(1)
	s.mu.Lock()
	*s.listeners = append(*(s.listeners), server)
	s.mu.Unlock()
	defer s.RemoveListener(server)

	<-ctx.Done()
	return nil
}

func (s *ServeAdmin) Statistics(interval *StatInterval, server Admin_StatisticsServer) error {
	timer := time.NewTimer(time.Duration(interval.IntervalSeconds) * time.Second)
	stat := &Stat{
		ByMethod:   make(map[string]uint64),
		ByConsumer: make(map[string]uint64),
	}

	s.mu.Lock()
	s.stats[server] = stat
	s.mu.Unlock()
	for {
		select {
		case <-server.Context().Done():
			s.mu.Lock()
			delete(s.stats, server)
			s.mu.Unlock()
			return nil
		case <-timer.C:
			err := server.Send(&Stat{
				Timestamp:  time.Now().Unix(),
				ByMethod:   stat.ByMethod,
				ByConsumer: stat.ByConsumer,
			})
			if err != nil {
				return err
			}
			stat.ByMethod = make(map[string]uint64)
			stat.ByConsumer = make(map[string]uint64)
			timer.Reset(time.Duration(interval.IntervalSeconds) * time.Second)
		}
	}
}

func (s *ServeAdmin) mustEmbedUnimplementedAdminServer() {
	//TODO implement me
	panic("implement me")
}
