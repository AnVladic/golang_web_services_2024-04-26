package main

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"sync"
)

type Serve struct {
	accessUserMap map[string][]string
	logChannel    chan *Event
	listeners     *[]Admin_LoggingServer
	mu            sync.Mutex
	stats         map[Admin_StatisticsServer]*Stat
}

func (s *Serve) authInterceptor(ctx context.Context) error {
	methodName, ok := grpc.Method(ctx)
	if !ok {
		return status.Error(codes.Unknown, "Unknown method")
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Error(codes.Unauthenticated, "metadata not found in context")
	}

	consumers := md["consumer"]
	if len(consumers) == 0 {
		return status.Error(codes.Unauthenticated, "metadata not found in context")
	}
	consumer := consumers[0]
	accessMethods := s.accessUserMap[consumer]
	if accessMethods == nil {
		return status.Error(codes.Unauthenticated, "access method not found in context")
	}

	if !isMethodAllowed(accessMethods, methodName) {
		return status.Error(
			codes.Unauthenticated, "access method not found for this consumer")
	}

	return nil
}

func (s *Serve) logsInterceptor(ctx context.Context) error {
	methodName, ok := grpc.Method(ctx)
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Error(codes.Unauthenticated, "metadata not found in context")
	}
	consumers := md["consumer"]
	if len(consumers) == 0 {
		return status.Error(codes.Unauthenticated, "metadata not found in context")
	}
	p, ok := peer.FromContext(ctx)
	if !ok {
		return status.Error(codes.Unauthenticated, "remove address not found")
	}
	consumer := consumers[0]

	for _, value := range s.stats {
		value.ByConsumer[consumer]++
		value.ByMethod[methodName]++
	}

	s.logChannel <- &Event{
		Consumer: consumer,
		Method:   methodName,
		Host:     p.Addr.String(),
	}
	return nil
}
